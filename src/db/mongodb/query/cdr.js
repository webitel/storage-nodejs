/**
 * Created by igor on 30.08.16.
 */

"use strict";

const conf = require(__appRoot + '/conf'),
    ObjectId = require('mongodb').ObjectId,
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

        getByObjId: (id, cb) => {
            return db
                .collection(cdrCollectionName)
                .findOne({_id: id}, cb);
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
        },

        search: (query, columns, sort, skip, limit, cb) => {
            return db
                .collection(cdrCollectionName)
                .find(query, columns)
                .sort(sort)
                .skip(skip)
                .limit(limit)
                .toArray(cb);
        },

        count: (query, cb) => {
            return db
                .collection(cdrCollectionName)
                .find(query)
                .count(cb);
        },

        buildFilterQuery: (filter) => {
            let filterArray = [];
            //filterArray.push(filterLegA); // TODO
            if (filter) {
                for (let key in filter) {
                    if (key == '_id' && ObjectId.isValid(filter[key])) {
                        filter[key] = ObjectId(filter[key]);
                        continue;
                    }
                    for (let item in filter[key]) {
                        // TODO ... parse _id
                        if (key == '_id' && ObjectId.isValid(filter[key][item])) {
                            filter[key][item] = ObjectId(filter[key][item]);
                        }
                    }
                }
                filterArray.push(filter)
            }

            return {
                "$and": filterArray
            };
        },

        remove: (uuid, cb) => {
            return db
                .collection(cdrCollectionName)
                .remove({"variables.uuid": uuid}, cb);

        },

        getIdFromUuid: (uuid, cb) => {
            return db
                .collection(cdrCollectionName)
                .findOne({"variables.uuid": uuid}, {_id: 1, "variables.domain_name": 1}, cb);
        }
    }
}