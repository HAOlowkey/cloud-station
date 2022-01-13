package aliyun_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/HAOlowkey/cloud-station/pro/store/aliyun"
	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	assert := assert.New(t)

	fmt.Println(ep, ak, sk, bn)
	uploader, err := aliyun.NewAliyunOssUploader(ep, ak, sk, bn)
	if assert.NoError(err) {
		downloadUrl, err := uploader.Upload("bbb", "bbb")
		if assert.NoError(err) {
			assert.NotEmpty(downloadUrl)
			fmt.Println(downloadUrl)
		}
	}

}

var (
	ep, ak, sk, bn string
)

func init() {
	ep = os.Getenv("ALI_OSS_ENDPOINT")
	ak = os.Getenv("ALI_AK")
	sk = os.Getenv("ALI_SK")
	bn = os.Getenv("ALI_BUCKET_NAME")
}
