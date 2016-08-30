/**
 * Created by igor on 23.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    http = require('http'),
    https = require('https'),
    fs = require('fs'),
    conf = require(__appRoot + '/config'),
    EventEmitter2 = require('eventemitter2').EventEmitter2
    ;
    
class Application extends EventEmitter2 {
    constructor () {
        super();
        this.DB = null;
        this.elastic = null;

        this.once('db:connect', this.configureServer);
        this.on('db:connect', (db) => {
            this.DB = db;
        });

        this.on('db:error', this.reconnectDB.bind(this));

        if (typeof gc == 'function') {
            setInterval(function () {
                gc();
                console.log('----------------- GC -----------------');
            }, 5000);
        }

        let elasticConf = conf.get('elastic');

        if (elasticConf && elasticConf.enabled.toString() == 'true') {
            let ElasticClient = require('./db/elastic');
            this.elastic = new ElasticClient(elasticConf);
        }
        
        process.nextTick(this.connectDB.bind(this));
    }

    reconnectDB () {
        setImmediate(this.connectDB.bind(this));
    }

    connectDB () {
        require('./db')(this);
    }

    configureServer () {
        this.startServer(require('./adapter/rest')(this));
    }

    startServer (api) {
        let scope = this,
            publicSrv,
            privateSrv;
        
        function onStartPublic() {
            log.info(`Public server listening on port ${publicSrv.address().port}`);
            scope.emit('public:start');
        }
        function onStartPrivate() {
            log.info(`Private server listening on port ${privateSrv.address().port}`);
            scope.emit('public:start');
        }

        if (conf.get('public:enabled').toString() == 'true') {
            let https_options = {
                key: fs.readFileSync(conf.get('public:ssl_key')),
                cert: fs.readFileSync(conf.get('public:ssl_cert'))
            };
            publicSrv = https.createServer(https_options, api.publicApi)
                .listen(
                    conf.get('public:port'),
                    conf.get('public:host'),
                    onStartPublic
                );
        } else {
            publicSrv = http.createServer(api.publicApi).listen(
                conf.get('public:port'),
                conf.get('public:host'),
                onStartPublic
            );
        }

        privateSrv = http.createServer(api.privateApi).listen(
            conf.get('private:port'),
            conf.get('private:host'),
            onStartPrivate
        );
    }

    stop () {
        // TODO close connections
    }
}

module.exports = new Application();