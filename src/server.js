/**
 * Created by igor on 23.08.16.
 */

"use strict";

// Include the cluster module
var cluster = require('cluster'),
    crashCount = 1;

// Code to run if we're in the master process
if (cluster.isMaster) {
    cluster.fork();

    // Listen for dying workers
    cluster.on('exit', function (worker) {

        // Replace the dead worker, we're not sentiment
        console.error('Worker ' + worker.id + ' died.');
        cluster.fork({
            "CRASH_WORKER_COUNT": (crashCount++)
        });
    });

// Code to run if we're in a worker process
} else {
    require('./worker');
    console.info('Worker ' + cluster.worker.id + ' running!');
}