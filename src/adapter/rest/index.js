/**
 * Created by igor on 23.08.16.
 */

"use strict";
    
const express = require('express'),
    publicApi = express(),
    privateApi = express(),
    log = require(__appRoot + '/lib/log')(module),
    bodyParser = require('body-parser')
;

module.exports = (application) => {

    publicApi.use(express.static(__appRoot + '/public'));
    require('./logger').addRoutes(publicApi, privateApi);
    require('./cors').addRoutes(publicApi, privateApi);

    publicApi.use(bodyParser.json());
    privateApi.use(bodyParser.json());

// API
    require('./private')(privateApi);
    require('./public')(publicApi);


// Error handle
    require('./error').addRoutes(publicApi, privateApi);

    log.info(`--------- Public API ---------`);
    printRoutes(publicApi);
    log.info(`--------- Private API ---------`);
    printRoutes(privateApi);

    function printRoutes(api) {
        let route, methods;
        api._router.stack.forEach(function(middleware){
            if(middleware.route){ // routes registered directly on the app
                route = middleware.route;
                methods = Object.keys(route.methods);
                log.info('Add: [%s]: %s', (methods.length > 1) ? 'ALL' : methods[0].toUpperCase(), route.path);

            } else if(middleware.name === 'router'){ // router middleware
                middleware.handle.stack.forEach(function(handler){
                    route = handler.route;
                    methods = Object.keys(route.methods);
                    log.info('Add: [%s]: %s', (methods.length > 1) ? 'ALL' : methods[0].toUpperCase(), route.path);
                });
            }
        });
    }
    return {
        publicApi: publicApi,
        privateApi: privateApi
    }
};