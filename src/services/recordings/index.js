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
    'b2': Storage.B2Storage
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

const Service = module.exports = {

    getFileFromUUID: (uuid, option, cb) => {
        application.DB._query.file.get(uuid, option.pathName, option.contentType, (err, res) => {
            if (err)
                return cb(err);

            let fileDb = res && res[0];
            if (!fileDb || !fileDb.path)
                return cb(new CodeError(404, `File ${uuid} not found!`));

            let providerName = FILE_TYPES[fileDb.type];

            if (!providerName)
                return cb(new CodeError(500, `Bad file provider.`));

            if (option.range) {
                option.range = httpUtil.readRangeHeader(option.range, fileDb.size);
            }

            let result = {
                source: null,
                totalLength: fileDb.size
            };

            if (fileDb.private) {
                application.DB._query.domain.getByName(fileDb.domain, 'storage', (err, domainConfig) => {
                    if (err)
                        return cb(err);

                    if (useDefaultStorage(domainConfig)) {
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
                "createdOn": Date.now(),
                "size": fileConf.contentLength
            };
            if (response.hasOwnProperty('bucketName'))
                doc.bucketName = response.bucketName;

            application.DB._query.file.insert(doc, (err) => {
                if (err)
                    return cb(err);
                return cb(err, fileConf);
            });
        }

    }
};


function getProvider(domainName, storageConf, nameProvider) {
    let name = !nameProvider ? storageConf.defaultProvider : nameProvider,
        id = `${domainName}_${name}`;

    let provider = cache.get(id);
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
    if (providers instanceof Array) {
        for (let config of providers)
            if (config.type === name)
                return config;
    }
}

function useDefaultStorage(domainConfig) {
    return !domainConfig || !domainConfig.storage || domainConfig.storage.defaultProvider == 'local' || !STORAGES[domainConfig.storage.defaultProvider]
}