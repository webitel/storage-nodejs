package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
	"io"
	"net/http"
	"path"
)

type S3FileBackend struct {
	BaseFileBackend
	name        string
	region      string
	accessKey   string
	accessToken string
	bucket      string
	endpoint    string
	pathPattern string
	sess        *session.Session
	svc         *s3.S3
	uploader    *s3manager.Uploader
}

func (self *S3FileBackend) Name() string {
	return self.name
}

func (self *S3FileBackend) GetStoreDirectory(domain int64) string {
	return path.Join(parseStorePattern(self.pathPattern, domain))
}

func (self *S3FileBackend) getEndpoint() *string {
	if self.endpoint == "amazonaws.com" {
		return nil
	} else if self.region != "" {
		return aws.String(fmt.Sprintf("%s.%s", self.region, self.endpoint))
	} else {
		return aws.String(fmt.Sprintf("%s", self.endpoint))
	}
}

func (self *S3FileBackend) TestConnection() *model.AppError {
	config := &aws.Config{
		Region:   aws.String(self.region),
		Endpoint: self.getEndpoint(),
		//DisableSSL: aws.Bool(true),
		//S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(self.accessKey, self.accessToken, ""),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		panic(err.Error())
	}

	self.sess = sess
	self.svc = s3.New(sess)
	self.uploader = s3manager.NewUploader(sess)

	if _, err := self.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(self.bucket),
	}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != s3.ErrCodeBucketAlreadyOwnedByYou && aerr.Code() != s3.ErrCodeBucketAlreadyExists {
				return model.NewAppError("S3FileBackend", "utils.file.s3.test_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		} else {
			return model.NewAppError("S3FileBackend", "utils.file.s3.test_connection.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}

func (self *S3FileBackend) Write(src io.Reader, file File) (int64, *model.AppError) {
	directory := self.GetStoreDirectory(file.Domain())
	location := path.Join(directory, file.GetStoreName())

	params := &s3manager.UploadInput{
		Bucket: &self.bucket,
		Key:    aws.String(location),
		Body:   src,
	}

	res, err := self.uploader.Upload(params)

	if err != nil {
		return 0, model.NewAppError("WriteFile", "utils.file.s3.writing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	file.SetPropertyString("location", location)
	wlog.Debug(fmt.Sprintf("[%s] create new file %s", self.name, res.Location))

	return file.GetSize(), nil
}

func (self *S3FileBackend) Remove(file File) *model.AppError {
	return nil
}

func (self *S3FileBackend) RemoveFile(directory, name string) *model.AppError {
	return nil
}

func (self *S3FileBackend) Reader(file File, offset int64) (io.ReadCloser, *model.AppError) {
	var rng *string = nil
	if offset > 0 {
		rng = aws.String(fmt.Sprintf("%d/*", offset))
	}

	params := &s3.GetObjectInput{
		Bucket: &self.bucket,
		Key:    aws.String(file.GetPropertyString("location")),
		Range:  rng,
	}

	out, err := self.svc.GetObject(params)
	if err != nil {

	}

	return out.Body, nil
}
