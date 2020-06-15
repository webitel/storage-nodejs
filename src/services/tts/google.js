const textToSpeech = require('@google-cloud/text-to-speech');
const client = new textToSpeech.TextToSpeechClient();
const Duplex = require('stream').Duplex;

module.exports = async (req, res, next) => {
    // Construct the request
    const rate = +(req.query.rate || 8000);
    const request = {
        input: {},
        // Select the language and SSML Voice Gender (optional)
        voice: {
            languageCode: req.query.language || 'en-US',
            ssmlGender: (req.query.gender || 'NEUTRAL').toUpperCase()
        },
        // Select the type of audio encoding
        audioConfig: {
            sampleRateHertz: rate,
            audioEncoding: getCodec(rate)
        },
    };

    if (req.query.audioEncoding) {
        request.audioConfig.audioEncoding = req.query.sampleRateHertz;
    }
    if (+req.query.sampleRateHertz) {
        request.audioConfig.sampleRateHertz = +req.query.sampleRateHertz;
    }
    if (+req.query.speakingRate) {
        request.audioConfig.speakingRate = +req.query.speakingRate;
    }
    if (+req.query.pitch) {
        request.audioConfig.pitch = +req.query.pitch;
    }
    if (+req.query.volumeGainDb) {
        request.audioConfig.volumeGainDb = +req.query.volumeGainDb;
    }
    if (req.query.effectsProfileId) {
        request.audioConfig.effectsProfileId = (req.query.effectsProfileId || '').split(',');
    }


    if (req.query.textType === 'ssml') {
        request.input.ssml = req.query.text
    } else {
        request.input.text = req.query.text;
    }

    // Performs the Text-to-Speech request
    try {
        const [response] = await client.synthesizeSpeech(request);
        res.setHeader("Content-Type", getContentType(rate));
        let stream = new Duplex();
        stream.push(response.audioContent);
        stream.push(null);
        res.status(200);
        stream.pipe(res);
    } catch (e) {
        return next(e);
    }
};

const getCodec = (rate = 8000) => {
    if (rate > 16000) {
        return "MP3"
    }
    return "LINEAR16"
};

function getContentType(rate = 8000) {
    if (rate > 16000) {
        return "audio/mpeg"
    }
    return "audio/ogg"
}
