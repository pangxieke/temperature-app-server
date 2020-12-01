package models

import (
	"image"
	"temperature/config"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

var OSSBucket *oss.Bucket
var ossClient *oss.Client

func initOSS(cfg config.OSSConfig) (err error) {
	ossClient, err = oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return
	}

	OSSBucket, err = ossClient.Bucket(cfg.Bucket)
	return
}

type PostParams struct {
	Key            string `json:"key"`
	Policy         string `json:"policy"`
	OssAccessKeyId string `json:"OSSAccessKeyId"`
	Signature      string `json:"Signature"`
}

type UploadParams struct {
	Provider   string `json:"provider"`
	Post       string `json:"post"`
	PostParams `json:"postParams"`
	Put        string `json:"put"`
}

func makePolicy(key string, expiry int64) string {
	template := `{
    "expiration": "%s",
    "conditions": [
      { "key": "%s" },
      { "bucket": "%s" },
      ["content-length-range", 1, 8000000]
    ]
  }`
	expiredAt := time.Unix(time.Now().Unix()+expiry, 0).Format(time.RFC3339Nano)
	policy := fmt.Sprintf(template, expiredAt, key, OSSBucket.BucketName)
	return base64.StdEncoding.EncodeToString([]byte(policy))
}

func sign(s string) string {
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(OSSBucket.Client.Config.AccessKeySecret))
	io.WriteString(h, s)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func NewPostParams(key string, expiry int64) PostParams {
	policy := makePolicy(key, expiry)
	return PostParams{Key: key, Policy: policy, OssAccessKeyId: OSSBucket.Client.Config.AccessKeyID, Signature: sign(policy)}
}

func NewUploadParams(key string, expiry int64) (*UploadParams, error) {
	return &UploadParams{
		Provider:   "oss",
		Put:        PutUrlFor(key, expiry),
		Post:       BucketURL(),
		PostParams: NewPostParams(key, expiry),
	}, nil
}

func BucketURL(name ...string) string {
	bucket := OSSBucket.BucketName
	if len(name) > 0 {
		bucket = name[0]
	}
	parts := strings.Split(OSSBucket.Client.Config.Endpoint, "//")
	proto, host := parts[0], parts[1]
	return fmt.Sprintf("%s//%s.%s", proto, bucket, host)
}

func UrlFor(key string) string {
	return fmt.Sprintf("%s/%s", BucketURL(), key)
}

func GetUrlFor(key string, expiry int64, process ...string) (res string) {
	var options []oss.Option
	if len(process) > 0 && process[0] != "" {
		options = append(options, oss.Process(process[0]))
	}
	res, _ = OSSBucket.SignURL(key, http.MethodGet, expiry, options...)
	return
}

func PutUrlFor(key string, expiry int64, contentType ...string) (res string) {
	var options []oss.Option
	if len(contentType) > 0 && contentType[0] != "" {
		options = append(options, oss.ContentType(contentType[0]))
	}
	res, _ = OSSBucket.SignURL(key, http.MethodPut, expiry, options...)
	return
}

func Resign(u string, expiry int64) (res string, err error) {
	aUrl, err := url.Parse(u)
	if err != nil {
		return
	}
	name := strings.Split(aUrl.Hostname(), ".")[0]
	key := aUrl.Path[1:]
	if name == "" || key == "" {
		err = errors.New("invalid url with bucket_name: %s, key: %s, url: u")
		return
	}

	bucket, err := ossClient.Bucket(name)
	if err != nil {
		return
	}
	res, _ = bucket.SignURL(key, http.MethodGet, expiry)
	return
}

func Down(key string) (img image.Image, err error) {
	body, err := OSSBucket.GetObject(key)
	if err != nil {
		err = errors.New("invalid url")
		return
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()
	img, _, err = image.Decode(body)
	return
}
