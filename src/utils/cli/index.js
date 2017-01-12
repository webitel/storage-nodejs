/**
 * Created by igor on 15.11.16.
 */

"use strict";

const MongoClient = require("mongodb").MongoClient,
    request = require('request'),
    async = require('async'),
    ProgressBar = require('progress'),
    ObjectId = require('mongodb').ObjectId,
    path = require('path')
    ;

//mongoimport --host 10.10.10.25:27017 --db webitel --collection location --file mongo_location.json

global.__appRoot = path.normalize(`${__dirname}/../..`);

const conf = require(`${__appRoot}/config`),
    recServices = require(`${__appRoot}/services/recordings`);

const args = process.argv.slice(2);
const fn = getCommands();
const config = parseArgs(args);

if (!(fn[args[0]] instanceof Function))
    throw `Bad function ${args[0]}`;

fn[args[0]](config);

function getCommands () {
    return {
        file: (config = {}) => {
            let {filter = {}, url, collection, mongo, lastInsertedId} = config;
            if (!mongo)
                throw `Bad mongo uri`;
            if (!collection)
                throw `Bad collection name`;

            if (!url)
                throw `Bad url`;

            async.waterfall(
                [
                    cb => createMongoInstance(mongo, {}, cb),
                    (mongoDb, cb) => {
                        const c = mongoDb.collection(collection);
                        c.count(filter || {}, {}, (e, count) => {
                            if (e)
                                return cb(e);

                            console.log(`Found ${collection} documents: ${count}`);
                            if (count === 0)
                                return cb(new Error(`Not found data`));
                            cb(null, mongoDb, count);
                        })
                    },
                    (mongoDb, count, cb) => {
                        const c = mongoDb.collection(collection);
                        const stream = c.find(filter || {})
                            .sort({"_id": 1})
                            .addCursorFlag('noCursorTimeout',true)
                            .stream();

                        console.log(`Start: ${new Date().toLocaleString()}`);

                        const bar = new ProgressBar(' exports [:bar] :percent :etas  failed :failedCount err :lastError', {
                            complete: '=',
                            incomplete: ' ',
                            width: 100,
                            total: count
                        });

                        let inserted = 0;
                        let failedCount = 0;
                        let lastError = "";
                        const tick = (err) => {
                            if (err) {
                                failedCount++;
                                lastError = err.message.substring(0, 20);
                            }

                            bar.tick({
                                failedCount: failedCount,
                                lastError: lastError
                            });

                            if (inserted++ > count)
                                return cb();
                            stream.resume();
                            // console.log(inserted);
                        };

                        const makeUrl = (doc) => {
                            return `${url}?domain=${doc.domain}&id=${doc.uuid}&name=${doc.name}&type=${getTypeFromContentType(doc['content-type'])}`
                        };

                        stream.on('data', doc => {
                            stream.pause();
                            lastInsertedId = doc._id;
                            recServices._getFile(doc, {}, (err, res) => {
                                // inserted++;
                                if (err) {
                                    // console.log(`Skip file: ${doc.uuid}, error ${err.message}`)
                                    tick(err);

                                } else {
                                    if (res.location) {
                                        request
                                            .get(res.location)
                                            .pipe(request.post(makeUrl(doc)), (e, v) => {
                                                console.log(v.statusCode);
                                                tick();
                                            })
                                    } else {
                                        res
                                            .source
                                            .pipe(request.post(makeUrl(doc), (e, v) => {
                                                if (e) {
                                                    cb(e);
                                                } else if (v.statusCode !== 200) {
                                                    tick(new Error(`Bad response status code ${v.statusCode}`));
                                                } else {
                                                    tick();
                                                }
                                            }))
                                    }
                                }

                            });
                        });

                        stream.on('error', e => {
                            throw e;
                        });

                        // stream.once('end', cb)
                    }
                ],
                (err, res) => {
                    if (err) {
                        console.error(err);
                        if (lastInsertedId) {
                            config.lastInsertedId = lastInsertedId;
                            config.filter = {_id: {$gte: ObjectId(lastInsertedId)}};
                            console.log(`Try recovery 5sec`);
                            setTimeout(() => {
                                fn.file(config);
                            }, 5000);
                            return;
                        }
                    }


                    console.log(`End: ${new Date().toLocaleString()}\n`);
                    process.exit(0)
                }

            );
        },
        cdr: (config = {}) => {
            let {filter = {}, url, collection, mongo, lastInsertedId} = config;
            if (!mongo)
                throw `Bad mongo uri`;
            if (!collection)
                throw `Bad collection name`;

            if (!url)
                throw `Bad url`;

            async.waterfall(
                [
                    cb => createMongoInstance(mongo, {}, cb),
                    (mongoDb, cb) => {
                        const c = mongoDb.collection(collection);
                        c.count(filter, {}, (e, count) => {
                            if (e)
                                return cb(e);

                            console.log(`Found ${collection} documents: ${count}`);
                            cb(null, mongoDb, count);
                        })
                    },
                    (mongoDb, count, cb) => {
                        const c = mongoDb.collection(collection);
                            // stream = c.find(filter)
                            //     .sort({"_id": 1})
                            //     .stream();

                        console.log(`Start: ${new Date().toLocaleString()}`);

                        const stream = c.aggregate([
                            {$match: filter || {}},
                            {
                                $sort : { "_id": 1}
                            },
                            {
                                $lookup: {
                                    from: conf.get('mongodb:collectionFile'),
                                    localField: "variables.uuid",
                                    foreignField: "uuid",
                                    as: "recordings"
                                }
                            },
                            {
                                $project: {"core-uuid": 1, "_id": 1, "switchname": 1, "variables": 1, "callflow": 1}
                            }
                        ], {allowDiskUse: true}).stream();

                        console.log(`Load stream data...`);
                        const bar = new ProgressBar(' exports [:bar] :percent :etas', {
                            complete: '=',
                            incomplete: ' ',
                            width: 100,
                            total: count
                        });

                        let inserted = 0;
                        stream.on('data', doc => {
                            stream.pause();
                            lastInsertedId = doc._id;
                            request({
                                method: "POST",
                                url: url,
                                useQuerystring: true,
                                qs: {
                                    skip_mongo: true
                                },
                                json: doc
                            }, (e, r) => {
                                if (e)
                                    return cb(e);

                                if (r.statusCode !== 200)
                                    return cb(new Error(`Bad response status code ${r.statusCode}`));

                                bar.tick();
                                if (++inserted >= count)
                                    return cb();
                                stream.resume();
                            });
                        });

                        stream.on('start', e => console.log('Ok data'));

                        stream.on('error', e => {
                            throw e;
                        });

                        // stream.once('end', cb)
                    }
                ],
                (err, res) => {
                    if (err) {
                        console.error(err.message);
                        if (lastInsertedId) {
                            config.lastInsertedId = lastInsertedId;
                            config.filter = {_id: {$gte: ObjectId(lastInsertedId)}};
                            console.log(`Try recovery 5sec`);
                            setTimeout(() => {
                                fn.cdr(config);
                            }, 5000);
                            return;
                        }
                    }


                    console.log(`End: ${new Date().toLocaleString()}\n`);
                    process.exit(0)
                }

            )
        },
        location: (config = {}) => {
            let {collection, mongo} = config;
            if (!mongo)
                throw `Bad mongo uri`;
            if (!collection)
                throw `Bad collection name`;

            console.log(`Start export location to collection ${collection}`);

            createMongoInstance(mongo, {}, (err, db) => {
                if (err) {
                    throw err;
                }
                const c = db.collection(collection);

                const location = require("./location.js");
                async.each(
                    location,
                    (item, cb) => {
                        console.log(`Inserting ${item.country} code ${item.code}`);
                        item._id = ObjectId(item._id);
                        c.insert(item, cb)
                    },
                    err => {
                        if (err) {
                            throw err;
                        }
                        process.exit(1);
                    }
                )
            });
        }
    }
}

