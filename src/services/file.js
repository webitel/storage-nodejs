/**
 * Created by igor on 25.08.16.
 */

"use strict";

const fs = require('fs'),
    crypto = require('crypto'),
    generateUuid = require('node-uuid')
;
    
const Service = module.exports = {
    requestStreamToCache: (fileName, req, cb) => {
        let pathFile = `cache/${generateUuid.v4()}_${fileName}`,
            stream = fs.createWriteStream(pathFile),
            result = {
                sha1: crypto.createHash('sha1'),
                path: pathFile,
                name: fileName,
                contentType: req.headers['content-type'],
                contentLength: +req.headers['content-length'] || 0
            }
            ;
        req.pipe(stream);
        req.on('data', (chunk) => {
            result.sha1.update(chunk);
        });
        req.on('end', () => {
            result.sha1 = result.sha1.digest('hex');
            return cb(null, result)
        });
        req.on('error', cb);
    },

    deleteFile: (filePath, cb) => {
        fs.unlink(filePath, cb)
    }
};