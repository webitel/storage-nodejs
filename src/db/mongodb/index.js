/**
 * Created by igor on 23.08.16.
 */

"use strict";

const MongoClient = require("mongodb").MongoClient,
    mongoClient = new MongoClient(),
    config = require(__appRoot + '/config'),
    log = require(__appRoot + '/lib/log')(module)
    ;


module.exports = initConnect;

function initConnect (server) {
    var options = {
        server: {
            auto_reconnect: true,
            slaveOk: true,
            replset: {
                socketOptions: {
                    keepAlive: 1,
                    connectTimeoutMS : 30000 ,
                    socketTimeoutMS: 90000
                }
            }
        }
    };
    mongoClient.connect(config.get('mongodb:uri'), options, function(err, db) {
        if (err) {
            log.error('Connect db error: %s', err.message);
            return server.emit('db:error', err);
        }
        log.info('Connected db %s ', config.get('mongodb:uri'));

        db._query = {
            auth: require('./query/auth').addQuery(db),
            file: require('./query/file').addQuery(db),
            domain: require('./query/domain').addQuery(db),
            email: require('./query/email').addQuery(db),
            cdr: require('./query/cdr').addQuery(db),
            media: require('./query/media').addQuery(db)
        };

        server.emit('db:connect', db);

        db.on('close', function () {
            log.warn('close MongoDB');
            server.emit('sys::closeDb', db);
        });
        db.on('reconnect', function () {
            log.info('Reconnect MongoDB');
            server.emit('sys::reconnectDb', db);
        });

        db.on('error', function (err) {
            log.error('close MongoDB: ', err);
        });
    });
}