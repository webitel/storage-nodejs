/**
 * Created by igor on 29.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    CodeError = require(__appRoot + '/lib/error'),
    helper = require('../helper'),
    B2 = require('./b2'),
    async = require('async'),
    TYPE_ID = 2;

module.exports = class B2Storage {

    constructor (conf, mask) {
        this.accountId = conf.accountId;
        this.bucketId = conf.bucketId;
        this.bucketName = conf.bucketName;
        this.applicationKey = conf.applicationKey;
        this.mask = mask || "$Y/$M/$D/$H";
        this._expireToken = 0;
        this.name = "b2";
        this._uploadQueue = async.queue( (fileConf, cb) => {
            B2.saveFile(this._authParams, fileConf, helper.getPath(this.mask, fileConf.domain, fileConf.name), cb);
        });
        this._uploadQueue.drain = () => {
            log.debug('All upload files done');
        };
    }

    _isAuth (cb) {
        if (this._expireToken <= Date.now()) {
            return this._auth(cb);
        } else {
            return cb(null)
        }
    }

    _auth (cb) {
        let conf = {
            accountId: this.accountId,
            applicationKey: this.applicationKey,
            bucketId: this.bucketId,
            bucketName: this.bucketName
        };

        B2.auth(conf, (err, authParam) => {
            // token validate 24h
            // storage re-auth 6h
            // 1000 * 60 * 60 * 6 = 21600000
            if (err) {
                log.error(err);
                return cb && cb(err);
            }

            this._expireToken = Date.now() + 21600000;
            log.trace(`Set new auth ${conf.accountId}`);
            this._authParams = authParam;
            return cb && cb();
        });
    }

    get (fileDb, options, cb) {
        this._isAuth((err) => {
            if (err)
                return cb(err);

            return B2.getFile(this._authParams, fileDb, options.range, cb);
        })
    }

    save (fileConf, option, cb) {
        this._isAuth((err) => {
            if (err)
                return cb(err);

            return this._uploadQueue.push(fileConf, cb)
        });
    }

    del (fileConf, cb) {
        this._isAuth((err) => {
            if (err)
                return cb(err);

            return B2.delFile(this._authParams, fileConf, cb);
        });
    }

    validateConfig (config) {
        console.log('validate b2');
        return !config || this.accountId != config.accountId || this.bucketId != config.bucketId
            || this.bucketName != config.bucketName || this.applicationKey != config.applicationKey

    }
};