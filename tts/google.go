package tts

import (
	"bytes"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"io"
	"io/ioutil"
)

func Google(params TTSParams) (io.ReadCloser, *string, error) {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: params.Language,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
			//SpeakingRate:     0,
			//Pitch:            0,
			//VolumeGainDb:     0,
			SampleRateHertz: 8000,
			//EffectsProfileId: nil,
		},
	}

	if params.SpeakingRate != 0 {
		req.AudioConfig.SpeakingRate = params.SpeakingRate
	}

	if params.Pitch != 0 {
		req.AudioConfig.Pitch = params.Pitch
	}

	if params.VolumeGainDb != 0 {
		req.AudioConfig.VolumeGainDb = params.VolumeGainDb
	}

	if params.EffectsProfileId != nil {
		req.AudioConfig.EffectsProfileId = params.EffectsProfileId
	}

	switch params.Voice {
	case "MALE":
		req.Voice.SsmlGender = texttospeechpb.SsmlVoiceGender_MALE
	case "FEMALE":
		req.Voice.SsmlGender = texttospeechpb.SsmlVoiceGender_FEMALE
	default:
		req.Voice.SsmlGender = texttospeechpb.SsmlVoiceGender_NEUTRAL
	}

	v := "audio/ogg"
	if params.Format == "mp3" {
		v = "audio/mp3"
		req.AudioConfig.SampleRateHertz = 22050
		req.AudioConfig.AudioEncoding = texttospeechpb.AudioEncoding_MP3
	}

	// Set the text input to be synthesized.
	if params.TextType == "ssml" {
		req.Input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: params.Text},
		}
	} else {
		req.Input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: params.Text},
		}
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return nil, nil, err
	}

	r := ioutil.NopCloser(bytes.NewReader(resp.GetAudioContent()))

	return r, &v, nil
}
