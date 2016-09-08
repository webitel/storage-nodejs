/**
 * Created by igor on 29.08.16.
 */

"use strict";

const Storage = require(__appRoot + '/lib/storage'),
    helper = require('./helper'),
    util = require('util'),
    CacheCollection = require(__appRoot + '/lib/collection'),
    log = require(__appRoot + '/lib/log')(module),
    cache = new CacheCollection('id'),
    CodeError = require(__appRoot + '/lib/error'),
    httpUtil = require(__appRoot + '/utils/http'),
    
    FILE_TYPES = ['local', 's3', 'b2'],
    DEF_ID = '_default_'
    ;


const STORAGES = {
    'local': Storage.LocalStorage,
    'b2': Storage.B2Storage,
    's3': Storage.S3Storage
};

if (!STORAGES.hasOwnProperty(helper.DEFAULT_PROVIDER_NAME))
    throw `Not found default provider ${helper.DEFAULT_PROVIDER_NAME}!`;

const DEFAULT_PROVIDER_CONFIG = helper.getProviderConfigByName(helper.DEFAULT_PROVIDER_NAME);

if (!DEFAULT_PROVIDER_CONFIG)
    throw `Please set config provider ${helper.DEFAULT_PROVIDER_NAME}!`;

const defProvider = getProvider(DEF_ID, helper.DEFAULT_PROVIDERS_CONF, helper.DEFAULT_PROVIDER_NAME);

if (!defProvider)
    throw `Error create default provider ${helper.DEFAULT_PROVIDER_NAME}`;
else log.info(`Set default provider ${helper.DEFAULT_PROVIDER_NAME} - OK`);

let executeRemoveNonExistentFiles = false;

