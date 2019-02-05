/**
 * Created by igor on 30.08.16.
 */

"use strict";

const elasticsearch = require('elasticsearch'),
    EventEmitter2 = require('eventemitter2').EventEmitter2,
    log = require(__appRoot + '/lib/log')(module),
    setCustomAttribute = require(__appRoot + '/utils/cdr').setCustomAttribute,
    moment = require("moment")
;

const CDR_NAME = 'cdr',
    CDR_TYPE_NAME = 'cdr';


const scriptAddPin = `if (ctx._source["pinned_items"] == null) {
    ctx._source.pinned_items = [params.rec]
} else {
    if (!ctx._source.pinned_items.contains(params.rec)) {
        ctx._source.pinned_items.add(params.rec)
    }
}`;

const scriptDelPin = `if (ctx._source["pinned_items"] != null) {
    for (int i = 0; i < ctx._source.pinned_items.size(); i++){
        if (ctx._source.pinned_items[i] == params.rec) {
            ctx._source.pinned_items.remove(i)
        }    
    }
}`;

const dropTimeFromTimestamp = (time) => {
    return (time - (time % 86400000))
};

const getIndexName = (template, index, domain, leg, startTimestamp) => {
    const time = moment(dropTimeFromTimestamp(startTimestamp)).utc();
    return template.replace(/\$\{([\s\S]*?)\}/g, (pattern) => {
        switch (pattern) {
            case "${INDEX}":
                return index;
            case "${LEG}":
                return leg;
            case "${YEAR}":
                return time.format('YYYY');
            case "${MONTH}":
                return time.format('MM');
            case "${WEEK}":
                return time.format('WW');
            case "${DAY}":
                return time.format('DD');
            case "${DOMAIN}":
                return domain

        }
        return ""
    }).toLowerCase()
};

class ElasticClient extends EventEmitter2 {
    constructor (config) {
        super();
        this.connected = false;
        this.config = config;

        this.client = new elasticsearch.Client(this.config);

        [
            this.indexNameCdr,
            this.indexNameCdrTemplate,
            this.indexNameAccounts,
            this.indexNameAccountsTemplate
        ] = [
            config.indexNameCdr,
            config.indexNameCdrTemplate,
            config.indexNameAccounts,
            config.indexNameAccountsTemplate
        ];

        if (!this.indexNameCdr || !this.indexNameCdrTemplate) {
            throw 'Config indexNameCdrTemplate is required'
        }

        this.ping();
        this.pingId = null;
    }

    ping () {
        if (this.pingId)
            clearTimeout(this.pingId);

        this.client.ping({}, (err) => {
            if (err) {
                log.error(err);
                if (err.statusCode === 403) {
                    throw err
                }
                this.pingId = setTimeout(this.ping.bind(this), 1000);
                return null;
            }

            log.info(`Connect to elastic - OK`);
            this.connected = true;
            //this.initTemplate();
            // this.initMaxResultWindow();

        })
    }

    search (data, cb) {
        return this.client.search(data, cb);
    }

    scroll (params = {}, cb) {
        return this.client.scroll(params, cb);
    }


    getCdrInsertParam (doc, createdAt, skipAttribute) {
        const index = getIndexName(this.indexNameCdrTemplate, this.indexNameCdr
            , ((doc.variables && doc.variables.domain_name) || ''), 'a', createdAt || Date.now());

        let _record = skipAttribute ? doc : setCustomAttribute(doc),
            _id = (_record.variables && _record.variables.uuid) || _record._id.toString();
        delete _record._id;
        return {
            index,
            type: CDR_TYPE_NAME,
            id: _id,
            body: {
                doc: _record,
              //  parent: _id,
                doc_as_upsert: true
            }
        };
    }

    insertPostProcess (doc, createdAt, cb) {
        //TODO
        const r = this.getCdrInsertParam(doc, createdAt, true);
        this.client.update(r, cb);
    }


