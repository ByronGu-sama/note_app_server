package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"math/rand/v2"
	"note_app_server/global"
	"note_app_server/utils"
	"time"
)

// UploadFileObject 用于将本地文件上传到OSS存储桶。
// @params bucketName - 存储空间名称。
// @params pathPrefix - 文件路径前缀
// @params objectName - Object完整路径，完整路径中不包含Bucket名称。
// @params file - 文件读取流
// @return 文件名，错误
func UploadFileObject(bucketName, pathPrefix string, file io.Reader, fileType string) (string, error) {
	client := global.OssClientPool.Get().(*oss.Client)
	defer global.OssClientPool.Put(client)
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	encodeFullFileName := utils.EncodePicsName(fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Int64())) + "." + fileType
	objectName := pathPrefix + encodeFullFileName

	// 上传文件。
	err = bucket.PutObject(objectName, file)
	if err != nil {
		return "", err
	}

	// 文件上传成功后，记录日志。
	log.Printf("File uploaded successfully to %s/%s", bucketName, objectName)
	return encodeFullFileName, nil
}

// GetOssObject 用于从OSS存储桶获取文件流。
// @params bucketName - 存储空间名称。
// @params objectName - Object完整路径，完整路径中不能包含Bucket名称。
func GetOssObject(bucketName, pathPrefix, objectName string) (io.ReadCloser, error) {
	client := global.OssClientPool.Get().(*oss.Client)
	defer global.OssClientPool.Put(client)
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	// 获取输入流
	object, err1 := bucket.GetObject(pathPrefix + objectName)
	if err1 != nil {
		return nil, err1
	} else {
		return object, nil
	}
}

// DeleteObject 用于删除OSS存储空间中的一个对象。
// @params bucketName - 存储空间名称。
// @params objectName - 要删除的对象名称。
func DeleteObject(bucketName, pathPrefix, objectName string) error {
	client := global.OssClientPool.Get().(*oss.Client)
	defer global.OssClientPool.Put(client)
	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 删除文件
	err = bucket.DeleteObject(pathPrefix + objectName)
	if err != nil {
		return err
	}

	// 文件删除成功后，记录日志。
	log.Printf("Object deleted successfully: %s/%s", bucketName, objectName)
	return nil
}

// CopyObjectToAnother 复制文件到另一个文件夹
// @params bucketName 桶名称
// @params originName 文件原始位置
// @params destName 文件目标存储位置
func CopyObjectToAnother(bucketName, originName, srcName string) error {
	client := global.OssClientPool.Get().(*oss.Client)
	defer global.OssClientPool.Put(client)
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	_, err1 := bucket.CopyObjectTo(bucketName, srcName, originName)
	if err1 != nil {
		return err1
	}
	return nil
}

// HasObject 判断文件是否存在
// @params bucketName 桶名称
// @params objectName 文件名
func HasObject(bucketName, pathPrefix, objectName string) (bool, error) {
	client := global.OssClientPool.Get().(*oss.Client)
	defer global.OssClientPool.Put(client)
	bucket, err1 := client.Bucket(bucketName)
	if err1 != nil {
		return false, err1
	}
	exist, err2 := bucket.IsObjectExist(pathPrefix + objectName)
	if err2 != nil {
		return false, err2
	}
	return exist, nil
}
