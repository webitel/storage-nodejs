/**
 * Created by igor on 29.08.16.
 */

"use strict";

const fs = require('fs-extra'),
    log = require(__appRoot + '/lib/log')(module),
    CodeError = require(__appRoot + '/lib/error'),
    path = require('path'),
    helper = require('./helper'),
    TYPE_ID = 0;

module.exports = class LocalStorage {

    constructor (conf, mask) {
        this.rootPath = conf.fileRoot;
        this.mask = mask;
        this.name = "local";
    }

    get (fileDb, options, cb) {
        fs.lstat(fileDb.path, function (err, stat) {
            if (err)
                return cb(err);

            if (!stat.isFile()) {
                return cb(new CodeError(404, 'Bad file.'))
            }

            let readable;

            if (options.range) {
                readable = fs.createReadStream(fileDb.path, {flags: 'r', start: options.range.start, end: options.range.end })
            } else {
                readable = fs.createReadStream(fileDb.path, {flags: 'r'});
            }

            readable.on('open', () => {
                return cb(null, readable);
            })
        });
    }

    save (fileConf, option, cb) {
        const pathFolder = path.join(this.rootPath, helper.getPath(this.mask, fileConf.domain));
        fs.ensureDir(pathFolder, (err) => {
            if (err)
                return cb(err);

            let filePath = pathFolder + '/' + fileConf.name;
            fs.copy(fileConf.path, filePath, {clobber: true}, function(err) {
                if (err)
                    return cb(err);
                log.trace(`Save file: ${filePath}`);
                return cb(null, {
                    path: filePath,
                    type: TYPE_ID
                })
            });
        })
    }

    validateConfig (config) {
        return !config || this.rootPath != config.rootPath || this.mask != config.maskPath;

    }
};