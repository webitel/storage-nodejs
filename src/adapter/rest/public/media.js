/**
 * Created by igor on 31.08.16.
 */

"use strict";

const mediaService = require(__appRoot + '/services/media'),
    streaming = require(__appRoot + '/utils/http').streaming,
    checkPermission = require(__appRoot + '/utils/acl'),
    TTSMiddleware = require(__appRoot + '/services/tts'),
    log = require(__appRoot + '/lib/log')(module)
    ;


module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.get('/api/v2/media/tts/:provider', generateFileFromTTS);
    api.post('/api/v2/media/:type', saveFile);
    api.get('/api/v2/media', listMedia);
    api.get('/api/v2/media/:type/:name', getFile);
    api.delete('/api/v2/media/:type/:name', deleteFile);
}

function generateFileFromTTS(req, res, next) {
    if (!checkPermission(req.webitelUser.acl, 'cdr/media', 'c'))
        return next(new CodeError(403, 'Forbidden!'));

    return TTSMiddleware(req, res)
}

function saveFile(req, res, next) {
    let type = req.params.type || 'mp3',
        domainName = req.query.domain
        ;


    mediaService.saveFilesFromRequest(
        req.webitelUser,
        req,
        {domain: domainName, type: type},
        (err) => {
            if (err)
                return next(err);

            return res
                .status(201)
                .end()
        }
    );
}

function listMedia(req, res, next) {
    let options = {
        domain: req.query.domain
    };

    return mediaService.list(req.webitelUser, options, (err, mediaArray) => {
        if (err)
            return next(err);

        return res
            .status(200)
            .json({
                "status": "OK",
                "data": mediaArray
            })
    })
}

function getFile(req, res, next) {
    let dispositionName = req.query.file_name;
    
    let options = {
        domain: req.query.domain,
        type: req.params.type,
        name: req.params.name,
        range: req.headers['range'],
        dispositionName: dispositionName
    };

    mediaService.get(req.webitelUser, options, (err, response) => {
        if (err) {
            return next(err);
        }

        if (!response || !response.source)
            return next(new Error(`No source stream.`));

        return streaming(response.source, res, {
            dispositionName: dispositionName,
            range: options.range,
            contentType: response.contentType,
            totalLength: response.totalLength
        });
    });
}

function deleteFile(req, res, next) {
    const options = {
        name: req.params.name,
        type: req.params.type,
        domain: req.query.domain
    };

    mediaService.delete(req.webitelUser, options, (err, result) => {
        if (err)
            return next(err);

        if (!result)
            return res
                .status(404)
                .json({
                    "status": "error",
                    "info": `File not found.`
                });

        return res
            .status(200)
            .json({
                "status": "OK",
                "info": result
            })
    })
}