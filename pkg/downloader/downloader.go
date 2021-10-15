package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob/s3blob"
	"gopkg.in/ini.v1"
)

// DownloadTimestep fetches a GRIB file from a given S3 bucket to a local directory
func DownloadTimestep(ctx context.Context, s3Bucket string, timestepID string, gribDir string) {
	start := time.Now()
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	s3cfgFilename := usr.HomeDir + "/.s3cfg"

	cfg, err := ini.Load(s3cfgFilename)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	s3accessKey := cfg.Section("").Key("access_key").String()
	s3secretKey := cfg.Section("").Key("secret_key").String()
	s3hostBase := cfg.Section("").Key("host_base").String()

	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String("default"),
		Endpoint:                      aws.String(s3hostBase),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials(s3accessKey, s3secretKey, ""),
		S3ForcePathStyle:              aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	options := s3blob.Options{UseLegacyList: true}

	bucket, err := s3blob.OpenBucket(ctx, sess, s3Bucket, &options)
	if err != nil {
		log.Fatalf("Can't open bucket: %v", err)
	}
	defer bucket.Close()

	r, err := bucket.NewReader(ctx, timestepID, nil)
	if err != nil {
		log.Fatalf("Can't create reader for %s: %v", timestepID, err)
	}
	defer r.Close()

	localGribFile, err := os.Create(fmt.Sprintf("%s/%s", gribDir, timestepID))
	if err != nil {
		log.Fatal(err)
	}
	defer localGribFile.Close()

	bytesWritten, err := io.Copy(localGribFile, r)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bytes Written: %d in %d\n", bytesWritten, time.Since(start))
}
