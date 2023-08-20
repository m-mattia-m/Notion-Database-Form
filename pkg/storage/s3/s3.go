package s3bucket

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
)

type Client struct {
	session *session.Session
}

type SessionConfig struct {
	S3StorageZone     *string
	S3StorageEndpoint *string
	S3StorageKey      *string
	S3StorageSecret   *string
	S3StorageBucket   *string
}

type FileConfig struct {
	FilePath string
	File     *multipart.FileHeader
}

func New() (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3-session: %s", err)
	}

	return &Client{
		session: &*sess,
	}, nil
}

func (s3client *Client) UploadFile(sessionConfig SessionConfig, fileConfig FileConfig) error {

	fileHandle, err := fileConfig.File.Open()
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	sess := s3client.session.Copy(&aws.Config{
		Region:           aws.String(*sessionConfig.S3StorageZone),
		Endpoint:         aws.String(*sessionConfig.S3StorageEndpoint),
		Credentials:      credentials.NewStaticCredentials(*sessionConfig.S3StorageKey, *sessionConfig.S3StorageSecret, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	uploader := s3manager.NewUploader(sess)
	upload, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(*sessionConfig.S3StorageBucket),
		Key:                aws.String(fileConfig.FilePath),
		ContentType:        aws.String(fileConfig.File.Header.Get("Content-Type")),
		ContentDisposition: aws.String("inline"),
		CacheControl:       aws.String("no-cache"),
		ACL:                aws.String("private"),
		Body:               fileHandle,
	})
	if err != nil {
		return err
	}

	fmt.Println(upload)
	return nil
}

//func (s3client *Client) DownloadBinaryFile(filename string) (models.S3File, error) {
//	objectResult, err := s3client.s3.GetObject(&s3.GetObjectInput{
//		Bucket: aws.String(*s3client.s3StorageBucket),
//		Key:    aws.String(filename),
//	})
//	if err != nil {
//		return models.S3File{}, err
//	}
//	defer func(Body io.ReadCloser) error {
//		err := Body.Close()
//		if err != nil {
//			return err
//		}
//		return nil
//	}(objectResult.Body)
//
//	data := make([]byte, int(*objectResult.ContentLength))
//	_, err = io.ReadFull(objectResult.Body, data)
//	if err != nil {
//		return models.S3File{}, err
//	}
//
//	return models.S3File{
//		FileName: filename,
//		FileType: *objectResult.ContentType,
//		FileSize: strconv.Itoa(int(*objectResult.ContentLength)),
//		FileMeta: objectResult.Metadata,
//		FileData: data,
//	}, nil
//}

//func (s3client *Client) GetS3SignedDownloadUrl(filename string) (models.S3FileUrl, error) {
//	objectResult, _ := s3client.s3.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: aws.String(*s3client.s3StorageBucket),
//		Key:    aws.String(filename),
//	})
//
//	url, err := objectResult.Presign(5 * time.Minute)
//	if err != nil {
//		return models.S3FileUrl{}, err
//	}
//
//	return models.S3FileUrl{
//		FileName: filename,
//		FileUrl:  url,
//	}, nil
//}

// https://s3.console.aws.amazon.com/s3/object/{your bucket}/{your file path}?region={the region of your bucket}&tab=overview
// https://sos-ch-dk-2.exo.io/s3/object/generated-notion-forms/mattiamueggler/text.txt?region=CH-DK-2&tab=overview
