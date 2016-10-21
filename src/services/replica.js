/**
 * Created by igor on 21.10.16.
 */

"use strict";

const log = require(`${__appRoot}/lib/log`)(module),
    conf = require(`${__appRoot}/config`),
    CodeError = require(`${__appRoot}/lib/error`),
    url = require('url'),
    cdrConf = conf.get('replica:cdr'),
    filesConf = conf.get('replica:files'),
    hostConf = conf.get('replica:host'),
    authConf = conf.get('replica:auth'),
    live = `${conf.get('replica:live')}` === 'true',
    headersConf = conf.get('replica:headers'),
    cdrCollectionName = conf.get('mongodb:collectionCDR'),
    fileCollectionName = conf.get('mongodb:collectionFile')
;

const Service = module.exports = {
    cdr: {},
    files: {},
    provider: null,

    sendCdr: (doc, parentId) => {
        try {
            if (!parentId)
                return log.error(`[cdr] Bad parent id`, doc);

            if (live) {
                Service._sendRequest(Service.getCdrRequestParams(), JSON.stringify(doc), (res) => {
                    if (res.statusCode !== 200)
                        return Service._onError(new Error(`[cdr] Bad response ${res.statusCode}`), cdrCollectionName, parentId);

                    log.trace(`[cdr] Ok send ${parentId}`);
                });
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
                Service._sendRequest(Service.getFilesRequestParams(), JSON.stringify(doc), (res) => {
                    if (res.statusCode !== 200)
                        return Service._onError(new Error(`[file] Bad response ${res.statusCode}`), fileCollectionName, parentId);

                    log.trace(`[file] Ok send ${parentId}`);
                });
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

    _sendRequest: (params = {}, body, cb) => {
        if (!Service.provider)
            return cb(new CodeError(500, `No initialize replica provider`));

        const req = Service.provider.request(params, cb);
        req.on('error', cb);
        if (body)
            req.write(body);
        req.end();
    },

    setAuth: req => {

    },

    getRequestOption: () => {
        return {

        }
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

    getFilesRequestParams: () => {
        if (!Service.files) return null;

        return {
            host: Service.files.host,
            port: Service.files.port,
            method: Service.files.method,
            path: Service.files.path,
            headers: Service.files.headers,
            agent: Service.files.agent
        }
    }
};