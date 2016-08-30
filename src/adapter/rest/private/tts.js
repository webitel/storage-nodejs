/**
 * Created by igor on 30.08.16.
 */

"use strict";

const TTSMiddleware = require(__appRoot + '/services/tts');

module.exports = {
    addRoutes: addRoutes
};

/**
 * Adds routes to the api.
 */
function addRoutes(api) {
    api.get('/sys/tts/:provider', TTSMiddleware);
}