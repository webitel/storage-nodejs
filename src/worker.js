/**
 * Created by igor on 23.08.16.
 */

"use strict";

const path = require('path');

global.__appRoot = path.resolve(__dirname);

const APPLICATION = require('./application'),
    application = global.application = APPLICATION;

// process.once('SIGINT', function() {
//     console.log('SIGINT received ...');
//     if (application) {
//         application.stop();
//     }
// });

module.exports = application;