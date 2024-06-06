package oss

import (
	"context"
	"file/config"
	"file/pkg/consts"
	fileLogger "file/pkg/utils/logger"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var _minioClient *minio.Client

func InitMinioClient() {
	// 基本配置信息
	ossConfig := config.Config.Minio
	endpoint := fmt.Sprintf("%s:%s", ossConfig.Host, ossConfig.Port)
	fileLogger.LogrusObj.Info("Load minio configuration successfully.")

	// 创建minio客户端对象
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ossConfig.AccessKeyId, ossConfig.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		fileLogger.LogrusObj.Error("Error occurred while creating the minio aiproxy:", err)
		panic(err)
	}

	// 检查头像bucket是否存在
	ctx := context.Background()
	isExist, err := client.BucketExists(ctx, consts.AvatarBucketName)
	if err != nil {
		fileLogger.LogrusObj.Error("Error occurred checking for the existence of the user's avatar bucket:", err)
		panic(err)
	}
	if !isExist {
		// 创建头像bucket
		fileLogger.LogrusObj.Warning("User avatar bucket does not exist yet, please create it in oss first!")
		panic("User avatar bucket does not exist yet, please create it in oss first!")
	}

	// 检查文件bucket是否存在
	isExist, err = client.BucketExists(ctx, consts.FileBucketName)
	if err != nil {
		fileLogger.LogrusObj.Error("Error occurred while checking if the file bucket exists:", err)
		panic(err)
	}
	if !isExist {
		// 创建头像bucket
		fileLogger.LogrusObj.Warning("The file bucket does not exist yet, please create it in oss first!")
		panic("The file bucket does not exist yet, please create it in oss first!")
	}

	// oss基本可用
	_minioClient = client
	fileLogger.LogrusObj.Info("Loading the minio aiproxy was successful.")
}

func GetMinioClient() *minio.Client {
	return _minioClient
}
