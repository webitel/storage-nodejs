/**
 * Created by igor on 21.10.16.
 */

"use strict";

const fileService = require(`${__appRoot}/services/file`),
    recordingsService = require(`${__appRoot}/services/recordings`),
    log = require(`${__appRoot}/lib/log`)(module)
    ;

module.exports = {
    addRoutes: api => {
        api.post('/api/v1/recordings', create);
        api.get('/api/v1/recordings/:id', create);
    }
};

const create = (req, res, next) => {
    let uuid = req.query.id,
        name = req.query.name || 'tmp',
        type = req.query.type,
        domainName = req.query.domain
        ;

    fileService.requestStreamToCache( `${uuid}_${name}.${type || 'mp3'}`, req, (err, file) => {
        if (err)
            return next(err);

        file.domain = domainName;
        file.uuid = uuid;
        file.queryName = name;

        recordingsService.saveToLocalProvider(file, req.query, (err, response = {}) => {

            const doc = recordingsService.getSchema(file, response);

            recordingsService.saveInfoToElastic(doc, (err) => {

                fileService.deleteFile(file.path, (err) => {
                    if (err)
                        log.error(err);
                });

                if (err)
                    return next(err);

                return res.status(204).end();

            });

        });
    });
};