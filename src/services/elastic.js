/**
 * Created by igor on 30.08.16.
 */

"use strict";

const checkPermission = require(__appRoot + '/utils/acl'),
    CodeError = require(__appRoot + '/lib/error')
    ;
    
const Service = module.exports = {
    search (caller, option, cb) {
        let acl = caller.acl,
            _ro = false
            ;

        let _readAll = checkPermission(acl, 'cdr', 'r');

        if (!_readAll && checkPermission(acl, 'cdr', 'ro', true)) {
            _ro = true;
        }

        if (!_ro && !_readAll) {
            return cb(new CodeError(403, "Permission denied!"));
        }

        let filter = {
            "bool": {
                "must": [

                ],
                "must_not": []
            }
        };

        if (option.filter)
            filter.bool.must.push(option.filter);

        if (_ro)
            filter.bool.must.push({
                "term": {"variables.presence_id": caller.id}
            });

        let columns = option.columns,
            columnsDate = option.columnsDate || [],
            query = option.query || "*",
            limit = parseInt(option.limit, 10) || 40,
            pageNumber = option.pageNumber,
            sort = (option.sort && Object.keys(option.sort).length > 0) ? option.sort : {"Call start time":{"order":"desc","unmapped_type":"boolean"}}
            ;

        application.elastic.search(
            {
                index: `cdr-*${caller.domain ? '-' + caller.domain : '' }`,
                size: limit,
                _source: columns,
                fields: columns,
                ignoreUnavailable: true,
                from: pageNumber > 0 ? ((pageNumber - 1) * limit) : 0, //Number — Starting offset (default: 0)
                body: {
                    "fielddata_fields": columnsDate,
                    "sort": [sort],
                    "query": {
                        "filtered": {
                            "query": {
                                "query_string": {
                                    "analyze_wildcard": true,
                                    //"default_operator": "AND",
                                    "query": query
                                }
                            },
                            "filter": filter
                        }
                    }
                }
            },
            cb
        );
    }
};