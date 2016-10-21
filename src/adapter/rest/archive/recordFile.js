/**
 * Created by igor on 21.10.16.
 */

"use strict";

module.exports = {
    addRoutes: api => {
        api.post('/api/v1/recordings', create);
        api.get('/api/v1/recordings/:id', create);
    }
};

const create = (req, res, next) => {
    return res.send('OK');
};