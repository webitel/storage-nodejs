/**
 * Created by I. Navrotskyj on 05.10.17.
 */

"use strict";

const recordingsService = require(`${__appRoot}/services/recordings`);
const log = require(`${__appRoot}/lib/log`)(module);

const Service = module.exports = {
    removeFile: (data = {}, api, cb) => {
        if (!data) return cb();

        const providerName = recordingsService.getProviderNameFromFile(data);
        if (!providerName) {
            log.error('Not found provider in: ', data);
            return cb();
        }

        recordingsService._delFile(providerName, data, {}, err => {
            if (err) {
                return log.error(err);
            }

            log.trace(`Remove file: ${data.name}`)
        });
        return cb();
    }
};