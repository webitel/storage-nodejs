/**
 * Created by I. Navrotskyj on 04.10.17.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    recordingsService = require('./recordings'),
    CodeError = require(__appRoot + '/lib/error');

const Service = module.exports = {

    getFile: (caller, options = {}, cb) => {
        if (caller.domain) {
            return cb(new CodeError(401, "Permission denied"))
        }

        if (!isFinite(options.id)) {
            return cb(new CodeError(400, "Bad id"))
        }

        application.PG.getQuery('tcpDump').getMetadata(options.id, (err, fileDb) => {
            if (err)
                return cb(err);

            recordingsService._getFile(fileDb, {}, cb)
        });
    }
};