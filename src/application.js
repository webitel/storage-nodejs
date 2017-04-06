/**
 * Created by igor on 23.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    http = require('http'),
    https = require('https'),
    fs = require('fs'),
    conf = require(__appRoot + '/conf'),
    EventEmitter2 = require('eventemitter2').EventEmitter2
    ;
    
class Application extends EventEmitter2 {
    constructor () {
        super();
        this.DB = null;
        this.elastic = null;
        this.replica = null;

        require(__appRoot + '/lib/broker');

        let elasticConf = conf.get('elastic');

        if (elasticConf && elasticConf.enabled.toString() == 'true') {
            let ElasticClient = require('./db/elastic');
            this.elastic = new ElasticClient(elasticConf);
        }

        // TODO
        if (~conf._getUseApi().indexOf('public') || ~conf._getUseApi().indexOf('private')) {
            this.on('db:connect', (db) => {
                this.DB = db;
                if (this.elastic)
                    require('./services/cdr').processSaveToElastic();
            });
            this.once('db:connect', this.configureServer);

            this.on('db:error', this.reconnectDB.bind(this));

            process.nextTick(this.connectDB.bind(this));
        } else {
            log.info(`Skip connect to mongodb`);
            process.nextTick(this.configureServer.bind(this));
        }

        if (typeof gc == 'function') {
            setInterval(function () {
                gc();
                console.log('----------------- GC -----------------');
                console.dir(process.memoryUsage())
            }, 5000);
        }
    }

    reconnectDB () {
        setImmediate(this.connectDB.bind(this));
    }

    connectDB () {
        require('./db')(this);
    }

    configureServer () {
        if (`${conf.get('replica:enabled')}` === 'true') {
            this.replica = require('./services/replica');
            this.replica._init();
        }

        this.startServer(require('./adapter/rest')(this));
    }

    startServer (api) {
        let scope = this;

        const createApi = (api, params = {}, cb) => {
            if (`${params.enabled}` !== 'true')
                return;

            if (`${params.ssl}` === 'true') {
                const sslOptions = {
                    key: fs.readFileSync(params.ssl_key),
                    cert: fs.readFileSync(params.ssl_cert)
                };

                return https.createServer(sslOptions, api)
                    .listen(params.port, params.host, cb);
            } else {
                return http.createServer(api)
                    .listen(params.port, params.host, cb);
            }
        };

        for (let key in api) {
            if (api.hasOwnProperty(key)) {
                let name = key.replace(/Api$/, '');
                createApi(api[key], conf.get(name), function () {
                    log.info(`${name} server listening on port ${this.address().port}`);
                    scope.emit(`${name}:start`);
                })
            }
        }
    }

    stop () {
        // TODO close connections
    }
}

module.exports = new Application();