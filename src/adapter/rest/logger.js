/**
 * Created by igor on 23.08.16.
 */

"use strict";


const log = require(__appRoot + '/lib/log')(module),
    getIp = require('request-ip').getClientIp;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes() {
    for (let api of arguments) {
        api.use(function(req, res, next) {
            log.trace(`Method: ${req.method}, url: ${req.url}, path: ${req.path}, ip: ${getIp(req)}`);
            next();
        });
    }
}