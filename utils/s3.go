package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
    "log"
    "mime/multipart"
)

func UploadImage(file multipart.File, header *multipart.FileHeader, userID string) (string, string, error) {
	aws_access_key_id := ""
	aws_secret_access_key := "YkzztVy2f74WA/"
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err := creds.Get()
	if err != nil {
		log.Printf("bad credentials: %s", err)
        return "", "", err
	}

	cfg := aws.NewConfig().WithRegion("us-west-2").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	defer file.Close()

	size, _ := file.Seek(0, 0)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

    bucketName := "yaopeng-photos"
    fileName := header.Filename
	path := "/images/" + userID + "/" + fileName
	fileType := header.Header.Get("Content-Type")

	params := &s3.PutObjectInput{
		Bucket:         aws.String(bucketName),
		Key:            aws.String(path),
		Body:           file,
		ContentLength:  aws.Int64(size),
		ContentType:    aws.String(fileType),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		log.Printf("bad response: %s", err)
	}
	log.Println("response" + awsutil.StringValue(resp))

    url := "https://s3-us-west-2.amazonaws.com/" + bucketName + path

    return fileName, url, nil
}
