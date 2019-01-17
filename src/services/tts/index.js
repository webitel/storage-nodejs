/**
 * Created by igor on 16.08.16.
 */

"use strict";

const https = require('https'),
    log = require(__appRoot + '/lib/log')(module),
    aws = require('./aws4'),
    conf = require(__appRoot + '/conf'),
    defSettings = conf.get('tts'),
    defProviders = new Map()
    ;

let defProviderName = null;

if (defSettings) {
    for (let key in defSettings) {
        if (key === 'defaultProvider') {
            defProviderName = defSettings[key]
        } else {
            defProviders.set(key, defSettings[key])
        }
    }
}

if (defProviderName && !defProviders.has(defProviderName)) {
    throw `No settings default provider ${defProviderName}`
}

function getDefaultSetting(providerName) {
    const provider = providerName || defProviderName;
    if (!defProviders.has(provider))
        return {};

    return defProviders.get(provider);
}
    
module.exports = (req, res, next) => {
    let provider = req.params.provider === 'default' ? defProviderName : req.params.provider;
    if (PROVIDERS.hasOwnProperty(provider) && PROVIDERS[provider] instanceof Function) {
        return PROVIDERS[provider](req, res, next);
    } else {
        log.warn(`Bad provider name ${provider}`);
        res.status(400).end();
    }
};

const PROVIDERS = {
    ivona: (req, res, next) => {
        const ivonaDefaultSettings = getDefaultSetting('ivona');
        const accessKeyId = req.query.accessKey || ivonaDefaultSettings.accessKey;
        const secretAccessKey = req.query.accessToken || ivonaDefaultSettings.accessToken;

        if (!accessKeyId || !secretAccessKey) {
            log.error(`Polly bad request accessKey or accessToken is required`);
            return res.status(400).send(`Ivona bad request accessKey or accessToken is required`);
        }

        let voiceSettings = {
            Input: {
                Data : req.query.text,
                Type : 'text/plain'
            },
            OutputFormat: {
                Codec : req.query.format === '.wav' ? 'OGG' : 'MP3',
                SampleRate : 22050
            },
            Parameters: {
                Rate : 'medium',
                Volume : 'medium',
                SentenceBreak : 500,
                ParagraphBreak : 650
            },
            Voice: {
                Name : req.query.name || 'Salli',
                Language : req.query.language || 'en-US',
                Gender : req.query.gender || 'Female'
            }
        };

        let requestParams = {
            path: `/CreateSpeech`,
            host: 'tts.eu-west-1.ivonacloud.com',
            service: 'tts',
            method: 'POST',
            region: 'eu-west-1',
            body: JSON.stringify(voiceSettings)
        };

        aws.sign(requestParams, {accessKeyId, secretAccessKey});

        return _sendRequest(requestParams, res);
    },
    polly: (req, res, next) => {

        const polyDefaultSettings = getDefaultSetting('polly');
        const accessKeyId = req.query.accessKey || polyDefaultSettings.accessKey;
        const secretAccessKey = req.query.accessToken || polyDefaultSettings.accessToken;
        const voice = req.query.voice || polyDefaultSettings.voice || 'Salli';

        if (!accessKeyId || !secretAccessKey) {
            log.error(`Polly bad request accessKey or accessToken is required`);
            return res.status(400).send(`Polly bad request accessKey or accessToken is required`);
        }

        let voiceSettings = {
            Text: req.query.text,
            TextType: req.query.textType || "text",
            OutputFormat: req.query.format === '.wav' ? 'ogg_vorbis' : 'mp3', // wav
            SampleRate: req.query.rate || "8000", // 8KHz
            VoiceId: voice
        };

        let requestParams = {
            path: `/v1/speech`,
            host: 'polly.eu-west-1.amazonaws.com',
            service: 'polly',
            method: 'POST',
            region: 'eu-west-1',
            body: JSON.stringify(voiceSettings)
        };

        aws.sign(requestParams, {accessKeyId, secretAccessKey});

        return _sendRequest(requestParams, res);
    },
    microsoft: (req, res, next) => {

        const microsoftDefaultSettings = getDefaultSetting('microsoft');
        const keyId = req.query.accessKey || microsoftDefaultSettings.accessKey;
        const keySecret = req.query.accessToken || microsoftDefaultSettings.accessToken;
        const region = req.query.region || microsoftDefaultSettings.region;

        if (!keyId || !keySecret || !region) {
            log.error(`Microsoft bad request accessKey, region or accessToken is required`);
            return res.status(400).send(`Microsoft bad request accessKey, region or accessToken is required`);
        }

        microsoftAccessToken(keyId, keySecret, region, (err, token) => {
            if (err || !token )
                return res.status(500).send('Bad response');

            let voice = {
                gender: req.query.gender || 'Female',
                lang: req.query.language || 'en-US'
            };
            const body =  new Buffer(`<speak version='1.0' xml:lang='${voice.lang}'>
                        <voice xml:lang='${voice.lang}' xml:gender='${voice.gender}' name='Microsoft Server Speech Text to Speech Voice (${microsoftLocalesNameMaping(voice.lang, voice.gender)})'>${req.query.text}
                        </voice>
                      </speak>`);
            let requestParams = {
                path: `/cognitiveservices/v1`,
                host: `${region}.tts.speech.microsoft.com`,
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/ssml+xml',
                    'X-Microsoft-OutputFormat': req.query.format === '.wav' ? 'riff-8khz-8bit-mono-mulaw' : 'audio-16khz-32kbitrate-mono-mp3',

                    // 'X-Search-AppId': appId,
                    // 'X-Search-ClientID': keyId,
                    'User-Agent': 'WebitelACR'
                },
                body

            };
            requestParams.headers["Content-Length"] = body.length;

            return _sendRequest(requestParams, res);

        });
    }
};

