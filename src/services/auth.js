/**
 * Created by igor on 23.08.16.
 */

"use strict";

const config = require(__appRoot + '/config'),
    CodeError = require(__appRoot + '/lib/error'),
    log = require(__appRoot + '/lib/log')(module),
    inSubnet = require('insubnet')
    ;

const MODE = config.get("uploadAcl:mode"),
    IPS = config.get("uploadAcl:ip")
    ;

let Service = {
    checkBannedIp: (ip) => {
        ip.replace(/^::ffff:/, '');

        if (MODE === 'allow' && inSubnet.Auto(ip, IPS)) {
            return null;
        } else {
            log.warn(`Unauthorized connection ip: ${ip}`);
            return new CodeError(401, `Unauthorized connection ip: ${ip}`)
        }
    },

    getUserByKey: (key, cb) => {
        try {
            var authDb = application.DB._query.auth;
            authDb.getByKey(key, cb);
        } catch (e){
            cb(e);
        }
    }
};

module.exports = Service;