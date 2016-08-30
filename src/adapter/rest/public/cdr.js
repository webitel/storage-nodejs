/**
 * Created by igor on 30.08.16.
 */

"use strict";

const elasticService = require(__appRoot + '/services/elastic'),
    log = require(__appRoot + '/lib/log')(module)
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.post('/api/v2/cdr/text', getElasticData);
}

function getElasticData(req, res, next) {
    return elasticService.search(req.webitelUser, req.body, (err, result) => {
        if (err)
            return next(err);

        res.json(result);
    });
}