/**
 * Created by igor on 23.03.17.
 */

"use strict";

const http = require('https'),
    rpc = require(__appRoot + '/lib/broker')(),
    conf = require(__appRoot + '/conf'),
    log = require(__appRoot + '/lib/log')(module),
    DEF_KEY = conf.get('stt:defaultKey'),
    {Encode} = require('base64-stream')
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

function sttStream(req, res) {
    const {
        callId,
        reply,
        lang = 'en-US',
        setVar,
        key = DEF_KEY,
        rate
    } = req.query;

    if (!callId || !reply || !rate) {
        log.error(`Bad request: callId, rate or reply is required, url: ${req.originalUrl}`);
        return sendResponse(res)
    }

    const body = JSON.stringify({
        "config": {
            "encoding": "LINEAR16", //TODO
            "sampleRateHertz": +rate,
            "languageCode": lang,
        },
        "audio": {
            "content": "@DATA@"
        },
    });

    let requestParams = {
        path: `/v1/speech:recognize?key=${key || DEF_KEY}`,
        host: 'speech.googleapis.com',
        method: 'POST',
        headers: {
            'Transfer-Encoding': 'chunked',
            'Content-Type': 'application/json'
        }
    };

    const r = http.request(requestParams, (resGoog) => {
        log.debug(`Response for call ${callId}: ${resGoog.statusCode} `);

        let resData = '';
        resGoog.on('data', c => resData += c);

        resGoog.on('end', () => {
            const stt = parseJson(resData);
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
    r.write(body.substring(0, body.indexOf('@DATA@')));

    const encoder = new Encode();
    encoder.on('data', c => {
        r.write(c.toString())
    });

    req.on('data', c => {
        encoder.write(c);
    });

    req.on('end', ()=> {
        r.write(body.substring(body.indexOf('@DATA@') + 6));
        r.end();
    });
}

function parseJson(str) {
    try {
        return JSON.parse(str);
    } catch (e) {
        return null;
    }
}