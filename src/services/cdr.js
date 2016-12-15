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

let _elasticConnect = true;
    
const Service = module.exports = {
    
    save: (cdrData, params, callback) => {
        if (typeof params === 'function') {
            callback = params;
            params = {}
        }

        let data = replaceVariables(cdrData);

        if (data.variables &&
            ( (data.variables.loopback_leg === "A" && data.variables.loopback_bowout_on_execute !== 'true')
                || (data.variables.loopback_bowout_on_execute === 'true' && data.variables.loopback_leg === "B") )) {
            log.debug(`Skip leg ${data.variables.loopback_leg} ${data.variables.uuid}`);
            return callback(null);
        }
        async.waterfall(
            [
                (cb) => {
                    if (data.callflow instanceof Array &&  /^u:/.test(data.callflow[0].caller_profile.destination_number)) {
                        data.callflow[0].caller_profile.destination_number = data.variables.presence_id;
                    }
                    if (data && data.variables && !data.variables.domain_name && /@/.test(data.variables.presence_id)) {
                        data.variables.domain_name = data.variables.presence_id.split('@')[1];
                    }

                    if (params && params.skipMongo === true) {
                        if (!data._id)
                            return cb(new CodeError(403, `Bad cdr. Field _id is require`));
                        return cb(null, data, data._id)
                    } else {
                        application.DB._query.cdr.insert(data, (err, result) => {
                            if (err)
                                return cb(err);

                            if (result && result.ops) {
                                return cb(null, result.ops[0], result.insertedIds && result.insertedIds[0])
                            }
                            return cb (new CodeError(500, `Bad create db record.`));
                        });
                    }
                },
                (result, newId, cb) => {
                    if (application.replica)
                        application.replica.sendCdr(data, newId);

                    if (application.elastic && result) {
                        let _id = result._id;
                        
                        application.elastic.insertCdr(result, (err) => {
                            if (err && !~err.message.indexOf('document_already_exists_exception')) {
                                log.warn(`no save elastic: ${err}`);
                                _elasticConnect = false;
                                return application.DB._query.cdr.setById(_id, {"_elasticExportError": true}, cb);
                            } else {
                                if (_elasticConnect === false)
                                    processSaveToElastic();
                                _elasticConnect = true;
                            }

                            return cb(err)
                        });
                    } else {
                        cb();
                    }
                }
            ],
            callback
        )
    },

    setValideAttrDoc: (doc = {}) => {
        let data = replaceVariables(doc);
        if (data.callflow instanceof Array &&  /^u:/.test(data.callflow[0].caller_profile.destination_number)) {
            data.callflow[0].caller_profile.destination_number = data.variables.presence_id;
        }
        if (data && data.variables && !data.variables.domain_name && /@/.test(data.variables.presence_id)) {
            data.variables.domain_name = data.variables.presence_id.split('@')[1];
        }
        return data;
    },

    saveToElastic: (doc, cb) => {
        application.elastic.insertCdr(Service.setValideAttrDoc(doc), cb);
    },

    search: (caller, option, cb) => {
        let _ro = false
            ;

        let columns = option.columns || DEF_COLUMNS,
            sort = option.sort  || {
                "callflow.times.created_time": -1
            },
            limit = parseInt(option.limit, 10) || 40,
            pageNumber = option.pageNumber
            ;

        let query = application.DB._query.cdr.buildFilterQuery(option.filter);

        if (caller.domain)
            query['$and'].push({
                "variables.domain_name": caller.domain
            });

        let _readAll = checkPermission(caller.acl, 'cdr', 'r');

        if (!_readAll && checkPermission(caller.acl, 'cdr', 'ro', true)) {
            query['$and'].push({
                "variables.presence_id": caller.id
            });
            _ro = true;
        }

        if (!_ro && !_readAll) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        application.DB._query.cdr.search(
            query,
            columns,
            sort,
            pageNumber > 0 ? ((pageNumber - 1) * limit) : 0,
            limit,
            cb
        );
    },

    count: (caller, option, cb) => {
        let _ro = false
            ;
        
        let query = application.DB._query.cdr.buildFilterQuery(option.filter);

        if (caller.domain)
            query['$and'].push({
                "variables.domain_name": caller.domain
            });

        let _readAll = checkPermission(caller.acl, 'cdr', 'r');

        if (!_readAll && checkPermission(caller.acl, 'cdr', 'ro', true)) {
            query['$and'].push({
                "variables.presence_id": caller.id
            });
            _ro = true;
        }

        if (!_ro && !_readAll) {
            return cb(new CodeError(403, "Permission denied!"))
        }

        application.DB._query.cdr.count(
            query,
            cb
        );
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
                                        application.elastic.removeCdr(docId, data[0]._index, err => cb(err));
                                    } else {
                                        return cb();
                                    }
                                });
                            } else {
                                return cb(null)
                            }
                        },

                        (cb) => {
                            application.DB._query.cdr.remove(uuid, cb);
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
    }
    
};

const DEF_COLUMNS = {
    "variables.uuid": 1,
    "callflow.caller_profile.caller_id_name": 1,
    "callflow.caller_profile.caller_id_number": 1,
    "callflow.caller_profile.callee_id_number": 1,
    "callflow.caller_profile.callee_id_name": 1,
    "callflow.caller_profile.destination_number": 1,
    "callflow.times.created_time": 1,
    "callflow.times.answered_time": 1,
    "callflow.times.bridged_time": 1,
    "callflow.times.hangup_time": 1,
    "variables.duration": 1,
    "variables.hangup_cause": 1,
    "variables.billsec": 1,
    "variables.direction": 1
};

function processSaveToElastic() {
    application.DB._query.cdr.find({"_elasticExportError": true}, (err, data) => {
        if (err) {
            return log.error(err);
        }

        if (data instanceof Array) {
            async.every(
                data,
                (doc, cb) => {
                    let _id = doc._id;
                    application.elastic.insertCdr(doc, (err) => {
                        if (err && !~err.message.indexOf('document_already_exists_exception'))
                            return cb(err);
                        log.debug(`Save elastic document: ${_id}`);
                        return application.DB._query.cdr.unsetById(_id, {"_elasticExportError": true}, cb);
                    })
                },
                (err) => {
                    if (err)
                        log.error(err);

                }
            )
        } else {
            log.error(`Bad response find no save elastic data`);
        }
    })
}

function replaceVariables(data) {

    for (let key in data.variables) {
        if (isFinite(data.variables[key]))
            data.variables[key] = +data.variables[key];

        if (/\.|\$/.test(key)) {
            data.variables[encodeKey(key)] = data.variables[key];
            delete data.variables[key];
        }
    }
    if (data.callflow instanceof Array) {
        for (let cf of data.callflow) {
            if (cf.hasOwnProperty('times')) {
                for (let key in cf.times) {
                    cf.times[key] = +cf.times[key];
                }
            }
        }
    }
    return data
}

function encodeKey(key) {
    return key.replace(/\\/g, "\\\\").replace(/\$/g, "\\u0024").replace(/\./g, "\\u002e")
}