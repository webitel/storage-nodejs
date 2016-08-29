/**
 * Created by igor on 25.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    recordingsService = require(__appRoot + '/services/recordings'),
    fileService = require(__appRoot + '/services/file')
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.put('/sys/formLoadFile?:id', saveFile); // deprecated
    api.put('/sys/records', saveFile);
}

function saveFile(req, res, next) {
    let uuid = req.query.id,
        name = req.query.name || 'tmp';

    fileService.requestStreamToCache( `${uuid}_${name}.${req.query.type || 'mp3'}`, req, (err, file) => {
        if (err)
            return next(err);

        file.domain = req.query.domain;
        file.uuid = req.query.id;
        file.queryName = name;

        recordingsService.saveFile(file, req.query, (err, result) => {
            fileService.deleteFile(file.path, (err) => {
                if (err)
                    log.error(err);
            });

            if (err)
                return next(err);
            
            res.status(204).end();

        })
    });
}