/**
 * Created by igor on 23.08.16.
 */

"use strict";
    
const express = require('express'),
    compression = require('compression'),
    log = require(__appRoot + '/lib/log')(module),
    conf = require(`${__appRoot}/conf`),
    bodyParser = require('body-parser'),
    API = ['public', 'private', 'archive']
;

const UseCompression = `${conf.get('application:compression')}` === 'true';

module.exports = (application) => {
    const api = conf._getUseApi();
    const result = {};

    api.forEach( name => {
        result[`${name}Api`] = express();

        if (UseCompression)
            result[`${name}Api`].use(compression());

        if (name === 'public')
            result[`${name}Api`].use(express.static(__appRoot + '/public'));

        require('./logger').addRoutes(result[`${name}Api`]);
        require('./cors').addRoutes(result[`${name}Api`]);

        result[`${name}Api`].use(bodyParser.json({limit: '1000kb'}));

        require(`./${name}`)(result[`${name}Api`]);
        require('./error').addRoutes(result[`${name}Api`]);

        log.info(`--------- ${name} API ---------`);
        printRoutes(result[`${name}Api`]);
    });

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
    return result
};