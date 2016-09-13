/**
 * Created by igor on 25.08.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    fileCollectionName = conf.get('mongodb:collectionFile');

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
        }
    }
}