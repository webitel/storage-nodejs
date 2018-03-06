/**
 * Created by igor on 30.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    checkPermission = require(__appRoot + '/utils/acl'),
    recordingsService = require('./recordings'),
    CodeError = require(__appRoot + '/lib/error'),
    async = require('async')
    ;
    
const Service = module.exports = {

    addPin: (caller, option = {}, cb) => {
        if (!application.elastic) {
            return cb(new CodeError(500, "No live connect to elastic!"))
        }
        if (!checkPermission(caller.acl, 'cdr', 'u')) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        if (!option.index) {
            return cb(new CodeError(400, "Index is required!"))
        }

        if (caller.domain && !~option.index.indexOf(caller.domain)) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        application.elastic.addPinCdr(option.id, option.index, caller.id, cb)
    },

    delPin: (caller, option = {}, cb) => {
        if (!application.elastic) {
            return cb(new CodeError(500, "No live connect to elastic!"))
        }
        if (!checkPermission(caller.acl, 'cdr', 'u')) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        if (!option.index) {
            return cb(new CodeError(400, "Index is required!"))
        }

        if (caller.domain && !~option.index.indexOf(caller.domain)) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        application.elastic.delPinCdr(option.id, option.index, caller.id, cb)
    },

    remove: (caller, option, callback) => {
        const {uuid} = option;
        let domain = caller.domain,
            docId;

        if (!checkPermission(caller.acl, 'cdr', 'd')) {
            return callback(new CodeError(403, "Permission denied!"))
        }

        recordingsService.getFileFromUUID(caller, uuid, {contentType: 'all'}, (err, files) => {
            if (err)
                return callback(err);


            const delCdr = () => {
                async.waterfall(
                    [

                        (cb) => {
                            if (application.elastic) {
                                application.elastic.findByUuid(uuid, domain || "", (err, res) => {
                                    if (err)
                                        return cb(err);

                                    const data = res && res.hits && res.hits.hits;
                                    if (data && data.length > 0) {
                                        docId = data[0]._id;
                                        application.elastic.removeCdr(uuid, data[0]._index, err => cb(err));
                                    } else {
                                        return cb();
                                    }
                                });
                            } else {
                                return cb(null)
                            }
                        },

                        (cb) => {
                            application.PG.getQuery('cdr').removeLegA(uuid, cb);
                        }
                    ],
                    callback
                );
            };

            if (files instanceof Array && files.length > 0) {
                async.eachSeries(files,
                    (item, cb) => {
                        const providerName = recordingsService.getProviderNameFromFile(item);
                        domain = item.domain;
                        if (!providerName) {
                            log.warn(`skip: not found provider`);
                            return cb(null)
                        }
                        recordingsService._delFile(providerName, item, {delDb: true}, (err) => {
                            if (err)
                                log.error(err);
                            return cb(null);
                        });
                    },
                    (err) => {
                        if (err)
                            log.error(err);

                        return delCdr();
                    }
                );
            } else {
                delCdr()
            }
        })
    },

    getItem: (caller, options = {}, cb) => {
        if (!checkPermission(caller.acl, 'cdr', 'r')) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        application.PG.getQuery('cdr').getLegByUuid(options.uuid, null, options.leg, cb);
    }
};