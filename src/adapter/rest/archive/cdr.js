/**
 * Created by igor on 21.10.16.
 */

"use strict";

const log = require(`${__appRoot}/lib/log`)(module),
    cdrService = require(`${__appRoot}/services/cdr`)
;

module.exports = {
    addRoutes: api => {
        api.post('/api/v1/cdr', create);
    }
};

const create = (req, res, next) => {

    if (req.body instanceof Array) {
        return cdrService.saveToElasticBulk(req.body, (err, result) => {
            if (err) {
                log.error(err);
                return next(err);
            }

            log.debug(`Ok save bulk count ${req.body.length}`);
            res.json(result);
        });
    }

    const uuid = req.body.variables  && req.body.variables.uuid;
    cdrService.saveToElastic(req.body, (err) => {
        if (err) {
            log.error(err);
            return next(err);
        }

        log.debug(`Ok save: ${uuid}`);
        res.status(200).end();
    })
};