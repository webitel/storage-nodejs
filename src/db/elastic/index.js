/**
 * Created by igor on 30.08.16.
 */

"use strict";

const elasticsearch = require('elasticsearch'),
    EventEmitter2 = require('eventemitter2').EventEmitter2,
    log = require(__appRoot + '/lib/log')(module),
    setCustomAttribute = require(__appRoot + '/utils/cdr').setCustomAttribute,
    async = require('async')
;

class ElasticClient extends EventEmitter2 {
    constructor (config) {
        super();
        this.connected = false;
        this.config = config;

        this.client = new elasticsearch.Client({
            host: config.host,
            keepAlive: config.keepAlive || true,
            requestTimeout: config.requestTimeout || 500000
        });

        this.ping();
        this.pingId = null;
    }

    ping () {
        if (this.pingId)
            clearTimeout(this.pingId);

        this.client.ping({}, (err) => {
            if (err) {
                log.error(err);
                this.pingId = setTimeout(this.ping.bind(this), 1000);
                return null;
            }

            log.info(`Connect to elastic - OK`);
            this.connected = true;
            this.initTemplate();
            this.emit('elastic:connect', this);
        })
    }

    initTemplate () {
        const client = this.client;

        client.indices.getTemplate((err, res) => {
            if (err) {
                return log.error(err);
            }

            let elasticTemplatesNames = Object.keys(res),
                templates = this.config.templates || [],
                tasks = [],
                delTemplate = [];

            templates.forEach(function (template) {
                if (elasticTemplatesNames.indexOf(template.name) > -1) {
                    delTemplate.push(function (done) {
                        client.indices.deleteTemplate(
                            template,
                            function (err) {
                                if (err) {
                                    log.error(err);
                                } else {
                                    log.debug('Template %s deleted.', template.name)
                                }
                                done();
                            }
                        );
                    });
                }

                tasks.push(function (done) {
                    client.indices.putTemplate(
                        template,
                        function (err) {
                            if (err) {
                                log.error(err);
                            } else {
                                log.debug('Template %s - created.', template.name);
                            }
                            done();
                        }
                    );
                });
            });

            if (tasks.length > 0) {
                async.waterfall([].concat(delTemplate, tasks) , (err) => {
                    if (err)
                        return log.error(err);
                    return log.info(`Replace elastic template - OK`);
                });
            }
        });
    }

    insertCdr (doc, cb) {
        let currentDate = new Date(),
            indexName = `cdr-${currentDate.getMonth() + 1}.${currentDate.getFullYear()}`,
            _record = setCustomAttribute(doc),
            _id = _record._id.toString();
        delete _record._id;

        this.client.create({
            index: (indexName + (doc.variables.domain_name ? '-' + doc.variables.domain_name : '')).toLowerCase(),
            type: 'collection',
            id: _id,
            body: _record
        }, cb);
    }
}
    
module.exports = ElasticClient;