/**
 * Created by igor on 31.08.16.
 */

"use strict";


const mediaService = require(__appRoot + '/services/media'),
    streaming = require(__appRoot + '/utils/http').streaming,
    log = require(__appRoot + '/lib/log')(module)
    ;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.get('/sys/media/:type/:name', getFile)
}

function getFile(req, res, next) {

    let options = {
        domain: req.query.domain,
        type: req.params.type,
        name: req.params.name,
        range: req.headers['range']
    };

    mediaService._get(req.query.domain, options, (err, response) => {
        if (err) {
            return next(err);
        }

        if (!response || !response.source)
            return next(`No source stream.`);

        return streaming(response.source, res, {
            range: response.range,
            contentType: response.contentType,
            totalLength: response.totalLength
        });
    });
}