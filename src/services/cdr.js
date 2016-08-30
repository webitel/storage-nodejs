/**
 * Created by igor on 30.08.16.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module),
    async = require('async')
    ;

let _elasticConnect = true;
    
const Service = module.exports = {
    save: (cdrData, callback) => {
        let data = replaceVariables(cdrData);

        async.waterfall(
            [
                (cb) => {
                    if (data.callflow instanceof Array &&  /^u:/.test(data.callflow[0].caller_profile.destination_number)) {
                        data.callflow[0].caller_profile.destination_number = data.variables.presence_id;
                    }
                    if (data && data.variables && !data.variables.domain_name && /@/.test(data.variables.presence_id)) {
                        data.variables.domain_name = data.variables.presence_id.split('@')[1];
                    }
                    application.DB._query.cdr.insert(data, cb);
                },
                (result, cb) => {
                    if (application.elastic && result && result.ops && result.ops[0]) {
                        let _id = result.ops[0]._id;
                        
                        application.elastic.insertCdr(result.ops[0], (err) => {
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
    }
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