/**
 * Created by igor on 25.08.16.
 */

"use strict";

const conf = require(__appRoot + '/conf'),
    domainCollectionName = conf.get('mongodb:collectionDomain');

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        getByName: function (name, columnName, cb) {
            let col = {};
            col[columnName] = 1;
            return db
                .collection(domainCollectionName)
                .findOne({"name": name}, col, cb);
        }
    }
}