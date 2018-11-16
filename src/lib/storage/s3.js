/**
 * Created by igor on 02.09.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    aws = require('aws-sdk'),
    fs  = require('fs'),
    helper = require('./helper'),
    TYPE_ID = 1;

module.exports = class S3Storage {
    constructor (conf = {}, mask) {
        this.name = 's3';
        this.mask = mask;
        this._conf = conf;
        this._bucketName = conf.bucketName;
        this._client = new aws.S3({
            signatureVersion: conf.signatureVersion || 'v4',
            accessKeyId: conf.accessKeyId,
            secretAccessKey: conf.secretAccessKey,
            region: conf.region
        });

    }

    checkConfig (conf = {}, mask) {
        return this.mask == mask && this._conf.accessKeyId == conf.accessKeyId && this._conf.secretAccessKey == conf.secretAccessKey
            && this._bucketName == conf.bucketName
    }

    get (fileDb, options, cb) {
        let params = {
            Bucket: this._bucketName,
            Key: fileDb.path
        };

        if (options.dispositionName) {
            params.ResponseContentDisposition = 'attachment;  filename=' + options.dispositionName;
        }

        if (options.stream) {
            this._client.getObject(params)
                .on('httpHeaders', function (statusCode) {
                    if (statusCode !== 200) {
                        return cb(new Error(`AWS response ${statusCode}`));
                    }

                    return cb(null, this.response.httpResponse.createUnbufferedStream());
                })
                .send();
        } else {
            this._client.getSignedUrl('getObject', params, function (err, url) {
                if (err)
                    return cb(err);

                log.debug(`try redirect to ${url}`);
                return cb(null, {
                    location: url,
                    statusCode: 302
                })
            });
        }

    }

    copyTo (fileDb, to, cb) {
        this.get(fileDb, {stream: true}, (err, stream) => {
            if (err)
                return cb(err);

            to.save(fileDb, {stream}, cb);
        });
    }


    save (fileConf, option = {}, cb) {
        let s3Path = helper.getPath(this.mask, fileConf.domain, fileConf.name);
        let re = option.stream || fs.createReadStream(fileConf.path);
        re.once('error', (err) => {
            log.error(err);
            return cb(err);
        });
        re.on('open', () => {
            this
                ._client
                .putObject(
                    {
                        Bucket: this._bucketName,
                        Key: s3Path,
                        Body: re,
                        ContentType: fileConf.contentType || 'audio/mpeg'
                    },
                    (err) => {
                        if (err)
                            return cb(err);

                        log.trace(`Save (${this._bucketName}) file: ${s3Path}`);
                        return cb(null, {
                            path: s3Path,
                            type: TYPE_ID,
                            bucketName: this._bucketName
                        })
                    }
                );

        });
    }

    del (fileConf, cb) {
        this._client.deleteObject({
            Bucket: this._bucketName,
            Key: fileConf.path
        }, cb);
    }

    existsFile (fileConf, cb) {
        let params = {
            Bucket: this._bucketName,
            Key: fileConf.path
        };

        this._client.headObject(params, (err, data) => {
            if (err && err.statusCode !== 404)
                return cb(err);

            return cb(null, !!data)
        })
    }
};