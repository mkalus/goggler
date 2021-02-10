package sthree

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"io/ioutil"
	"log"
)

type S3Cache struct {
	debug  bool
	bucket string
	client *minio.Client
}

// return file path for a certain hash value
func (c S3Cache) getPathForHash(hash string) string {
	return hash[0:1] + "/" + hash[1:2] + "/" + hash[2:3] + "/" + hash + ".png"
}

func (c S3Cache) Save(hash string, data []byte) error {
	_, err := c.client.PutObject(context.Background(), c.bucket, c.getPathForHash(hash), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "image/png"})
	return err
}

func (c S3Cache) Get(hash string, maxAge int) ([]byte, error) {
	reader, err := c.client.GetObject(context.Background(), c.bucket, c.getPathForHash(hash), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// we do not need to clean stuff here because S3 implements life cycle itself

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		if c.debug {
			log.Printf("Error reading hash %s: %s", hash, err)
		}

		return nil, err
	}

	return data, nil
}

func (c S3Cache) RunCleanUp(maxAge int) {
	// we do not need to clean stuff here because S3 implements life cycle itself
	if c.debug {
		log.Print("Cleanup ignored since S3 can manage this itself.")
	}
}

// init S3 cache and populate with data
func InitS3Cache(url string, region string, bucket string, accessKey string, secretKey string, expirationDays int, ssl bool, createBucket bool, debug bool) (*S3Cache, error) {
	// Initialize minio client object.
	s3Client, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Region: region,
		Secure: ssl,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	found, err := s3Client.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalln(err)
	}
	if !found {
		if createBucket {
			// try to create bucket
			err = s3Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region})
			if err != nil {
				log.Fatalf("Bucket %s not found and could not be created", bucket)
			}
		} else {
			log.Fatalf("Bucket %s not found", bucket)
		}
	}

	// set bucket life cycle
	config := lifecycle.NewConfiguration()
	if expirationDays > 0 { // only set if larger than 0, otherwise existing rules will be deleted
		config.Rules = []lifecycle.Rule{
			{
				ID:     "expire-bucket",
				Status: "Enabled",
				Expiration: lifecycle.Expiration{
					Days: lifecycle.ExpirationDays(expirationDays),
				},
			},
		}
	}
	err = s3Client.SetBucketLifecycle(ctx, bucket, config)
	if err != nil {
		log.Fatalln(err)
	}

	return &S3Cache{
		debug:  debug,
		bucket: bucket,
		client: s3Client,
	}, nil
}
