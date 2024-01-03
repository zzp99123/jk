package article

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

func TestS3(t *testing.T) {
	// 腾讯云中对标 s3 和 OSS 的产品叫做 COS
	cosId, ok := os.LookupEnv("COS_APP_ID")
	if !ok {
		panic("没有找到环境变量 COS_APP_ID ")
	}
	cosKey, ok := os.LookupEnv("COS_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 COS_APP_SECRET")
	}
	//初始化sess 可以频繁的aws进行交互
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cosId, cosKey, ""),
		Region:      ekit.ToPtr[string]("ap-nanjing"),
		Endpoint:    ekit.ToPtr[string]("https://cos.ap-nanjing.myqcloud.com"),
		// 强制使用 /bucket/key 的形态
		S3ForcePathStyle: ekit.ToPtr[bool](true),
	})
	assert.NoError(t, err)
	client := s3.New(sess)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//往s3中放数据
	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      ekit.ToPtr[string]("webook-1314583317"),
		Key:         ekit.ToPtr[string]("2102"),
		Body:        bytes.NewReader([]byte("测试内容1111")),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	assert.NoError(t, err)
	//取s3中的数据
	res, err := client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: ekit.ToPtr[string]("webook-1314583317"),
		Key:    ekit.ToPtr[string]("测试文件"),
	})
	assert.NoError(t, err)
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	t.Log(string(data))
}
