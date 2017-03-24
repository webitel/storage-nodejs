/**
 * Created by igor on 23.03.17.
 */

"use strict";

const http = require('http'),
    rpc = require(__appRoot + '/lib/broker')(),
    conf = require(__appRoot + '/conf'),
    log = require(__appRoot + '/lib/log')(module),
    DEF_KEY = conf.get('stt:defaultKey')
    ;

module.exports = {
    addRoutes: addRoutes
};

function addRoutes(api) {
    api.put('/sys/stt', sttStream);
}

function sendResponse(res, args, queueName) {
    if (rpc && queueName && args) {
        rpc.sendToQueue(queueName, {"exec-api": "channel.sttResponse", "exec-args": args}, {}, e => {
            if (e)
                return log.error(e);
        });
    }

    return res.send('OK');
}

function sttStream(req, res, next) {
    const {
        callId,
        reply,
        lang = 'en-US',
        setVar,
        key = DEF_KEY,
        codec = 'audio/l16',
        rate
    } = req.query;

    if (!callId || !reply || !rate) {
        log.error(`Bad request: callId, rate or reply is required, url: ${req.originalUrl}`);
        return sendResponse(res)
    }

    let requestParams = {
        path: `/speech-api/v2/recognize?client=chromium&output=json&lang=${lang}&key=${key || DEF_KEY}`,
        host: 'www.google.com',
        method: 'POST',
        headers: {
            'Content-Type': `${codec}; rate=${rate}`
        }
    };

    const r = http.request(requestParams, (resGoog) => {
        log.debug(`Response for call ${callId}: ${resGoog.statusCode} `);

        let resData = '';
        resGoog.on('data', c => resData += c);

        resGoog.on('end', () => {
            const text = /.*\n(.*)/.exec(resData);
            let stt = null;
            if (text && text[1]) {
                stt = parseJson(text[1]);

            }
            log.trace(`Call ${callId} stt: ${resData}`);
            return sendResponse(res, {stt, callId, setVar}, reply);
        });

        resGoog.on('error', e => {
            log.error(e);
            return sendResponse(res, {stt: null, callId, setVar}, reply);
        });

    });

    r.on('error', e => {
        log.error(e);
    });

    req.pipe(r)
}

function parseJson(str) {
    try {
        return JSON.parse(str);
    } catch (e) {
        return null;
    }
}