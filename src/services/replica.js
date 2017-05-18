/**
 * Created by igor on 21.10.16.
 */

"use strict";

const log = require(`${__appRoot}/lib/log`)(module),
    conf = require(`${__appRoot}/conf`),
    CodeError = require(`${__appRoot}/lib/error`),
    url = require('url'),
    recordingsService = require('./recordings'),
    cdrConf = conf.get('replica:cdr'),
    filesConf = conf.get('replica:files'),
    hostConf = conf.get('replica:host'),
    authConf = conf.get('replica:auth'),
    live = `${conf.get('replica:live')}` === 'true',
    headersConf = conf.get('replica:headers'),
    cdrCollectionName = "cdr", //TODO rename collectionCDR
    fileCollectionName = "cdrFile", //TODO rename collectionFile
    Scheduler = require(`${__appRoot}/lib/scheduler`),
    async = require('async')
;

const getDoc = {
    [cdrCollectionName]: (objId, cb) => {
        application.DB._query.cdr.getByObjId(objId, cb);
    },
    [fileCollectionName]: (objId, cb) => {
        application.DB._query.file.getByObjId(objId, cb);
    }
};

const sendDoc = {
    [cdrCollectionName]: _sendCdr,
    [fileCollectionName]: _sendFile
};

const Service = module.exports = {
    cdr: {},
    files: {},
    provider: null,
    scheduler: null,

    sendCdr: (doc, parentId) => {
        try {
            if (!parentId)
                return log.error(`[cdr] Bad parent id`, doc);

            if (live) {
                _sendCdr(doc, parentId);
            } else {
                Service._saveReplica(cdrCollectionName, parentId)
            }
        } catch (e) {
            Service._onError(e, cdrCollectionName, parentId);
        }
    },

    sendFile: (doc, parentId) => {
        try {
            if (!parentId)
                return log.error(`[file] Bad parent id`, doc);

            if (live) {
                _sendFile(doc, parentId);
            } else {
                Service._saveReplica(fileCollectionName, parentId)
            }
        } catch (e) {
            Service._onError(e, fileCollectionName, parentId);
        }
    },

    _onError: (err, collection, docId) => {
        log.error(err);
        Service._saveReplica(collection, docId);
    },

    _saveReplica: (collection, docId) => {
        application.DB._query.replica.insert({
            createdOn: Date.now(),
            collection: collection,
            docId: docId
        }, e => {
            if (e)
                log.error(e)
        });
    },

    _sendRequest: (params, body, cb) => {
        if (!Service.provider)
            return cb(new CodeError(500, `No initialize replica provider`));

        const req = Service.provider.request(params, cb);
        req.on('error', cb);
        if (body)
            req.write(body);
        req.end();
    },

    _sendStream: (params, stream, cb) => {
        try {
            if (!Service.provider)
                return cb(new CodeError(500, `No initialize replica provider`));

            const req = Service.provider.request(params, cb);
            req.on('error', cb);
            stream.pipe(req);
        } catch (e) {
            return cb(e);
        }
    },

    _huntOne: (cb) => {

        async.waterfall([
            application.DB._query.replica.getOne,
            (doc, cb) => {
                const val = doc && doc.value;
                if (!val)
                    return cb(new CodeError(404, `Not found`));

                if (typeof getDoc[val.collection] !== "function" || typeof sendDoc[val.collection] !== "function")
                    return cb(new CodeError(500, `Bad collection name ${val.collection}`));

                const fnName = val.collection;

                getDoc[fnName](val.docId, (err, res) => {
                    if (err)
                        return cb(err);

                    if (!res)
                        return cb(new CodeError(404, `Not found ${val.docId}`));

                    sendDoc[fnName](res, res._id, cb);
                })
            }
        ], cb);
    },

    process: (cb) => {
        log.trace(`Start export.`);
        const cdrRequest = (data, cb) => {
            _sendCdr(data, null, (err, res) => {
                if (err)
                    return cb(err);

                let response = '';
                res.on('error', cb);
                res.on('data', chunk => response += chunk);
                res.on('end', () => {
                    let json;
                    try {
                        json = JSON.parse(response);
                        if (!(json.items instanceof Array))
                            return cb(new Error("Bad response replica client"));

                        if (!json.errors) // OK save to client;
                            return cb(null, null);

                        const errorsIds = [];
                        json.items.forEach(item => {
                            if (item.update.status < 200 || item.update.status > 299) {
                                errorsIds.push(item.update._id)
                            }
                        });
                        return cb(null, errorsIds);
                    } catch (e) {
                        log.error(e);
                        return cb(e);
                    }
                });
            });
        };
        const maxCountCdr = 10000;
        const exportCdr = (cb) => {
            console.time(`Replica cdr`);
            application.DB._query.replica.sync(
                cdrCollectionName,
                maxCountCdr,
                cdrRequest,
                (err, count) => {
                    console.timeEnd(`Replica cdr`);
                    if (err) {
                        log.error(err);
                        return cb(err);
                    }
                    if (count && count > 0 && count >= maxCountCdr) {
                        exportCdr();
                    } else {
                        log.debug(`End export CDR data.`);
                        return cb();
                    }
                }
            );
        };

        const exportOtherFile = (cb) => {
            const process = err => {
                if (err && err.status === 404) {
                    return cb();
                } else if (err) {
                    log.error(err);
                    return cb();
                } else {
                    Service._huntOne(process);
                }
            };
            Service._huntOne(process);
        };

        async.waterfall(
            [
                exportCdr,
                exportOtherFile
                
            ],
            err => {
                if (err)
                    log.error(err);
                log.debug(`End export replica data.`);
                cb();
            }
        );
    },

    _init: () => {
        const cdr = {headers: {}},
            files = {headers: {}},
            hostParams = url.parse(hostConf)
        ;

        for (let key in headersConf)
            if (headersConf.hasOwnProperty(key))
                cdr.headers[key] = files.headers[key] = headersConf[key];

        // TODO
        if (authConf && authConf.type === 'base') {
            cdr.headers['Authorization'] = files.headers['Authorization']
                = `Basic ${new Buffer(`${authConf.login}:${authConf.password}`).toString('base64')}`;
        }

        cdr.headers['Content-Type'] = "application/json";

        Service.provider = ~hostParams.protocol.indexOf('https:') ? require('https') : require('http');

        if (`${cdrConf.keepAlive}` === 'true')
            cdr.agent = new Service.provider.Agent({keepAlive: true});

        if (`${filesConf.keepAlive}` === 'true')
            files.agent = new Service.provider.Agent({keepAlive: true});

        cdr.host = files.host = hostParams.hostname;
        cdr.port = files.port = hostParams.port;
        
        cdr.method = (cdrConf.method || "POST").toUpperCase();
        cdr.path = cdrConf.path;
        
        files.method = (filesConf.method || "POST").toUpperCase();
        files.path = filesConf.path;

        Service.cdr = cdr;
        Service.files = files;

        // Service.process( () => {});
        Service.scheduler = new Scheduler(conf.get('replica:cronJob'), Service.process);
    },

    getCdrRequestParams: () => {
        if (!Service.cdr) return null;

        return {
            host: Service.cdr.host,
            port: Service.cdr.port,
            method: Service.cdr.method,
            path: Service.cdr.path,
            headers: Service.cdr.headers,
            agent: Service.cdr.agent
        }
    },

    getFilesRequestParams: (params) => {
        if (!Service.files) return null;

        let path = Service.files.path,
            headers = Service.files.headers;
        for (let key in params)
            if (params.hasOwnProperty(key) && ~path.indexOf('${' + key + '}'))
                path = path.replace('${' + key + '}', params[key]);

        if (params.contentType)
            headers['content-type'] = params.contentType;

        if (params.contentLength)
            headers['content-length'] = params.contentLength;

        return {
            host: Service.files.host,
            port: Service.files.port,
            method: Service.files.method,
            path: path,
            headers: Service.files.headers,
            agent: Service.files.agent
        }
    }
};


