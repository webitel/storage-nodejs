/**
 * Created by igor on 30.08.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    cdrCollectionName = conf.get('mongodb:collectionCDR');

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        insert: (doc, cb) => {
            return db
                .collection(cdrCollectionName)
                .insert(doc, cb);
        },

        setById: (id, properties, cb) => {
            return db
                .collection(cdrCollectionName)
                .updateOne({_id: id}, {$set: properties}, cb);
        },

        unsetById: (id, properties, cb) => {
            return db
                .collection(cdrCollectionName)
                .updateOne({_id: id}, {$unset: properties}, cb);
        },

        find: (query, cb) => {
            return db
                .collection(cdrCollectionName)
                .find(query)
                .toArray(cb);
        }
    }
}