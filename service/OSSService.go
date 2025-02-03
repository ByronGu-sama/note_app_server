package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"note_app_server/global"
)

var client = global.OssClient

// CreateBucket 创建bucket
func CreateBucket(bucketName string) error {
	checkClientIfNil()
	err := client.CreateBucket(bucketName)
	if err != nil {
		return err
	}
	// 存储空间创建成功后，记录日志。
	log.Printf("Bucket created successfully: %s", bucketName)
	return nil
}

// UploadFile 用于将本地文件上传到OSS存储桶。
// 参数：
//
//	bucketName - 存储空间名称。
//	objectName - Object完整路径，完整路径中不包含Bucket名称。
//	localFileName - 本地文件的完整路径。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func UploadFile(bucketName, objectName, localFileName string) error {
	checkClientIfNil()
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return err
	}

	// 文件上传成功后，记录日志。
	log.Printf("File uploaded successfully to %s/%s", bucketName, objectName)
	return nil
}

// GetOssObject 用于从OSS存储桶获取文件流。
//
//	bucketName - 存储空间名称。
//	objectName - Object完整路径，完整路径中不能包含Bucket名称。
func GetOssObject(bucketName, objectName string) (io.ReadCloser, error) {
	checkClientIfNil()
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	// 获取输入流
	object, err1 := bucket.GetObject(objectName)
	if err1 != nil {
		return nil, err1
	} else {
		return object, nil
	}
}

// ListObjects 用于列举OSS存储空间中的所有对象。
// 参数：
//
//	bucketName - 存储空间名称。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，打印所有对象；否则，返回错误。
func ListObjects(bucketName string) error {
	checkClientIfNil()
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 列举文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			return err
		}

		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			log.Printf("Object: %s", object.Key)
		}

		if !lsRes.IsTruncated {
			break
		}
		marker = lsRes.NextMarker
	}
	return nil
}

// DeleteObject 用于删除OSS存储空间中的一个对象。
// 参数：
//
//	bucketName - 存储空间名称。
//	objectName - 要删除的对象名称。
//	endpoint - Bucket对应的Endpoint。
//
// 如果成功，记录成功日志；否则，返回错误。
func DeleteObject(bucketName, objectName string) error {
	checkClientIfNil()
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 删除文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}

	// 文件删除成功后，记录日志。
	log.Printf("Object deleted successfully: %s/%s", bucketName, objectName)
	return nil
}

func checkClientIfNil() {
	if client == nil {
		client = global.OssClient
	}
}
