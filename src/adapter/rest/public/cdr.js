/**
 * Created by igor on 30.08.16.
 */

"use strict";

const elasticService = require(__appRoot + '/services/elastic'),
    cdrService = require(__appRoot + '/services/cdr'),
    log = require(__appRoot + '/lib/log')(module),
    getElasticData = require('./helper').getElasticData
    ;

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {

    api.post('/api/v2/cdr/text', getElasticCDRData);
    api.post('/api/v2/cdr/text/scroll', scrollElasticData);

    api.put(   '/api/v2/cdr/:uuid/pinned', addPin);
    api.delete('/api/v2/cdr/:uuid/pinned', removePin);

    api.get('/api/v2/cdr/:uuid', getByUuid); // ?leg=ba //
    api.get('/api/v2/cdr/:uuid/b', getLegsByAUuid); // all call

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

function getElasticCDRData(req, res, next) {
    const options = req.body;

    switch (req.query.leg) {
        case "b":
            options.index = 'cdr-b';
            break;
        case "*":
            options.index = 'cdr-*';
            break;
        default:
            options.index = 'cdr-a'
    }

    return getElasticData(req.webitelUser, options, res, next);
}

function scrollElasticData(req, res, next) {
    return elasticService.scroll(req.webitelUser, req.body, (err, result) => {
        if (err)
            return next(err);

        res.json(result);
    });
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

function getByUuid(req, res, next) {
    const options = {
        uuid: req.params.uuid,
        leg: req.query.leg,
        domain: req.query.domain
    };

    cdrService.getItem(req.webitelUser, options, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}

function getLegsByAUuid(req, res, next) {
    const options = {
        uuid: req.params.uuid,
        domain: req.query.domain
    };

    cdrService.getLegsByAUuid(req.webitelUser, options, (err, data) => {
        if (err)
            return next(err);

        res.json(data);
    })
}