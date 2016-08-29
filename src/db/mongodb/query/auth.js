/**
 * Created by igor on 23.08.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    authCollectionName = conf.get('mongodb:collectionAuth');

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        getByKey: function (key, cb) {
            return db
                .collection(authCollectionName)
                .findOne({"key": key}, cb);
        }
    }
}