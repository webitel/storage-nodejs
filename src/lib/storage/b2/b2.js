/**
 * Created by igor on 29.08.16.
 */

"use strict";

/**
 * Created by igor on 25.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    CodeError = require(__appRoot + '/lib/error'),
    https = require('https'),
    url = require('url'),
    fs = require('fs'),
    crypto = require('crypto')
    ;

const HOST = 'api.backblaze.com',
    AUTH_PATH = '/b2api/v1/b2_authorize_account',
    TYPE_ID = 2;

const B2 = module.exports = {
    auth: (args, cb) => {
        if (!args || !args.accountId || !args.applicationKey || !args.bucketId || !args.bucketName)
            return cb(new CodeError(401, 'Bad auth parameters'));

        let requestParams = {
            method: 'GET',
            path: AUTH_PATH,
            host: HOST,
            headers: {
                'Authorization': `Basic ${new Buffer(args.accountId + ':' + args.applicationKey).toString('base64')}`
            }
        };
        log.debug(`Try auth B2`);
        return getJson(requestParams, (err, auth) => {
            if (err)
                return cb(err);

            auth.apiHost = (auth.apiUrl || '').replace('https://', '');
            return B2.setUploadUrl(auth, args.bucketId, args.bucketName, cb)
        })
    },

    setUploadUrl: (credential, bucketId, bucketName, cb) => {
        let body = `{"bucketId": "${bucketId}"}`;
        let requestParams = {
            method: 'POST',
            path: '/b2api/v1/b2_get_upload_url',
            host: credential.apiHost,
            headers: {
                'Authorization': credential.authorizationToken,
                'Content-Length': body.length
            },
            body: body

        };
        return getJson(requestParams, (err, res) => {
            if (err)
                return cb(err);
            credential.uploadParams = url.parse(res.uploadUrl);
            credential.uploadToken = res.authorizationToken;
            credential.bucketId = bucketId;
            credential.bucketName = encodeBucketName(bucketName);
            credential.uploadUrl = res.uploadUrl;
            return cb(null, credential)
        })
    },

    getFile: (credential, fileConf, range, cb) => {
        let option = {
            method: 'GET',
            path: '/file/' + fileConf.bucketName + '/' + getUrlEncodedFileName(fileConf.path),
            host: credential.downloadUrl.replace('https://', ''),
            headers: {
                'Authorization': credential.authorizationToken
            }
        };
        if (range && range.Start < range.End) {
            option.headers['Range'] = `bytes=${range.Start}-${range.End}`;
        }

        let request = https.request(option, (res) => {
            log.trace(`Response code get file: ${res.statusCode}`);

            if (res.statusCode !== 200 && res.statusCode !== 206)
                return cb(new CodeError(res.statusCode, 'Auth error.'));

            cb(null, res);
        });
        request.on('error', (e) => {
            log.error(e);
        });
        request.end();

    },

    delFile: (credential, fileConf, cb) => {

        let data = `{"fileId":"${fileConf.storageFileId}","fileName":"${fileConf.path.replace(/^\//, '')}"}`;

        var requestParams = {
            method: 'POST',
            path: '/b2api/v1/b2_delete_file_version',
            host: credential.apiHost,
            headers: {
                'Authorization': credential.authorizationToken,
                'Content-Type': 'application/json; charset=utf-8',
                'Content-Length': data.length
            },
            body: data
        };

        getJson(requestParams, cb);
    },

    saveFile: (credential, fileConf, fileName, cb) => {
        let mime = fileConf.contentType,
            data = fileConf.data
            ;

        var requestParams = {
            method: 'POST',
            path: credential.uploadParams.path, //`/file/${credential.bucketName}/${filename}`,
            host: credential.uploadParams.host,
            headers: {
                'Authorization': credential.uploadToken,
                'Content-Type': mime || 'b2/x-auto',
                'Content-Length': fileConf.contentLength,
                'X-Bz-File-Name': getUrlEncodedFileName(fileName),
                'X-Bz-Content-Sha1': fileConf.sha1 ? fileConf.sha1 : null
            }
        };

        let request = https.request(requestParams, (res) => {
            log.trace(`B2 upload status code `, res.statusCode);

            log.trace(`Save storage file path: ${fileName}`);

            let data = '';
            res.on('data', function(chunk) {
                data += chunk;
            });
            res.on('end', function() {
                try {
                    if (res.statusCode != 200) {
                        return cb(new CodeError(res.statusCode, data || "Internal error."));
                    }

                    let json = JSON.parse(data);
                    return cb(null, {
                        path: fileName,
                        type: TYPE_ID,
                        bucketName: credential.bucketName,
                        storageFileId: json.fileId
                    });
                } catch (e) {
                    return cb (e)
                }
            });
        });
        request.on('error', (e) => {
            log.error(e);
        });

        let rd = fs.createReadStream(fileConf.path);
        rd.on("error", function(e) {
            log.error(e);
            return cb(e);
        });

        rd.on('open', () => rd.pipe(request));
    }
};

function getUrlEncodedFileName (fileName) {
    return fileName.replace(/^\//, '').split('/')
        .map(encodeURIComponent)
        .join('/');
}

function getJson(requestParams, cb) {
    let data = '';
    let request = https.request(requestParams, (res) => {
        // if (res.statusCode !== 200)
        //     return cb(new CodeError(res.statusCode, 'Auth error.'));

        res.once('error', cb);
        res.on('data', function(chunk) {
            data += chunk;
        });
        res.on('end', function() {
            try {
                return cb(res.statusCode !== 200 ? new Error(data) : null, JSON.parse(data));
            } catch (e) {
                return cb (e)
            }
        });
    });
    if (requestParams.body) request.write(requestParams.body);
    request.end();
}

function encodeBucketName(name) {
    return (name || '').split('').map( i => i.charCodeAt() == 8209 ? '-' : i  ).join('');
}