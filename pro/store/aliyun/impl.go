package aliyun

import (
	"fmt"
	"time"

	"github.com/HAOlowkey/cloud-station/pro/store"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type Impl struct {
	Endpoint   string `validate:"required`
	Ak         string `validate:"required"`
	Sk         string `validate:"required"`
	BucketName string `validate:"required"`
	listener   oss.ProgressListener
}

func (p *ProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	// fmt.Println(event.EventType, event.TotalBytes, event.RwBytes)

	// switch event.EventType {
	// case oss.TransferStartedEvent:
	// 	p.bar = progressbar.DefaultBytes(
	// 		event.TotalBytes,
	// 		"文件上传中",
	// 	)
	// case oss.TransferDataEvent:
	// 	p.bar.Add64(event.RwBytes)
	// case oss.TransferCompletedEvent:
	// 	fmt.Printf("\n上传完成\n")
	// case oss.TransferFailedEvent:
	// 	fmt.Printf("\n上传失败\n")
	// default:
	// }

	switch event.EventType {
	case oss.TransferStartedEvent:
		p.bar = progressbar.NewOptions64(event.TotalBytes,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			// progressbar.OptionSetWriter(os.Stdout),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(30),
			progressbar.OptionSetDescription("开始上传:"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
		p.startAt = time.Now()
		fmt.Printf("文件大小: %d\n", event.TotalBytes)
	case oss.TransferDataEvent:
		p.bar.Add64(event.RwBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\n上传完成: 耗时%d秒\n", int(time.Since(p.startAt).Seconds()))
	case oss.TransferFailedEvent:
		fmt.Printf("\n上传失败: \n")
	default:
	}
}

type ProgressListener struct {
	bar     *progressbar.ProgressBar
	startAt time.Time
}

func NewProgressListener() *ProgressListener {
	return &ProgressListener{}
}

func NewAliyunOssUploader(ep, ak, sk, bn string) (store.OssUploader, error) {
	uploader := &Impl{
		Endpoint:   ep,
		Ak:         ak,
		Sk:         sk,
		BucketName: bn,
		listener:   NewProgressListener(),
	}
	validate := validator.New()
	err := validate.Struct(uploader)
	if err != nil {
		return nil, fmt.Errorf("validator err, %s", err)
	}
	return uploader, nil
}

func (i *Impl) Upload(fileName, objectKey string) (downloadUrl string, err error) {
	client, err := oss.New(i.Endpoint, i.Ak, i.Sk)
	if err != nil {
		err = fmt.Errorf("new oss client error, %s", err)
		return
	}

	bucket, err := client.Bucket(i.BucketName)
	if err != nil {
		err = fmt.Errorf("get bucket %s error, %s", i.BucketName, err)
		return
	}

	err = bucket.PutObjectFromFile(objectKey, fileName, oss.Progress(i.listener))
	if err != nil {
		err = fmt.Errorf("put file %s error, %s", fileName, err)
		return
	}

	downloadUrl, err = bucket.SignURL(objectKey, oss.HTTPGet, 60*60*24*3)
	if err != nil {
		err = fmt.Errorf("generate download url error, %s", err)
		return
	}

	return
}
