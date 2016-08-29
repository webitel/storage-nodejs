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
                        "uuid": id,
                        "content-type": contentType
                    };
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
        }
    }
}