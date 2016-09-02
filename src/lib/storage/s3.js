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
        this.mask = mask;
        this._bucketName = conf.bucketName;
        this._client = new aws.S3({
            signatureVersion: conf.signatureVersion || 'v4',
            accessKeyId: conf.accessKeyId,
            secretAccessKey: conf.secretAccessKey,
            region: conf.region
        });
    }

    get (fileDb, options, cb) {
        let params = {
            Bucket: fileDb.bucketName,
            Key: fileDb.path
        };

        if (options.dispositionName) {
            params.ResponseContentDisposition = 'attachment;  filename=' + options.dispositionName;
        }
        
        this._client.getSignedUrl('getObject', params, function (err, url) {
            log.debug(`try redirect to ${url}`);
            return cb(null, {
                location: url,
                statusCode: 302
            })
        });

    }

    save (fileConf, option, cb) {
        let s3Path = helper.getPath(this.mask, fileConf.domain, fileConf.name);
        let re = fs.createReadStream(fileConf.path);
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
            Bucket: fileConf.bucketName,
            Key: fileConf.path
        }, cb);
    }
};