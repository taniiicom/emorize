package bucket

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/samborkent/uuidv7"
)

type R2Uploader struct{}

func NewR2Uploader() *R2Uploader {
	return &R2Uploader{}
}

func (uploader *R2Uploader) UploadToBucket(ctx context.Context, filepath string) (string, error) {
	// .env -> 環境変数
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
		return "", err
	}

	var (
		BUCKET_ACCESS_KEY = os.Getenv("BUCKET_ACCESS_KEY")
		BUCKET_SECRET_KEY = os.Getenv("BUCKET_SECRET_KEY")
		BUCKET_URL        = os.Getenv("BUCKET_URL")
		BUCKET_PUBLIC_URL = os.Getenv("BUCKET_PUBLIC_URL")

		BUCKET_NAME = os.Getenv("BUCKET_NAME")
		region      = "auto" // R2では'region'は'auto'で問題ない
	)

	key := uuidv7.New().String()

	// aws-sdk のセッションを設定
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(BUCKET_ACCESS_KEY, BUCKET_SECRET_KEY, ""),
		Endpoint:    aws.String(BUCKET_URL),
	})
	if err != nil {
		fmt.Println("Session error:", err)
		return "", err
	}

	// アップロードするファイルを開く
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()

	// S3サービスクライアントを作成
	s3Client := s3.New(sess)

	// ファイルをアップロード
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}

	fmt.Println("File uploaded successfully.")

	bucketObjectUrl := fmt.Sprintf("%s/%s", BUCKET_PUBLIC_URL, key)

	return bucketObjectUrl, nil
}
