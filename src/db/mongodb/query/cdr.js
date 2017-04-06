/**
 * Created by igor on 30.08.16.
 */

"use strict";

const conf = require(__appRoot + '/conf'),
    ObjectId = require('mongodb').ObjectId,
    generateUuid = require('node-uuid'),
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

        _setTryToElastic: (count, iterator, bulkCb, finishCb) => {
            const collection = db
                .collection(cdrCollectionName);

            const operationId = generateUuid.v4();
            const bulk = collection.initializeUnorderedBulkOp();

            const len = count < 10000 ? count : 10000;
            for (let i = 0; i < len; i++) {
                bulk.find( { _elasticExportError: true} ).updateOne( { $set: { _elasticExportError: operationId } } );
            }

            const endCallProcess = (err, errIds) => {
                if (err) {
                    resetProcessBulk(collection, operationId, true);
                    return finishCb(err);
                }

                if (errIds && errIds.length > 0) {
                    return resetProcessBulkFromIds(collection, operationId, errIds, err => {
                        finishCb(err);
                    })
                } else {
                    resetProcessBulk(collection, operationId, null)
                }

                return finishCb(null);
            };

            console.time(`Bulk operation export to elastic ${operationId}`);
            bulk.execute((err) => {
                if (err)
                    return finishCb(err);

                console.timeEnd(`Bulk operation export to elastic ${operationId}`);
                const cursor = collection.find({_elasticExportError: operationId});
                let data = [];
                cursor.each( (err, doc) => {
                    if (err) {
                        return endCallProcess(err);
                    }

                    if (doc) {
                        iterator(doc, data);
                    } else {
                        bulkCb(data, endCallProcess);
                    }
                })

            });

        },

        getIdFromUuid: (uuid, cb) => {
            return db
                .collection(cdrCollectionName)
                .findOne({"variables.uuid": uuid}, {_id: 1, "variables.domain_name": 1}, cb);
        }
    }
}

function resetProcessBulk(collection, operationId, value) {
    collection.update({ _elasticExportError: operationId}, {$set: { _elasticExportError: value}}, {multi: true});
}


function resetProcessBulkFromIds(collection, operationId, ids, cb) {
    const bulk = collection.initializeUnorderedBulkOp();
    for (let i = 0; i < ids.length; i++) {
        bulk.find( { _elasticExportError: operationId, _id: ObjectId(ids[i])} ).updateOne( { $set: { _elasticExportError: true } } );
    }
    bulk.execute(cb);
}