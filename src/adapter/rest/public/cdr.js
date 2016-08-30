/**
 * Created by igor on 30.08.16.
 */

"use strict";

const elasticService = require(__appRoot + '/services/elastic'),
    cdrService = require(__appRoot + '/services/cdr'),
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
    api.post('/api/v2/cdr/searches', searches);
    api.post('/api/v2/cdr/counts', count);
}

function getElasticData(req, res, next) {
    return elasticService.search(req.webitelUser, req.body, (err, result) => {
        if (err)
            return next(err);

        res.json(result);
    });
}

function searches(req, res, next) {

    let option = {
        columns: req.body.columns,
        filter: req.body.filter,
        limit: req.body.limit,
        pageNumber: req.body.pageNumber,
        sort: req.body.sort
    };

    cdrService.search(req.webitelUser, option, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}
function count(req, res, next) {

    let option = {
        filter: req.body.filter
    };

    cdrService.count(req.webitelUser, option, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}