function _sendRequest(requestParams, res) {
    let request = https.request(requestParams, (responseTTS) => {
        return _handleResponseTTS(responseTTS, res);
    });
    request.on('error', (error) => {
        log.error(error);
    });

    if (requestParams.body) request.write(requestParams.body);
    request.end();
}

function _handleResponseTTS(responseTTS, res) {
    log.trace(`response TTS status code: ${responseTTS.statusCode}`);
    responseTTS.on('error', (error) => {
        log.error(error);
    });

    if (responseTTS.statusCode !== 200) {
        log.error(responseTTS.headers);
        let d = '';
        responseTTS.on('data', c => d+=c);
        responseTTS.on('end', () => {
            log.error(d);
        });

        return res.status(responseTTS.statusCode).send('Bad response');
    }
    responseTTS.pause();
    res.writeHeader(responseTTS.statusCode, responseTTS.headers);
    responseTTS.pipe(res);
    responseTTS.resume();
}

function microsoftAccessToken(clientId, clientSecret, region, cb) {
    let requestParams = {
        path: `/sts/v1.0/issueToken`,
        host: `${region}.api.cognitive.microsoft.com`,
        method: 'POST',
        headers: {
            'Content-Type': 'text/plain',
            'Ocp-Apim-Subscription-Key': clientId,
        }
    };
    let data = '';

    let request = https.request(requestParams, (res) => {
        log.trace(`response microsoft auth status code: ${res.statusCode} ${res.statusMessage}`);

        res.on('data', function(chunk) {
            data += chunk;
        });

        res.on('end', function() {
            try {
                cb(null, data);
            } catch (e) {
                log.error(e);
                cb(e);
            }
        });

        res.once('error', cb);

    });
    request.once('error', cb);

    if (requestParams.body) request.write(requestParams.body);
    request.end();
}

