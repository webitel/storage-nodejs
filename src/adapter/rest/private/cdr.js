/**
 * Created by igor on 26.08.16.
 */

"use strict";

const CodeError = require(__appRoot + '/lib/error')
    ;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    // TODO
    api.post('/sys/cdr', (req, res) => res.end() )
}