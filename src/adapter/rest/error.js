/**
 * Created by Admin on 04.08.2015.
 */

'use strict';

var log = require(__appRoot + '/lib/log')(module);

module.exports = {
    addRoutes: addRoutes
};

function addRoutes() {
    for (let api of arguments) {
        api.use(function(req, res){
            res.status(404);
            res.json({
                "status": "error",
                "info": req.originalUrl + ' not found.'
            });
        });

        api.use(function(err, req, res, next){
            log.warn('Api status: %s, message: %s', err.status, err.message);
            res.status(err.status || 500);
            res.json({
                "status": "error",
                "info": err.message,
                "code": err.code
            });
        });
    }
}