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
    api.post('/api/v2/cdr/text/scroll', scrollElasticData);
    api.put('/api/v2/cdr/:uuid/pinned', addPin);
    api.delete('/api/v2/cdr/:uuid/pinned', removePin);

    api.post('/api/v2/cdr/searches', searches);
    api.post('/api/v2/cdr/counts', count);

    api.delete('/api/v2/cdr/:uuid', remove);

    api.post('/api/v2/cdr/:uuid/post', savePostProcess);
}

function savePostProcess(req, res, next) {
    const uuid = req.params.uuid;
    const data = {
        _id: uuid,
        variables: {uuid, domain_name: req.webitelUser.domain},
        post_data: req.body
    };

    application.elastic.insertPostProcess(data, (err) => {
        if (err) {
            log.error(err);
            return next(err);
        }

        log.debug(`Ok save: ${uuid}`);
        res.status(200).end();
    });
}

function addPin(req, res, next) {
    const option = {
        id: req.params.uuid,
        index: req.query.index
    };

    cdrService.addPin(req.webitelUser, option, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}

function removePin(req, res, next) {
    const option = {
        id: req.params.uuid,
        index: req.query.index
    };

    cdrService.delPin(req.webitelUser, option, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}

function getElasticData(req, res, next) {
    return elasticService.search(req.webitelUser, req.body, (err, result) => {
        if (err)
            return next(err);

        res.json(result);
    });
}

function scrollElasticData(req, res, next) {
    return elasticService.scroll(req.webitelUser, req.body, (err, result) => {
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

function remove(req, res, next) {
    let option = {
        uuid: req.params.uuid
    };

    cdrService.remove(req.webitelUser, option, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}