/**
 * Created by igor on 31.08.16.
 */

"use strict";


const CodeError = require(__appRoot + '/lib/error'),
    mediaService = require(__appRoot + '/services/media'),
    log = require(__appRoot + '/lib/log')(module)
    ;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.get('/sys/media/:type/:id', getFile)
}

function getFile(req, res, next) {
    res.status(200).end();
}