/**
 * Created by igor on 26.10.16.
 */

"use strict";

const parser = require('cron-parser'),
    log = require(`${__appRoot}/lib/log`)(module)
    ;
    
class Scheduler {
    constructor (cronFormat, fn) {
        let _timer = null;
        const interval = parser.parseExpression(cronFormat);
        const _c = cronFormat;
        log.info(`Create job: ${fn.name || ''}`);

        (function shed() {
            if (_timer)
                clearTimeout(_timer);

            let n = -1;
            do {
                n = interval.next().getTime() - Date.now()
            } while (n < 0);

            log.trace(`Next exec schedule: ${fn.name || ''} ${n}`);
            _timer = setTimeout( function tick() {
                log.trace(`Exec schedule: ${fn.name || ''} ${_c}`);
                fn.apply(null, [shed]);
            }, n);
        })()
    }
}

module.exports = Scheduler;
