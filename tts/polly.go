package tts

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
	"net/http"
)

type TTSEngine interface {
	GetStream(TTSParams) (io.ReadCloser, string, error)
}

type TTSParams struct {
	Key, Token     string
	Format         string
	Voice          string
	Region         string
	Text, TextType string
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

func Microsoft(req TTSParams) (io.ReadCloser, *string, error) {
	microsoftToken(req.Key, req.Token, req.Region)

	return nil, nil, errors.New("TODO")
}

func microsoftToken(key, secret, region string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken", region), nil)
	if err != nil {
		return
	}

	req.Header.Set("Context-Type", "application/ssml+xml")
	req.Header.Set("User-Agent", "WebitelACR")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

}
