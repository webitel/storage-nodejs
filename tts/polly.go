package tts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
)

type TTSEngine interface {
	GetStream(TTSParams) (io.ReadCloser, string, error)
}

type TTSParams struct {
	Key, Token     string
	Format         string
	Voice          string
	Region         string
	Language       string
	Text, TextType string

	//google
	SpeakingRate     float64
	Pitch            float64
	VolumeGainDb     float64
	EffectsProfileId []string
}

func Poly(req TTSParams) (io.ReadCloser, *string, error) {
	config := &aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewStaticCredentials(req.Key, req.Token, ""),
	}

	if req.Region != "" {
		config.Region = aws.String(req.Region)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, nil, err
	}

	p := polly.New(sess)
	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(polly.OutputFormatMp3),
		SampleRate:   aws.String("22050"),
		Text:         aws.String(req.Text),
		VoiceId:      aws.String(polly.VoiceIdEmma),
	}

	if req.Format == "ogg" {
		params.SetOutputFormat(polly.OutputFormatOggVorbis)
	} else {
		params.SetOutputFormat(polly.OutputFormatMp3)
	}

	if req.TextType != "" {
		params.TextType = aws.String(req.TextType)
	}

	if out, err := p.SynthesizeSpeech(params); err != nil {
		return nil, nil, err
	} else {
		return out.AudioStream, out.ContentType, nil
	}
}
