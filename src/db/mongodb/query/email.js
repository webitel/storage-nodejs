/**
 * Created by igor on 30.08.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    emailCollectionName = conf.get('mongodb:collectionEmail');

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        getByDomain: function (name, cb) {
            return db
                .collection(emailCollectionName)
                .findOne({"domain": name}, cb);
        }
    }
}