function _sendCdr (doc, parentId, cb) {
    Service._sendRequest(Service.getCdrRequestParams(), JSON.stringify(doc), (res) => {
        if (res.statusCode !== 200 && res.statusCode !== 204) {
            if (parentId)
                Service._onError(new Error(`[cdr] Bad response ${res.statusCode}`), cdrCollectionName, parentId);

            return cb && cb(new CodeError(500, `Bad response code ${res.statusCode}`));
        }

        log.trace(`[cdr] Ok send ${parentId || 'bulk data'}`);
        return cb && cb(null, res);
    });
}

function _sendFile (doc, parentId, cb) {
    recordingsService._getFile(doc, {}, (err, res) => {
        if (err) {
            Service._onError(err, fileCollectionName, parentId);
            return cb && cb(err);
        }

        if (res && res.source && res.source.pipe) {
            const option = {
                uuid: doc.uuid,
                type: getTypeFromContentType(doc['content-type']),
                contentType: doc['content-type'],
                contentLength: doc.size,
                name: doc.name,
                domain: doc.domain
            };
            Service._sendStream(Service.getFilesRequestParams(option), res.source, (resDest) => {
                if (resDest.statusCode !== 200 && resDest.statusCode !== 204) {
                    Service._onError(new Error(`[file] Bad response ${resDest.statusCode}`), fileCollectionName, parentId);
                    return cb && cb(new CodeError(500, `Bad response code ${resDest.statusCode}`));
                }

                log.trace(`[file] Ok send ${parentId}`);
                return cb && cb(null);
            })
        } else {
            log.error(new Error(`Bad file stream`, res));
            return cb && cb(new CodeError(500, `Bad file stream`));
        }
    });
}

const getTypeFromContentType = (contentType = "") => {
    switch (contentType) {
        case 'application/pdf':
            return 'pdf';
        case 'audio/wav':
            return 'wav';
        case 'audio/mpeg':
        default:
            return 'mp3'
    }
};