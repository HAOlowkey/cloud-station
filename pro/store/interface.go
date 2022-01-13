package store

type OssUploader interface {
	Upload(fileName, objectKey string) (downloadUrl string, err error)
}
