/**
 * Created by igor on 23.08.16.
 */

"use strict";

const authService = require(__appRoot + '/services/auth'),
    jwt = require('jwt-simple'),
    config = require(__appRoot + '/config'),
    CodeError = require(__appRoot + '/lib/error'),
    tokenSecretKey = config.get('application:auth:tokenSecretKey')
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.all('/api/v2/*', checkAllow);
}

function checkAllow(req, res, next) {
    const token = (req.body && req.body.access_token) || (req.query && req.query.access_token) || req.headers['x-access-token'],
        key = (req.body && req.body.x_key) || (req.query && req.query.x_key) || req.headers['x-key'];

    if (token && key) {
        try {
            var decoded = jwt.decode(token, tokenSecretKey);
        } catch (e) {
            return next(new CodeError(401, "Invalid Token or Key"));
        }

        if (decoded.exp <= Date.now()) {
            return next(new CodeError(401, "Token Expired"));
        }

        // Authorize the user to see if s/he can access our resources

        authService.getUserByKey(key, function (err, dbUser) {
            if (dbUser && dbUser.token == token) {
                req.webitelUser = {
                    id: dbUser.username,
                    domain: dbUser.domain,
                    role: dbUser.role,
                    roleName: dbUser.roleName,
                    expires: dbUser.expires,
                    acl: dbUser.acl
                    //testLeak: new Array(1e6).join('X')
                };
                next(); // To move to next middleware
            } else {
                // No user with this name exists, respond back with a 401
                return next(new CodeError(401, "Invalid User"));
            }
        });
    } else {
        return next(new CodeError(401, "Invalid Token or Key"));
    }
}