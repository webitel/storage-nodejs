/**
 * Created by igor on 23.08.16.
 */

"use strict";

const conf = require(__appRoot + '/config'),
    authCollectionName = conf.get('mongodb:collectionAuth'),
    domainCollectionName = conf.get('mongodb:collectionDomain'),
    aclCollectionName = conf.get('mongodb:collectionAcl')
    ;

module.exports = {
    addQuery: addQuery
};

function addQuery(db) {
    return {
        getByKey: function (key, cb) {
            return db
                .collection(authCollectionName)
                .findOne({"key": key}, cb);
        },

        getDomainToken: (domainName, uuid, cb) => {
            return db
                .collection(domainCollectionName)
                .aggregate([
                    {$match: {"name": domainName, "tokens": {$elemMatch: {"uuid": uuid, "enabled": true}}}},
                    {$unwind: "$tokens"},
                    {$match: {"name": domainName, "tokens.uuid": uuid, "tokens.enabled": true}},
                    {
                        $lookup:
                        {
                            from: aclCollectionName,
                            localField: "tokens.roleName",
                            foreignField: "roles",
                            as: "acl"
                        }
                    },
                    {$unwind: "$acl"},
                    {
                        $addFields: {
                            "aclList": "$acl.allows",
                            "roleName": "$tokens.roleName"
                        }
                    },
                    {$project: {"aclList": 1, "roleName": 1}}
                ], cb)
        }
    }
}