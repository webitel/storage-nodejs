/**
 * Created by igor on 23.08.16.
 */

"use strict";

const getClientIp = require('request-ip').getClientIp,
    log = require(__appRoot + '/lib/log')(module),
    authService = require(__appRoot + '/services/auth')
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.all('/sys/*', checkAllow);
}

function checkAllow(req, res, next) {
    let ip = getClientIp(req) || '',
        err = authService.checkBannedIp(ip);

    if (err) {
        return next(err);
    }
    req.webitelUser = {
        role: "GOD",
        attr: {
            acl: {
                'cdr/files': ['*'],
                'cdr/media': ['*'],
                'cdr': ['*']
            }
        }
    };
    return next();
}