const GENDER_FEMALE = 'Female';
function microsoftLocalesNameMaping(locale, gender) {
    switch (locale) {
        case 'ar-EG':
            return "ar-EG, Hoda";

        case 'ar-SA':
            return "ar-SA, Naayf";

        case 'ca-ES':
            return "ca-ES, HerenaRUS";

        case 'cs-CZ':
            return "cs-CZ, Vit";

        case 'da-DK':
            return "da-DK, HelleRUS";

        case 'de-AT':
            return "de-AT, Michael";

        case 'de-CH':
            return "de-CH, Karsten";

        case 'de-DE':
            if (isFemale(gender))
                return "de-DE, Hedda";
            else return "de-DE, Stefan, Apollo";

        case 'el-GR':
            return "el-GR, Stefanos";

        case 'en-AU':
            return "en-AU, Catherine";

        case 'id-ID':
            return 'id-ID, Andika';


        case 'en-CA':
            return "en-CA, Linda";

        case 'en-GB':
            if (isFemale(gender))
                return "en-GB, Susan, Apollo";
            else return "en-GB, George, Apollo";

        case 'en-IE':
            return "en-IE, Shaun";

        case 'en-IN':
            if (isFemale(gender))
                return "en-IN, Heera, Apollo";
            else return "en-IN, Ravi, Apollo";

        case 'en-US':
            if (isFemale(gender))
                return "en-US, ZiraRUS";
            else return "en-US, BenjaminRUS";

        case 'es-ES':
            if (isFemale(gender))
                return "es-ES, Laura, Apollo";
            else return "es-ES, Pablo, Apollo";

        case 'es-MX':
            if (isFemale(gender))
                return "es-MX, HildaRUS";
            else return "es-MX, Raul, Apollo";

        case 'fi-FI':
            return "fi-FI, HeidiRUS";

        case 'fr-CA':
            if (isFemale(gender))
                return "fr-CA, Caroline";
            else return "fr-CH, Guillaume";

        case 'fr-FR':
            if (isFemale(gender))
                return "fr-FR, Julie, Apollo";
            else return "fr-FR, Paul, Apollo";

        case 'he-IL':
            return "he-IL, Asaf";

        case 'hi-IN':
            if (isFemale(gender))
                return "hi-IN, Kalpana, Apollo";
            else return "hi-IN, Hemant";

        case 'hu-HU':
            return "hu-HU, Szabolcs";

        case 'it-IT':
            return "it-IT, Cosimo, Apollo";

        case 'ja-JP':
            if (isFemale(gender))
                return "ja-JP, Ayumi, Apollo";
            else return "ja-JP, Ichiro, Apollo";

        case 'ko-KR':
            return "ko-KR, HeamiRUS";

        case 'nb-NO':
            return "nb-NO, HuldaRUS";

        case 'nl-NL':
            return "nl-NL, HannaRUS";

        case 'pl-PL':
            return "pl-PL, PaulinaRUS";

        case 'pt-BR':
            if (isFemale(gender))
                return "pt-BR, HeloisaRUS";
            else return "pt-BR, Daniel, Apollo";

        case 'pt-PT':
            return "pt-PT, HeliaRUS";

        case 'ro-RO':
            return "ro-RO, Andrei";

        case 'ru-RU':
            if (isFemale(gender))
                return "ru-RU, Irina, Apollo";
            else return "ru-RU, Pavel, Apollo";

        case 'sk-SK':
            return "sk-SK, Filip";

        case 'sv-SE':
            return "sv-SE, HedvigRUS";

        case 'th-TH':
            return "th-TH, Pattara";

        case 'tr-TR':
            return "tr-TR, SedaRUS";

        case 'zh-CN':
            if (isFemale(gender))
                return "zh-CN, Yaoyao, Apollo";
            else return "zh-CN, Kangkang, Apollo";

        case 'zh-HK':
            if (isFemale(gender))
                return "zh-HK, Tracy, Apollo";
            else return "zh-HK, Danny, Apollo";

        case 'zh-TW':
            if (isFemale(gender))
                return "zh-TW, Yating, Apollo";
            else return "zh-TW, Zhiwei, Apollo";

        case 'vi-VN':
            return "vi-VN, An";

        default:
            log.error(`unknown local: ${locale}`);
            return "";

    }
}

function isFemale(gender) {
    return GENDER_FEMALE == gender;
}