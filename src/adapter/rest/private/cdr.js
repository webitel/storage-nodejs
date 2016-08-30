/**
 * Created by igor on 26.08.16.
 */

"use strict";

const CodeError = require(__appRoot + '/lib/error'),
    cdrService = require(__appRoot + '/services/cdr'),
    log = require(__appRoot + '/lib/log')(module)
    ;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.post('/sys/cdr', save)
}

function save(req, res, next) {
    const uuid = req.query.uuid;
    cdrService.save(req.body, (err) => {
        if (err) {
            log.error(err);
            return next(err);
        }

        log.debug(`Ok save: ${uuid}`);
        res.status(200).end();
    })
}