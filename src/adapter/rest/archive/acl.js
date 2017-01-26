/**
 * Created by igor on 21.10.16.
 */

"use strict";

const CodeError = require(`${__appRoot}/lib/error`),
    conf = require(`${__appRoot}/conf`),
    BASE_LOGIN = conf.get('archive:login'),
    BASE_PASSWORD = conf.get('archive:password'),
    BASE_TOKEN = new Buffer(`${BASE_LOGIN}:${BASE_PASSWORD}`).toString('base64')
    ;

module.exports = {
    addRoutes: api => {
        api.all('/api/v1/*', check)
    }
};

const check = (req, res, next) => {
    const header = req.headers['authorization'] || '',
        token = header.split(/\s+/).pop() || '';

    if (BASE_TOKEN !== token) {
        return next(new CodeError(401, 'Invalid credentials'));
    }
    // TODO
    req.webitelUser = {
        role: "GOD",
        acl: {
            'cdr/files': ['*'],
            'cdr/media': ['*'],
            'cdr': ['*']
        },
        domain: null
    };
    return next();
};