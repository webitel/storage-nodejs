/**
 * Created by igor on 01.09.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    mediaCollectionName = conf.get('mongodb:collectionMedia');

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        insert: (doc, cb) => {
            return db
                .collection(mediaCollectionName)
                .insert(doc, cb)
        },

        listFromDomain: (domainName, cb) => {
            return db
                .collection(mediaCollectionName)
                .find({domain: domainName})
                .toArray(cb);
        },

        getItem: (type, name, domainName, cb) => {
            return db
                .collection(mediaCollectionName)
                .findOne({
                    name: name,
                    domain: domainName,
                    type: type
                }, cb)
        },

        deleteById: (id, cb) => {
            return db
                .collection(mediaCollectionName)
                .remove({
                    _id: id
                }, cb)
        }
    }
}