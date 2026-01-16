package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// UploadToCOS 上传文件到腾讯云 COS
// objectKey: 例如 "images/abcde.jpg"
func UploadToCOS(file io.Reader, objectKey string) (string, error) {
	// 1. 定义 Bucket 地址
	// 建议在云托管控制台设置一个环境变量叫 COS_BUCKET_URL，值为：https://<你的bucket>.cos.<地域>.myqcloud.com
	// 或者直接在这里写死
	bucketURLStr := os.Getenv("COS_BUCKET_URL")
	if bucketURLStr == "" {
		// 如果没配置环境变量，请在这里填入你的 Bucket 地址
		bucketURLStr = "https://7072-prod-0gr2o3qpe533f1fb-1352691102.cos.ap-shanghai.myqcloud.com"
	}

	u, _ := url.Parse(bucketURLStr)
	b := &cos.BaseURL{BucketURL: u}

	// 2. 获取微信云托管注入的临时密钥
	// 文档重点：必须使用 OS_TEMP_TOKEN
	secretID := os.Getenv("OS_TEMP_SECRET_ID")
	secretKey := os.Getenv("OS_TEMP_SECRET_KEY")
	token := os.Getenv("OS_TEMP_TOKEN")

	// 3. 初始化客户端
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     secretID,
			SecretKey:    secretKey,
			SessionToken: token, // 关键：必须传 SessionToken
		},
	})

	// 4. 执行上传
	_, err := client.Object.Put(context.Background(), objectKey, file, nil)
	if err != nil {
		return "", fmt.Errorf("COS上传失败: %v", err)
	}

	// 5. 返回访问链接
	// 如果是私有读写，这里返回的链接无法直接访问，需要生成预签名URL
	// 如果是公有读私有写(推荐)，直接返回 URL 即可
	return fmt.Sprintf("%s/%s", bucketURLStr, objectKey), nil
}
