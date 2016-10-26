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
        
        (function shed() {
            if (_timer)
                clearTimeout(_timer);

            _timer = setTimeout( function tick() {
                log.trace(`Exec schedule ${_c}`);
                fn.apply(null, [shed]);
            }, interval.next().getTime() - Date.now());
        })()
    }
}

module.exports = Scheduler;
