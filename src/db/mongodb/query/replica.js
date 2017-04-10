/**
 * Created by igor on 21.10.16.
 */

"use strict";

const conf = require(__appRoot + '/conf'),
    log = require(__appRoot + '/lib/log')(module),
    uuid = require('node-uuid'),
    replicaCollectionName = conf.get('mongodb:collectionReplica');

module.exports = {
    addQuery: (db) => {
        return {
            insert: (doc, cb) => {
                return db
                    .collection(replicaCollectionName)
                    .insert(doc, cb)
            },

            getOne: cb => {
                return db
                    .collection(replicaCollectionName)
                    .findOneAndDelete(
                        {},
                        {
                            sort: {createdOn: 1}
                        },
                        cb
                    )
            },

            sync: (collectionName, maxCount = 1000, send, cb) => {
                const processId = uuid.v4();

                const collectionReplica = db
                    .collection(replicaCollectionName);

                const bulk = collectionReplica.initializeUnorderedBulkOp();
                const mapIds = new Map();

                const cursor = collectionReplica
                    .find({"collection": collectionName, process: null}, {_id: 1})
                    .sort({createdOn: 1})
                    .limit(maxCount);

                cursor.each( (err, item) => {
                    if (err) {
                        return end(err);
                    }

                    if (item) {
                        bulk.find({_id: item._id, process: null}).updateOne({$set: {process: processId}})
                    } else {
                        if (!bulk.length) {
                            return end(new Error('No found data'));
                        }
                        bulk.execute( (err, res) => {
                            if (err)
                                return end(err);
                            return getData();
                        });
                    }
                });


                const end = (err, errorsIds) => {
                    if (err) {
                        errorProcess();
                        return cb(err);
                    }

                    if (errorsIds && errorsIds.length > 0) {
                        return resetIds(errorsIds, err => {
                            if (err) {
                                log.error(err);
                                return cb(err);
                            }
                            return setOkProcess(cb);

                        })
                    } else {
                        return setOkProcess(cb);
                    }
                };

                const errorProcess = (cb) => {
                    mapIds.clear();
                    collectionReplica.update({process: processId}, {$set: {process: null}}, {multi: true}, cb);
                };

                const resetIds = (ids, cb) => {
                    const bulk = collectionReplica.initializeUnorderedBulkOp();

                    ids.forEach( id => {
                        if (mapIds.has(id)) {
                            bulk.find({_id: mapIds.get(id)}).updateOne({$set: {process: null}})
                        } else {
                            log.warn(`No found bad id ${id}`);
                        }
                    });


                };

                const setOkProcess = (cb) => {
                    mapIds.clear();
                    collectionReplica.remove({process: processId}, {multi: true}, cb)
                };

                const getData = () => {
                    const data = [];
                    const c = collectionReplica
                        .aggregate([
                            {$match: {"collection": collectionName, process: processId}},
                            {
                                $lookup: {
                                    from: collectionName,
                                    localField: "docId",
                                    foreignField: "_id",
                                    as: "doc"
                                }
                            },
                            {$unwind: "$doc"}
                        ]);

                    c.each( (e, item) => {
                        if (e) {
                            return end(e);
                        }
                        if (item && item.doc) {
                            data.push(item.doc);
                            if (item.doc.variables) {
                                mapIds.set(item.doc.variables.uuid, item._id);
                            }
                        } else {
                            send(data, end);
                        }
                    })
                };
            },

            setProcessId: (collectionName, process, _id, cb) => {
                return db
                    .collection(replicaCollectionName)
                    .updateOne({_id}, {$set: {process}}, cb);
            },

            removeById: id => {
                return db
                    .collection(replicaCollectionName)
                    .removeOne({_id: id})
            }
        }
    }
};
