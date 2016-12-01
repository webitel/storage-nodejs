/**
 * Created by igor on 21.10.16.
 */

"use strict";

const fileService = require(`${__appRoot}/services/file`),
    recordingsService = require(`${__appRoot}/services/recordings`),
    streaming = require(`${__appRoot}/utils/http`).streaming,
    log = require(`${__appRoot}/lib/log`)(module)
    ;

module.exports = {
    addRoutes: api => {
        api.post('/api/v1/recordings', create);
        api.get('/api/v1/recordings/:id', getResource);
    }
};

const getResource = (req, res, next) => {
    const dispositionName = req.query.file_name;
    const params = {
        domain: req.query.domain,
        dispositionName: dispositionName,
        hash: req.query.hash
    };

    if (req.headers.range)
        params.range = req.headers.range;

    const uuid = req.params.id;

    recordingsService.getFileFromHash(req.webitelUser, uuid, params, (err, response) => {
        if (err)
            return next(err);
        
        if (!response || !response.source)
            return next(`No source stream.`);

        return streaming(response.source, res, {
            range: params.range,
            dispositionName: dispositionName,
            totalLength: response.totalLength,
            contentType: response.contentType
        });
    })
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
            if (err) {
                log.error(err);
                fileService.deleteFile(file.path, (err) => {
                    if (err)
                        log.error(err);
                });

                return next(err);
            }

            const doc = recordingsService.getSchema(file, response);

            recordingsService.saveInfoToElastic(doc, (err) => {

                fileService.deleteFile(file.path, (err) => {
                    if (err)
                        log.error(err);
                });

                if (err)
                    return next(err);

                return res.status(200).end();

            });

        });
    });
};