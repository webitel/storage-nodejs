/**
 * Created by I. Navrotskyj on 04.10.17.
 */

"use strict";

const fileService = require(__appRoot + '/services/file');
const recordingsService = require(`${__appRoot}/services/recordings`);
const log = require(`${__appRoot}/lib/log`)(module);

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.put('/sys/tcp_dump', saveDump);
    api.put('/sys/tcp_dump/:id/error', setError);
}

function setError(req, res, next) {
    let error = '';

    req.on('data', c => error += c);
    req.on('end', () => {
        application.PG.getQuery('tcpDump').setFile(req.params.id, {error: error || 'unknown'}, (err) => {
            if (err)
                return next(err);
        });

        res.end();
    });
}

function saveDump(req, res, next) {
    const name = `${req.query.id || Date.now()}.${req.query.type || 'pcap'}`;
    fileService.requestStreamToCache(name, req, (err, file) => {
        if (err)
            return next(err);

        file.domain = 'tcp_dump';
        file.uuid = name;
        file.queryName = name;

        recordingsService.saveToLocalProvider(file, {}, (err, response = {}) => {
            fileService.deleteFile(file.path, (err) => {
                if (err)
                    log.error(err);
            });

            if (err) {
                log.error(err);
                return next(err)
            }

            const doc = recordingsService.getSchema(file, response);

            application.PG.getQuery('tcpDump').setFile(req.query.id, doc, (err) => {
                if (err)
                    return next(err);

                res.end();
            });
        })
    })
}