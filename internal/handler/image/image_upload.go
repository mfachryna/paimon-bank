package image

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mfachryna/paimon-bank/internal/common/response"
	"github.com/mfachryna/paimon-bank/internal/common/utils/validation"
	"github.com/mfachryna/paimon-bank/internal/entity"
	"go.uber.org/zap"
)

func (im *ImageHandler) Store(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		im.log.Info("required fields are missing or invalid", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}
	defer file.Close()

	// validate file Mime type
	if err := validation.ValidateImageFileType(fileHeader); err != nil {
		im.log.Info("failed to validate image file type", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	// validate file size
	if fileHeader.Size > (2 * 1024 * 1024) { // 2 MB
		im.log.Info("file size exceeds the limit (2MB)")
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "file size exceeds the limit (2MB)",
		}).GenerateResponse(w)
		return
	}

	imageUrl, err := im.UploadImageToS3(fileHeader.Filename, file)
	if err != nil {
		im.log.Info("failed to upload image to s3", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	data := &entity.Image{
		ImageUrl: imageUrl,
	}

	(&response.Response{
		HttpStatus: http.StatusOK,
		Message:    "File uploaded sucessfully",
		Data:       data,
	}).GenerateResponse(w)
}

func (im *ImageHandler) UploadImageToS3(fileName string, image multipart.File) (string, error) {
	bucketName := im.cfg.S3.BucketName
	s3Id := im.cfg.S3.ID
	s3SecretKey := im.cfg.S3.SecretKey
	s3Region := im.cfg.S3.Region

	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(
			s3Id,
			s3SecretKey,
			"",
		),
	})
	if err != nil {
		im.log.Info("failed to create new session aws", zap.Error(err))
		return "", err
	}

	fileName = generateRandomString(10) + time.Now().Format("20060102150405") + "-" + fileName

	svc := s3.New(ses)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   image,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		im.log.Info("failed to put object", zap.Error(err))
		return "", err
	}

	imageUrl := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, "ap-southeast-1", fileName)

	return imageUrl, nil
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
