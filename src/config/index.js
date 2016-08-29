/**
 * Created by igor on 23.08.16.
 */

"use strict";

const nconf = require('nconf'),
    path = require('path');

//nconf
nconf.argv()
    .env()
    .add('elastic', {
        type: "file",
        file: path.join(__dirname, 'elastic.json')
    })
    .file({
        file: path.join(__dirname, 'config.json')
    })
;

module.exports = nconf;