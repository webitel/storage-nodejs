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
        return getJson(requestParams, (err, auth) => {
            if (err)
                return cb(err);

            return B2.setUploadUrl(auth, args.bucketId, args.bucketName, cb)
        })
    },

    setUploadUrl: (credential, bucketId, bucketName, cb) => {
        let body = `{"bucketId": "${bucketId}"}`;
        let requestParams = {
            method: 'POST',
            path: '/b2api/v1/b2_get_upload_url',
            host: credential.apiUrl.replace('https://', ''),
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
        if (range)
            option.headers['Range'] = `bytes=${range.Start}-${range.End}`;

        let request = https.request(option, (res) => {
            log.trace(`Response code get file: ${res.statusCode}`);

            if (res.statusCode !== 200)
                return cb(new CodeError(res.statusCode, 'Auth error.'));

            cb(null, res);
        });
        request.on('error', (e) => {
            log.error(e);
        });
        request.end();

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
                'Content-Length': data.length,
                'X-Bz-File-Name': fileName,
                'X-Bz-Content-Sha1': data ? sha1(data) : null
            },
            body: data
        };

        let request = https.request(requestParams, (res) => {
            log.trace(`B2 upload status code `, res.statusCode);

            if (res.statusCode !== 200) {
                log.error(`Upload params: ${requestParams}`);
                return cb(new CodeError(res.statusCode, 'Auth error.'));
            }

            log.trace(`Save storage file path: ${fileName}`);

            return cb(null, {
                path: fileName,
                type: TYPE_ID,
                bucketName: credential.bucketName
            });
            // let data = '';
            // res.on('data', function(chunk) {
            //     data += chunk;
            // });
            // res.on('end', function() {
            //     try {
            //         return cb(null, JSON.parse(data));
            //     } catch (e) {
            //         return cb (e)
            //     }
            // });
        });
        request.on('error', (e) => {
            log.error(e);
        });

        if (requestParams.body) request.write(requestParams.body);

        request.end();
    }
};

function sha1(str) {
    let hash = crypto.createHash('sha1');
    hash.update(str);
    return hash.digest('hex');
}

function getUrlEncodedFileName (fileName) {
    return fileName.split('/')
        .map(encodeURIComponent)
        .join('/');
}

function getJson(requestParams, cb) {
    let data = '';
    let request = https.request(requestParams, (res) => {
        if (res.statusCode !== 200)
            return cb(new CodeError(res.statusCode, 'Auth error.'));

        res.once('error', cb);
        res.on('data', function(chunk) {
            data += chunk;
        });
        res.on('end', function() {
            try {
                return cb(null, JSON.parse(data));
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