function parseArgs(args = []) {
    const parsed = {
        url: "",
        collection: null,
        filter: null,
        elastic: null,
        mongo: conf.get('mongodb:uri'),
        index: "cdr"
    };

    if (args[0] === 'file') {
        parsed.collection = conf.get('mongodb:collectionFile')
    } else if (args[0] === 'cdr') {
        parsed.collection = conf.get('mongodb:collectionCDR')
    } else if (args[0] === 'location') {
        parsed.collection = conf.get('mongodb:collectionLocation')
    }

    let i = 0;

    while (i < args.length) {
        switch  (args[i]) {
            case "-u":
                parsed.url = args[++i];
                break;

            case "-c":
                parsed.collection = args[++i];
                break;

            case "-i":
                parsed.index = args[++i];
                break;

            case "-e":
                parsed.elastic= args[++i];
                break;

            case "-m":
                parsed.mongo = args[++i];
                break;

            case "-f":
                try {
                    parsed.filter = JSON.parse(args[++i]);
                } catch (e) {
                    console.error(`Filter must JSON string`);
                    throw e;
                }
                break;

            case "--help":
                console.log(`
<commands> [params];
params:
    -m: MongoDB URI;
                `);
                return process.exit(0);
                
            default:
                i++;
        }
    }
    return parsed;
}

function createMongoInstance(uri, option = {}, cb) {
    const mongoClient = new MongoClient();
    mongoClient.connect(uri, option, cb);
}

const getTypeFromContentType = (contentType = "") => {
    switch (contentType) {
        case 'application/pdf':
            return 'pdf';
        case 'audio/wav':
            return 'wav';
        case 'audio/mpeg':
        default:
            return 'mp3'
    }
};