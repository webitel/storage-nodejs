/**
 * Created by igor on 21.10.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    replicaCollectionName = conf.get('mongodb:collectionReplica');

module.exports = {
    addQuery: (db) => {
        return {
            insert: (doc, cb) => {
                return db
                    .collection(replicaCollectionName)
                    .insert(doc, cb)
            }
        }
    }
};
