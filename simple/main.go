package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type config struct {
	endpoint   string
	ak         string
	sk         string
	bucketname string
}

func init() {
	flag.StringVar(&fileName, "f", "", "specify the filename which you want to upload")
	flag.BoolVar(&help, "h", false, "print help info")
	flag.Usage = func() {
		fmt.Println(`version: 0.0.1
Usage: cloud-station -f [-h]
Option:`)
		flag.PrintDefaults()
	}
}

func NewDefaultConfig(ep, ak, sk, bn string) *config {
	return &config{
		endpoint:   ep,
		ak:         ak,
		sk:         sk,
		bucketname: bn,
	}
}

var (
	conf     = NewDefaultConfig("oss-cn-shanghai.aliyuncs.com", "LTAI5tKDTfkPwAycQNkdfZ8K", "gooSoRSB1LNmkUaiCxuYzkHUgyME3v", "cloud-station")
	help     bool
	fileName string
)

func uploadfile(filename string) (downloadUrl string, err error) {
	client, err := oss.New(conf.endpoint, conf.ak, conf.sk)
	if err != nil {
		err = fmt.Errorf("new oss client error, %s", err)
		return
	}

	bucket, err := client.Bucket(conf.bucketname)
	if err != nil {
		err = fmt.Errorf("get bucket %s error, %s", conf.bucketname, err)
		return
	}

	err = bucket.PutObjectFromFile(filename, filename)
	if err != nil {
		err = fmt.Errorf("put file %s error, %s", filename, err)
		return
	}

	downloadUrl, err = bucket.SignURL(filename, oss.HTTPGet, 60*60*24*3)
	if err != nil {
		err = fmt.Errorf("generate download url error, %s", err)
		return
	}

	return
}

func loadParams() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	loadParams()
	downloadUrl, err := uploadfile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("下载地址：", downloadUrl)
}
