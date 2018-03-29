/**
 * Created by igor on 29.08.18.
 */

"use strict";

const getElasticData = require('./helper').getElasticData;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */

function addRoutes(api) {
    api.post('/api/v2/statistic/:index', getStatistics);
}

function getStatistics(req, res, next) {
    const options = req.body;
    options.index = req.params.index;

    return getElasticData(req.webitelUser, options, res, next);
}