    addPinCdr (id, index, userId, cb) {
        this.client.update({
            index: index,
            type: CDR_TYPE_NAME,
            id: id,
            body: {
                "script" : {
                    "inline": scriptAddPin,
                    "lang": "painless",
                    "params" : {
                        "rec" : userId
                    }
                }
            }
        }, cb);
    }

    delPinCdr (id, index, userId, cb) {
        this.client.update({
            index: index,
            type: CDR_TYPE_NAME,
            id: id,
            body: {
                "script" : {
                    "inline": scriptDelPin,
                    "lang": "painless",
                    "params" : {
                        "rec" : userId
                    }
                }
            }
        }, cb);
    }


    insertFile (doc, cb) {
        const _record = doc,
            _id = _record && _record.uuid;

        const index = getIndexName(this.indexNameCdrTemplate, this.indexNameCdr, (doc.domain_name || ''), 'a', doc.created_on || Date.now());
        log.trace(`Try save recordings to elastic index ${index} from id ${_id}`);
        this.client.update({
            index,
            type: CDR_TYPE_NAME,
            id: _id,
            body: {
                // parent: _id,
                "script" : {
                    "inline": "if (ctx._source[\"recordings\"] == null) { ctx._source.recordings = params.rec } else {  ctx._source.recordings.add(params.rec[0])}",
                    "lang": "painless",
                    "params" : {
                        "rec" : doc.recordings instanceof Array ? doc.recordings : [doc.recordings]
                    }
                },
                upsert: _record
            }
        }, cb);
    }

    //todo
    removeFile (uuid = "", _id = "", domain = "", cb) {
        const indexName = `${CDR_NAME}*`;

        this.client.updateByQuery({
            index: (indexName + (domain ? '-' + domain : '')).toLowerCase(),
            type: CDR_TYPE_NAME,
            body: {
                "query": {
                    "bool": {
                        "must": [
                            {
                                "term": {
                                    "recordings._id": _id.toString()
                                }
                            },
                            {
                                "term": {
                                    "uuid": uuid
                                }
                            }
                        ],
                        "must_not": []
                    }
                },
                "script" : {
                    "inline": "if (ctx._source[\"recordings\"] != null) { for (int i = 0; i < ctx._source.recordings.size(); i++) { if (ctx._source.recordings[i]._id == params.fileId) {ctx._source.recordings.remove(i); break;}} }",
                    "lang": "painless",
                    "params" : {
                        "fileId" : _id
                    }
                }
            }
        }, cb);
    }

    //todo
    findByUuid (uuid, domain, cb) {
        this.client.search(
            {
                index: `${CDR_NAME}-a-*${domain ? '-' + domain : '' }`,
                size: 1,
                _source: false,
                body: {
                    "query": {
                        "term": {
                            "uuid": uuid
                        }
                    }
                }
            },
            cb
        )
    }
    /*
     * TODO add UUID ??
     */
    findRecFromHash (hash, domain, cb) {
        this.client.search(
            {
                index: `${CDR_NAME}-a-*${domain ? '-' + domain : '' }`,
                size: 1,
                _source: ["recordings.*"],
                body: {
                    "query": {
                        "constant_score" : {
                            "filter" : {
                                "bool" : {
                                    "must" : [
                                        { "term" : { "recordings.hash" : hash } }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            (err, res) => {
                if (err)
                    return cb(err);

                const data = res && res.hits && res.hits.hits;
                if (!data || data.length !== 1)
                    return cb();

                for (let rec of data[0]._source.recordings) {
                    if (rec.hash === hash)
                        return cb(null, rec);
                }

                return cb();
            }
        )
    }

    //todo
    removeCdr (id, indexName = "", cb) {
        this.client.deleteByQuery({
            index: 'cdr*',
            body: {
                query: {
                    bool: {
                        should: [
                            {
                                term: {uuid: id}
                            },
                            {
                                term: {parent_uuid: id}
                            }
                        ]
                    }
                }
            }
        }, function (err, res) {
            return cb(err, res)
        });
    }
}
    
module.exports = ElasticClient;