const Service = module.exports = {

    getFileFromUUID: (caller, uuid, option, cb) => {
        application.DB._query.file.get(uuid, option.pathName, option.contentType, (err, res) => {
            if (err)
                return cb(err);

            if (option.contentType === 'all') {
                return cb(null, res);
            }

            let fileDb = res && res[0];
            if (!fileDb || !fileDb.path)
                return cb(new CodeError(404, `File ${uuid} not found!`));

            if (caller.domain && caller.domain != fileDb.domain)
                return cb(new CodeError(403, "Permission denied!"));

            let providerName = FILE_TYPES[fileDb.type];

            if (!providerName)
                return cb(new CodeError(500, `Bad file provider.`));

            let result = {
                source: null,
                totalLength: fileDb.size
            };

            if (option.range) {
                option.range = httpUtil.readRangeHeader(option.range, fileDb.size);
            }

            if (fileDb.private) {
                application.DB._query.domain.getByName(fileDb.domain, 'storage', (err, domainConfig) => {
                    if (err)
                        return cb(err);

                    if (!domainConfig || !domainConfig.storage) {
                        return cb(new CodeError(400, `Please set domain storage config!`))
                    } else {
                        let provider = getProvider(fileDb.domain, domainConfig.storage, providerName);

                        if (!provider)
                            return cb(new CodeError(400, `Bad provider config.`));

                        return provider.get(fileDb, option, sendResponse);
                    }
                })
            } else {
                let provider = getProvider(DEF_ID, helper.DEFAULT_PROVIDERS_CONF, providerName);
                if (!provider)
                    return cb(new CodeError(400, `Bad provider config.`));
                return provider.get(fileDb, option, sendResponse);
            }

            function sendResponse(err, source) {
                if (err)
                    return cb(err);
                result.source = source;
                return cb(null, result)
            }
        });
    },

    saveFile: (fileConf, option, cb) => {
        application.DB._query.domain.getByName(option.domain, 'storage', (err, domainConfig) => {
            if (err)
                log.error(err);

            fileConf.domain = option.domain;

            let provider,
                domainId,
                config;

            if (useDefaultStorage(domainConfig)) {
                domainId = DEF_ID;
                config = helper.DEFAULT_PROVIDERS_CONF;
            } else {
                fileConf.private = true;
                domainId = option.domain;
                config = domainConfig.storage;
            }

            provider = getProvider(domainId, config);

            if (!provider)
                return cb(new CodeError(400, `Bad provider config.`));

            return provider.save(fileConf, option, saveToDb);
        });

        function saveToDb(err, response = {}) {
            if (err)
                return cb(err);

            let doc = {
                "uuid": fileConf.uuid,
                "name": fileConf.queryName,
                "path": response.path,
                "domain": fileConf.domain,
                "private": fileConf.private === true,
                "content-type": fileConf.contentType,
                "type": response.type,
                "createdOn": new Date(),
                "size": fileConf.contentLength
            };
            if (response.hasOwnProperty('bucketName'))
                doc.bucketName = response.bucketName;

            if (response.hasOwnProperty('storageFileId'))
                doc.storageFileId = response.storageFileId;

            application.DB._query.file.insert(doc, (err) => {
                if (err)
                    return cb(err);
                return cb(err, fileConf);
            });
        }

    },
    
    // Public

    removeNonExistentFiles: (caller, option, cb) => {
        // TODO add acl
        if (executeRemoveNonExistentFiles)
            return cb(null, 'Process running.');

        if (caller.domain)
            return cb(new CodeError(403, "Permission denied!"));


        let startDate = new Date(option.from),
            endDate = new Date(option.to),
            deleteFileCount = 0;


        if (!startDate || !endDate)
            return cb(new CodeError(400, 'Bad date parameters'));

        let stream = application.DB._query.file.getStreamByDateRange(startDate, endDate);
        
        function done(err) {
            if (err) {
                log.error(err);
                executeRemoveNonExistentFiles = false;
                return stream.close();
            }
            stream.resume();
        }

        stream.on('error', (err) => {
            log.error(err);
            executeRemoveNonExistentFiles = false;
        });
        stream.on('end', () => {
            log.debug(`End stream.`);
            executeRemoveNonExistentFiles = false;
        });

        stream.on('data', (fileDb) => {
            let providerName = FILE_TYPES[fileDb.type];
            if (!providerName)
                return done(new Error(`Bad provider type: ${fileDb.type}`));

            stream.pause();

            if (fileDb.private) {
                application.DB._query.domain.getByName(fileDb.domain, 'storage', (err, domainConfig) => {
                    if (err)
                        return done(err);

                    if (!domainConfig || !domainConfig.storage) {
                        log.warn(`Skip ${fileDb._id}, no configure domain storage ${fileDb.domain}`);
                        return done(null);
                    } else {
                        existsData(getProvider(fileDb.domain, domainConfig.storage, providerName));
                    }
                })
            } else {
                existsData(getProvider(DEF_ID, helper.DEFAULT_PROVIDERS_CONF, providerName));
            }

            function existsData(provider) {
                if (!provider) {
                    log.warn(`Skip ${fileDb._id} not found provider  ${providerName}`);
                    return done();
                }
                log.trace(`try exists file ${fileDb._id} from provider ${provider.name}`);
                provider.existsFile(fileDb, (err, exists) => {
                    if (err)
                        return done(err);

                    log[exists ? 'trace' : 'warn'](`file ${fileDb.uuid} ${fileDb.name}@${fileDb.domain} (${fileDb.private || false}) exists: ${exists}`);
                    if (!exists) {
                        application.DB._query.file.deleteById(fileDb._id, (err) => {
                            if (err)
                                return done(err);
                            deleteFileCount++;
                            log.debug(`delete file ${fileDb.uuid} ${fileDb.name}@${fileDb.domain}`);
                            done(null);
                        });
                    } else {
                        done(null);
                    }
                });
            }
        });

        executeRemoveNonExistentFiles = true;
        cb(null, `Process start.`);
    },

    stats: (caller, option, cb)  => {
        application.DB._query.file.getFilesStats(option.uuid, caller.domain || option.domain, option, cb);
    },

    delFile: (caller, option, cb) => {
        let uuid = option.uuid;
        application.DB._query.file.get(uuid, option.pathName, option.contentType, (err, res) => {
            if (err)
                return cb(err);

            let fileDb = res && res[0];
            if (!fileDb || !fileDb.path)
                return cb(new CodeError(404, `File ${uuid} not found!`));

            if (caller.domain && caller.domain != fileDb.domain)
                return cb(new CodeError(403, "Permission denied!"));

            let providerName = FILE_TYPES[fileDb.type];

            if (!providerName)
                return cb(new CodeError(500, `Bad file provider.`));

            if (fileDb.private) {
                application.DB._query.domain.getByName(fileDb.domain, 'storage', (err, domainConfig) => {
                    if (err)
                        return cb(err);

                    if (!domainConfig || !domainConfig.storage) {
                        return cb(new CodeError(400, `Please set domain storage config!`))
                    } else {
                        let provider = getProvider(fileDb.domain, domainConfig.storage, providerName);

                        if (!provider)
                            return cb(new CodeError(400, `Bad provider config.`));

                        return provider.del(fileDb, sendResponse);
                    }
                })
            } else {
                let provider = getProvider(DEF_ID, helper.DEFAULT_PROVIDERS_CONF, providerName);
                if (!provider)
                    return cb(new CodeError(400, `Bad provider config.`));
                return provider.del(fileDb, sendResponse);
            }

            function sendResponse(err) {
                if (err)
                    return cb(err);

                if (option.delDb) {
                    return application.DB._query.file.deleteById(fileDb._id, (err) => {
                        if (err)
                            return cb(err);
                        return cb(null, fileDb);
                    });
                } else {
                    return cb(null, fileDb);
                }
            }
        });
    }
};


function getProvider(domainName, storageConf, nameProvider) {
    let name = !nameProvider ? storageConf.defaultProvider : nameProvider,
        id = `${domainName}_${name}`;

    let provider = cache.get(id);
    if (provider) {
        let newConf = findProviderConfigByName(storageConf.providers, name);
        if (!provider.checkConfig(newConf, storageConf.maskPath)) {
            log.debug(`recreate storage id: ${id}`);
            cache.remove(id);
            provider = null;
        }
    }

    if (!provider && STORAGES[name]) {
        let configProvider = findProviderConfigByName(storageConf.providers, name);

        if (!configProvider)
            return null;

        provider = new STORAGES[name](configProvider, configProvider.maskPath || storageConf.maskPath);
        cache.add(id, provider);
    }

    return provider
}

function findProviderConfigByName(providers, name) {
    if (providers && providers.hasOwnProperty(name))
        return providers[name];
}

function useDefaultStorage(domainConfig) {
    return !domainConfig || !domainConfig.storage || !domainConfig.storage.defaultProvider
        || domainConfig.storage.defaultProvider == 'local' || !STORAGES[domainConfig.storage.defaultProvider]
}