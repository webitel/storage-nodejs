/**
 * Created by igor on 23.08.16.
 */

"use strict";

var util = require('util');

class CodeError extends Error {
    constructor (status, message) {
        super();
        this.name = 'CodeError';
        this.status = status;
        this.message = message;
        Error.captureStackTrace(this, CodeError);
    }
}

module.exports = CodeError;