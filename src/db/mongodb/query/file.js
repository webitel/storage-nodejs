/**
 * Created by igor on 25.08.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    ObjectId = require('mongodb').ObjectId,
    fileCollectionName = conf.get('mongodb:collectionFile'),
    domainsCollectionName = conf.get('mongodb:collectionDomain');

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        get: (id, pathName, contentType,  cb) => {
            let query;

            switch (contentType) {
                case 'all':
                    query = {
                        "uuid": id
                    };
                    break;
                case 'audio/mpeg':
                    let _r = [];

                    if (pathName) {
                        _r.push({
                            "name": pathName
                        });
                    } else {
                        _r.push({
                            "content-type": contentType
                        });
                        _r.push({
                            "content-type": {
                                "$exists": false
                            }
                        });
                    }

                    query = {
                        "$and": [
                            {
                                "uuid": id
                            },{
                                "$or": _r
                            }
                        ]
                    };
                    break;
                default:
                    query = {
                        "uuid": id
                    };
                    if (contentType)
                        query["content-type"] = contentType;

                    if (pathName)
                        query['name'] = pathName;

                    break;
            }

            return db
                .collection(fileCollectionName)
                .find(query)
                .toArray(cb);
        },

        getByObjId: (id, cb) => {
            return db
                .collection(fileCollectionName)
                .findOne({_id: id}, cb);
        },

        insert: (doc, cb) => {
            return db
                .collection(fileCollectionName)
                .insert(doc, cb)
        },

        deleteById: (id, cb) => {
            return db
                .collection(fileCollectionName)
                .remove({_id: id}, cb)
        },

        getFilesStats: (uuid, domain, option, cb) => {
            let _q = {
                    "$and": []
                },
                $and = [],
                _date = {}
                ;

            if (option['start']) {
                _date['$gte'] = option['start'];
            }

            if (option['end']) {
                _date['$lte'] = option['end'];
            }

            if (Object.keys(_date).length > 0) {
                $and.push({
                    "createdOn": _date
                });
            }

            if (domain) {
                $and.push({
                    "domain": domain
                });
            }

            _q['$and'] = $and;

            if (uuid) {
                _q['$and'].push({
                    "uuid": uuid
                });

                return db
                    .collection(fileCollectionName)
                    .findOne(_q, {"size": 1, "_id": 0}, cb);
            } else {

                let aggr = [];
                if ($and.length > 0)
                    aggr.push({
                        "$match": _q
                    });

                aggr = aggr.concat(
                    {"$group": {"_id": null, "size": {"$sum": "$size"}}},
                    {"$project": {"_id": 0, "size": 1}}
                );

                return db
                    .collection(fileCollectionName)
                    .aggregate(aggr, cb);
            }
        },

        updateFile: (uuid, id, data, cb) => {
            const $set = {};
            for (let key in data) {
                if (data.hasOwnProperty(key) && key !== "domain")
                    $set[key] = data[key];
            }

            return db
                .collection(fileCollectionName)
                .update({
                    _id: ObjectId.isValid(id) ? ObjectId(id) : id,
                    uuid: uuid
                }, {$set}, cb);
        },

        getStreamByDateRange: (start, end, columns) => {
            return db
                .collection(fileCollectionName)
                .find(
                    {
                        // "type": 0,
                        "createdOn": {
                            "$gte": start,
                            "$lte": end
                        }
                    },
                    columns
                )
                .stream()
        },

        getStreamByAggregateOldFile: (configServerDays = 30) => {
            return db
                .collection(fileCollectionName)
                .aggregate([
                    // Skip lock file
                    {
                        $match: {"_lock": {$ne: true}}
                    },
                    {
                        $unwind: "$domain"
                    },
                    {
                        $lookup:
                        {
                            from: domainsCollectionName,
                            localField: "domain",
                            foreignField: "name",
                            as: "_domain"
                        }
                    },
                    // Calc times expires days
                    {
                        $project: {
                            "uuid": 1,

                            "private": 1,
                            "domain": 1,
                            "bucketName": 1,
                            "path": 1,
                            "storageFileId": 1,
                            "type": 1,

                            "_expiresDays": {
                                $arrayElemAt: ["$_domain.storage.expiresDays", 0]
                            },
                            "createdOn": 1
                        }
                    },
                    // set server value if no domain config
                    {
                        $project: {
                            "uuid": 1,

                            "private": 1,
                            "domain": 1,
                            "bucketName": 1,
                            "path": 1,
                            "storageFileId": 1,
                            "type": 1,

                            "_expiresDays": {
                                $ifNull: ["$_expiresDays", configServerDays]
                            },
                            "createdOn": 1
                        }
                    },
                    // set 0 - never delete
                    {
                        $match: {"_expiresDays": {$ne: 0}}
                    },
                    // calc deadline
                    {
                        $project: {
                            "uuid": 1,

                            "private": 1,
                            "domain": 1,
                            "bucketName": 1,
                            "path": 1,
                            "storageFileId": 1,
                            "type": 1,

                            "_deadlineMS": {
                                $multiply: ["$_expiresDays", 86400000] //24*60*60000
                            },
                            "createdOn": 1
                        }
                    },
                    {
                        $project: {
                            "uuid": 1,
                            "createdOn": 1,

                            "private": 1,
                            "domain": 1,
                            "bucketName": 1,
                            "path": 1,
                            "storageFileId": 1,
                            "type": 1,

                            "deadlineDate": {
                                $add: ["$createdOn", "$_deadlineMS"]
                            }
                        }
                    },
                    {
                        $project: {
                            "uuid": 1,
                            "createdOn": 1,

                            "private": 1,
                            "domain": 1,
                            "bucketName": 1,
                            "path": 1,
                            "storageFileId": 1,
                            "type": 1,

                            "deadline": {
                                $subtract: ["$deadlineDate", new Date()]
                            }
                        }
                    },
                    {
                        $match: {"deadline": {$lte: 0}}
                    }

                ])
                .stream()
        }
    }
}