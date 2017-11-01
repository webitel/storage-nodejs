/**
 * Created by igor on 25.01.17.
 */

"use strict";

const log = require(__appRoot + '/lib/log')(module);

const Service = module.exports = {
    saveToElastic: (data, api, cb) => {
        if (application.elastic) {
            application.elastic.insertUserStatus(data, (e, res) => {
                if (e)
                    log.error(e);
            })
        }
        return cb()
    }
};