/**
 * Created by igor on 12.09.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    fs = require('fs'),
    helper = require('./helper'),
    https = require('https'),
    CodeError = require(__appRoot + '/lib/error'),
    TYPE_ID = 4
    ;

module.exports = class DropBoxStorage {
    constructor (conf, mask) {
        this.name = 'DropBox';
        this.mask = mask;
        this._token = conf.accessToken;
    }

    /**
     *
     * @param fileDb
     * @param options
     * @param cb
     */
    get (fileDb, options, cb) {
        let requestParams = {
            method: 'GET',
            path: '/2/files/download',
            host: 'content.dropboxapi.com',
            headers: {
                'Authorization': 'Bearer ' + this._token,
                'Dropbox-API-Arg': getJson({path: fileDb.path})
            }
        };
        let range = options && options.range;
        if (range && range.Start < range.End) {
            requestParams.headers['Range'] = `bytes=${range.Start}-${range.End}`;
        }

        let request = https.request(requestParams, (res) => {
            log.trace(`Response code get file: ${res.statusCode}`);

            if (res.statusCode !== 200 && res.statusCode !== 206)
                return cb(new CodeError(res.statusCode, 'Auth error.'));

            cb(null, res);
        });
        request.on('error', (e) => log.error(e));
        request.end();

    }

    copyTo (fileDb, to, cb) {
        this.get(fileDb, {}, (err, stream) => {
            if (err)
                return cb(err);

            to.save(fileDb, {stream}, cb);
        });
    }

    /**
     *
     * @param domain
     * @param fileName
     */
    getFilePath(domain, fileName) {

    }

    /**
     *
     * @param fileConf
     * @param option
     * @param cb
     */
    save (fileConf, option = {}, cb) {
        const fileName =  '/' + helper.getPath(this.mask, fileConf.domain, fileConf.name);

        const requestParams = {
            method: 'POST',
            path: '/2/files/upload',
            host: 'content.dropboxapi.com',
            headers: {
                'Authorization': 'Bearer ' + this._token,
                'Content-Type': "application/octet-stream",
                'Dropbox-API-Arg': getJson({
                    path: fileName,
                    mode: "add",
                    autorename: true,
                    mute: false
                }),
                'Content-Length': fileConf.contentLength
            }
        };

        const request = https.request(requestParams, (res) => {
            log.trace(`DropBox upload status code `, res.statusCode);

            let data = '';

            res.once('error', cb);
            res.on('data', function(chunk) {
                data += chunk;
            });
            res.on('end', function() {
                try {
                    if (res.statusCode !== 200) {
                        return cb(new CodeError(res.statusCode, data || "Internal error."));
                    }

                    log.trace(`Save storage file path: ${fileName}`);
                    let json = JSON.parse(data);
                    return cb(null, {
                        path: json.path_lower,
                        type: TYPE_ID
                    });
                } catch (e) {
                    return cb (e)
                }
            });


        });

        request.on('error', (e) => log.error(e));

        let rd = option.stream || fs.createReadStream(fileConf.path);

        rd.on("error", function(e) {
            log.error(e);
            return cb(e);
        });

        rd.on('open', () => rd.pipe(request));
    }

    /**
     *
     * @param fileConf
     * @param cb
     */
    del (fileDb, cb) {
        let requestParams = {
            method: 'POST',
            path: '/2/files/delete',
            host: 'api.dropboxapi.com',
            headers: {
                'Authorization': 'Bearer ' + this._token,
                'Content-Type': 'application/json'
            }
        };


        log.debug(`try delete file: ${fileDb.uuid}`);

        let request = https.request(requestParams, (res) => {
            log.trace(`Response code delete file: ${res.statusCode}`);

            if (res.statusCode !== 200)
                return cb(new CodeError(res.statusCode, 'Auth error.'));

            cb(null, res);
        });
        request.on('error', (e) => log.error(e));
        request.write(getJson({
            "path": fileDb.path
        }));
        request.end();
    }

    /**
     *
     * @param fileConf
     * @param cb
     */
    existsFile (fileDb, cb) {

        let requestParams = {
            method: 'POST',
            path: '/2/files/get_metadata',
            host: 'api.dropboxapi.com',
            headers: {
                'Authorization': 'Bearer ' + this._token,
                'Content-Type': 'application/json'
            }
        };

        let request = https.request(requestParams, (res) => {
            log.trace(`Response code get metadata file: ${res.statusCode}`);

            if (res.statusCode === 200) {
                return cb(null, true)
            } else if (res.statusCode === 404) {
                return cb(null, false)
            } else {
                return cb(new CodeError(500, 'Internal error'))
            }
        });
        request.on('error', (e) => log.error(e));
        request.write(getJson({
            "path": fileDb.path
        }));
        request.end();
    }

    /**
     *
     * @param conf
     * @param mask
     */
    checkConfig (conf = {}, mask) {
        return this.mask === mask && this._token === conf.accessToken;
    }
};

function getJson(data = {}) {
    try {
        return JSON.stringify(data)
    } catch (e) {
        log.error(e);
        return ""
    }
}