/**
 * Created by I. Navrotskyj on 30.10.17.
 */

"use strict";

const checkPermission = require(__appRoot + '/utils/acl');
const elasticService = require(__appRoot + '/services/elastic');

const Service = module.exports = {
    request: function(data = {}, api, cb) {

        elasticService.search(checkPermission.GOD, data, (err, res) => {
            cb(null, {
                "exec-args": {
                    "callId": api.properties.correlationId,
                    "data": res
                }
            })
        });

    }
};
