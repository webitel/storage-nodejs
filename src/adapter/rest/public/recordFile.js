/**
 * Created by igor on 25.08.16.
 */

"use strict";

const fileService = require(__appRoot + '/services/recordings'),
    streaming = require(__appRoot + '/utils/http').streaming,
    log = require(__appRoot + '/lib/log')(module)
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.get('/api/v2/files/stats/:id?', stats);
    api.get('/api/v2/files/:id', getFile);
    api.delete('/api/v2/files/:id', delFile);
    // api.delete('/api/v2/files/cache/:domain', flushCache);
    api.delete('/api/v2/files/utils/removeNonExistentFiles', removeNonExistentFiles);
}

const stats = (req, res, next) => {
    let option = req.query;
    option.uuid = req.params.id;

    fileService.stats(req.webitelUser, option, (err, result) => {
        if (err)
            return next(err);

        let _size = result instanceof Array ? result[0] : result;

        if (!_size) {
            _size = {
                "size": 0
            }
        }

        res
            .status(200)
            .json(_size);
    });
};

const getFile = (req, res, next) => {
    let uuid = req.params.id,
        contentType = req.query.type || 'audio/mpeg',
        pathName = req.query.name,
        dispositionName = req.query['file_name']
        ;

    let params = {
        contentType: contentType,
        pathName: pathName,
        range: req.headers['range'],
        dispositionName: dispositionName
    };

    fileService.getFileFromUUID(
        req.webitelUser,
        uuid,
        params,
        (err, response) => {
            if (err)
                return next(err);

            if (contentType === 'all')
                return res
                    .status(200)
                    .json(response)
                    ;

            if (!response || !response.source)
                return next(`No source stream.`);

            return streaming(response.source, res, {
                dispositionName: dispositionName,
                range: params.range,
                contentType: contentType,
                totalLength: response.totalLength
            });
        }
    );
};

const delFile = (req, res, next) => {
    let options = {
        uuid: req.params.id,
        delDb: req.query.db == 'true',
        pathName: req.query.name,
        domain: req.query.domain
    };

    fileService.delFile(req.webitelUser, options, (err, fileDb) => {
        if (err)
            return next(err);

        return res
            .status(200)
            .json(fileDb);
    })
};

const removeNonExistentFiles = (req, res, next) => {
    let options = {
        from: req.body.from,
        to: req.body.to
    };

    fileService.removeNonExistentFiles(req.webitelUser, options, (err, result) => {
        if (err)
            return next(err);

        return res
            .status(200)
            .json({
                "status": "OK",
                "info": result
            })
    })
};

const flushCache = (req, res, next) => {
    
};