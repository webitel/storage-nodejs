/**
 * Created by igor on 25.08.16.
 */

"use strict";

const fs = require('fs'),
    generateUuid = require('node-uuid')
;
    
const Service = module.exports = {
    requestStreamToCache: (fileName, req, cb) => {
        let pathFile = `cache/${generateUuid.v4()}_${fileName}`,
            stream = fs.createWriteStream(pathFile),
            result = {
                data: new Buffer([]),
                path: pathFile,
                name: fileName,
                contentType: req.headers['content-type'],
                contentLength: +req.headers['content-length'] || 0
            }
            ;
        req.pipe(stream);
        req.on('data', (chunk) => {
            result.data = Buffer.concat([result.data, chunk], result.data.length + chunk.length);
        });
        req.on('end', () => {
            return cb(null, result)
        });
        req.on('error', cb);
    },

    deleteFile: (filePath, cb) => {
        fs.unlink(filePath, cb)
    }
};