/**
 * Created by igor on 31.08.16.
 */

"use strict";

const Storage = require(__appRoot + '/lib/storage'),
    helper = require('./helper'),
    util = require('util'),
    log = require(__appRoot + '/lib/log')(module),
    checkPermission = require(__appRoot + '/utils/acl'),
    async = require('async'),
    formidable = require('formidable'),
    CodeError = require(__appRoot + '/lib/error'),
    fileService = require('../file'),
    httpUtil = require(__appRoot + '/utils/http'),

    FILE_TYPES = ['local', 's3', 'b2']
    ;


const STORAGES = {
    'local': Storage.LocalStorage,
    'b2': Storage.B2Storage
};

if (!STORAGES.hasOwnProperty(helper.DEFAULT_PROVIDER_NAME))
    throw `Not found default provider ${helper.DEFAULT_PROVIDER_NAME}!`;

let providers = {};

for (let config of helper.DEFAULT_PROVIDERS_CONF.providers) {
    if (providers.hasOwnProperty(config.type))
        throw `Already exists config provider ${config.type}`;

    providers[config.type] = new STORAGES[config.type](config, helper.mask);

    log.info(`init provider ${config.type}`);
}
log.info(`use default provider ${helper.DEFAULT_PROVIDER_NAME}`);


const Service = module.exports = {

    saveFilesFromRequest: (caller, req, options, cb) => {
        if (!checkPermission(caller.acl, 'cdr/media', 'c'))
            return cb(new CodeError(403, 'Forbidden!'));

        let domain = caller.domain || options.domain;

        if (!domain)
            return (cb(new CodeError(400, 'Domain is required.')));

        if (req.headers['content-length'] && +req.headers['content-length'] > helper.maxFileSize)
            return cb(new CodeError(400, `Bad files size.`));

        let form = new formidable.IncomingForm()
            ;

        form.maxFieldsSize = helper.maxFileSize;
        form.uploadDir = 'cache/';
        form.hash = 'sha1';
        form.on('error', (err) => {
            log.error(err);
            return cb(err)
        });

        form.parse(req, function (err, fields, files) {
            let _files = [];
            for (let key in files) {
                let f = {
                    path: files[key].path,
                    name: files[key].name,
                    contentType: files[key].type,
                    contentLength: files[key].size,
                    domain: domain,
                    type: options.type || 'mp3',
                    sha1: files[key].hash
                };

                _files.push(f);
            }

            async.eachSeries(
                _files,
                (file, cb) => {
                    async.waterfall([
                        (cb) => {
                            providers[helper.DEFAULT_PROVIDER_NAME].save(file, {}, cb);
                        },
                        (fileResponse, cb) => {
                            let doc = {
                                "type": file.type,
                                "name": file.name,
                                "domain": file.domain,
                                "size": file.contentLength,
                                "format": file.contentType,
                                "providerId": fileResponse.type,
                                "path": fileResponse.path
                            };

                            if (fileResponse.hasOwnProperty('bucketName'))
                                doc.bucketName = fileResponse.bucketName;

                            if (fileResponse.hasOwnProperty('storageFileId'))
                                doc.storageFileId = fileResponse.storageFileId;

                            application.DB._query.media.insert(doc, cb);
                        }
                    ], (err) => {
                        fileService.deleteFile(file.path);
                        return cb(err)
                    })
                },
                cb
            )

        });
    },

    list: (caller, options, cb) => {
        if (!checkPermission(caller.acl, 'cdr/media', 'r'))
            return cb(new CodeError(403, 'Forbidden!'));

        let domain = caller.domain || options.domain;

        if (!domain)
            return cb(new CodeError(400, 'Domain is required.'));

        // TODO add pagination
        application.DB._query.media.listFromDomain(domain, cb);
    },

    get: (caller, options, cb) => {
        if (!checkPermission(caller.acl, 'cdr/media', 'r'))
            return cb(new CodeError(403, 'Forbidden!'));

        return Service._get(
            caller.domain || options.domain,
            options,
            cb
        );
    },

    _get: (domain, options, cb) => {

        if (!domain)
            return cb(new CodeError(400, 'Domain is required.'));

        application.DB._query.media.getItem(options.type, options.name, domain, (err, data) => {
            if (err)
                return cb(err);

            if (!data)
                return cb(new CodeError(404, `Not found file ${options.name}`));

            let provider = providers[FILE_TYPES[data.providerId || 0]];
            if (!provider)
                return cb(new CodeError(500, `Not found provider.`));

            let opt = {
                range: null
            };
            if (options.range) {
                opt.range = httpUtil.readRangeHeader(options.range, data.size);
            }

            // TODO
            if (data.format == "audio/mp3")
                data.format = "audio/mpeg";

            let result = {
                source: null,
                range: opt.range,
                totalLength: data.size,
                contentType: data.format
            };

            return provider.get({
                path: (data.path || provider.getFilePath(domain, data.name)),
                bucketName: data.bucketName,
                storageFileId: data.storageFileId
            }, opt, (err, source) => {
                if (err)
                    return cb(err);
                result.source = source;
                return cb(null, result)
            });
        })
    },

    delete: (caller, options, cb) => {
        if (!checkPermission(caller.acl, 'cdr/media', 'd'))
            return cb(new CodeError(403, 'Forbidden!'));

        let domain = caller.domain || options.domain;

        if (!domain)
            return cb(new CodeError(400, 'Domain is required.'));

        application.DB._query.media.getItem(options.type, options.name, domain, (err, data) => {
            if (err)
                return cb(err);

            if (!data)
                return cb(new CodeError(404, `Not found file ${options.name}`));

            let provider = providers[FILE_TYPES[data.providerId || 0]];
            if (!provider)
                return cb(new CodeError(500, `Not found provider.`));

            return provider.del({
                path: (data.path || provider.getFilePath(domain, data.name)),
                bucketName: data.bucketName,
                storageFileId: data.storageFileId
            }, delDb);

            function delDb(err, res) {
                if (err)
                    return cb(err);

                return application.DB._query.media.deleteById(data._id, cb);
            }
        });
    }
};
