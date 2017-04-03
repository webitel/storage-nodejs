/**
 * Created by igor on 30.08.16.
 */

"use strict";

const checkPermission = require(__appRoot + '/utils/acl'),
    CodeError = require(__appRoot + '/lib/error')
    ;
    
const Service = module.exports = {
    search (caller, option = {}, cb) {
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
                    {
                        "query_string": {
                            "analyze_wildcard": true,
                            "query": option.query || "*"
                        }
                    }
                ],
                "must_not": []
            }
        };

        if (option.filter) {
            if (option.filter instanceof Array) {
                filter.bool.must = filter.bool.must.concat(option.filter)
            } else {
                filter.bool.must.push(option.filter);
            }
        }

        if (_ro)
            filter.bool.must.push({
                "term": {"variables.presence_id": caller.id}
            });

        let columns = option.columns,
            columnsDate = option.columnsDate || [],
            scroll = option.scroll,
            limit = parseInt(option.limit, 10) || 40,
            pageNumber = option.pageNumber,
            sort = (option.sort && Object.keys(option.sort).length > 0) ? option.sort : {"Call start time":{"order":"desc","unmapped_type":"boolean"}}
            ;

        const domain = caller.domain || option.domain;
        
        application.elastic.search(
            {
                index: `cdr-*${domain ? '-' + domain : '' }`,
                size: limit,
                storedFields: columns,
                docvalueFields: columns,
                ignoreUnavailable: true,
                scroll,
                from: pageNumber > 0 ? ((pageNumber - 1) * limit) : 0, //Number — Starting offset (default: 0)
                body: {
                    "sort": [sort],
                    "docvalue_fields": columnsDate,
                    "query": filter
                }
            },
            cb
        );
    },

    scroll: (caller, option = {}, cb) => {
        if (!caller)
            return cb(new CodeError(401, ""));

        application.elastic.scroll(
            {
                scroll_id: option.scrollId,
                scroll: option.scroll
            },
            cb
        );
    }
};