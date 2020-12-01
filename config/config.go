package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type MySQLConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
	DB       string
	Debug    bool
}

type OSSConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
}

type GRPCConfig struct {
	Preprocess           string
	Matting              string
	Beauty               string
	Openpose             string
	Segment              string
	Chicona              string
	Fusion               string
	Hair                 string
	UnionGaussianPredict string
}

type RedisConfig struct {
	Address string
}

type WechatPlatformConfig struct {
	Address      string
	DefaultAppId string
	ShareAuthUrl string
}

type SMSConfig struct {
	Account string
	Password string
}

var (
	MySQL          MySQLConfig
	OSS            OSSConfig
	GRPC           GRPCConfig
	Env            string
	Redis          RedisConfig
	WechatPlatform WechatPlatformConfig
	BaiDu BaiDuAIConfig
	SMS SMSConfig
)

type BaiDuAIConfig struct {
	ApiKey         string
	SecretKey      string
	Host         string
}

func Init(configPaths ...string) (err error) {
	if err := setup(configPaths...); err != nil {
		return err
	}
	if err := initMySQL(); err != nil {
		return err
	}

	if err := initRedis(); err != nil {
		return err
	}
	//if err := initWeiXinPlatform(); err != nil {
	//	return err
	//}
	if err := initOSS(); err != nil {
		return err
	}
	//if err := initGRPC(); err != nil {
	//	return err
	//}

	if err := initBaiDu(); err != nil {
		return err
	}
	if err := initSMS(); err != nil {
		return err
	}
	return
}

func setup(paths ...string) (err error) {
	Env = os.Getenv("GO_ENV")
	if "" == Env {
		Env = "development"
	}
	godotenv.Load(".env." + Env)
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read config file (but environment config still affected), err = %+v\n", err)
		err = nil
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return
}

func initMySQL() (err error) {
	MySQL.Host = viper.GetString("mysql.host")
	MySQL.Port = viper.GetUint("mysql.port")
	MySQL.User = viper.GetString("mysql.user")
	MySQL.Password = viper.GetString("mysql.password")
	MySQL.DB = viper.GetString("mysql.db")
	if MySQL.DB == "" {
		return errors.New("mysql.db should not be empty")
	}
	MySQL.Debug = viper.GetBool("mysql.debug")
	return
}

func initOSS() (err error) {
	OSS.Endpoint = viper.GetString("oss.endpoint")
	if OSS.Endpoint == "" {
		return errors.New("oss.endpoint should not be empty")
	}
	OSS.AccessKeyId = viper.GetString("oss.access_key_id")
	if OSS.AccessKeyId == "" {
		return errors.New("oss.access_key_id should not be empty")
	}
	OSS.AccessKeySecret = viper.GetString("oss.access_key_secret")
	if OSS.AccessKeySecret == "" {
		return errors.New("oss.access_key_secret should not be empty")
	}
	OSS.Bucket = viper.GetString("oss.bucket")
	if OSS.Bucket == "" {
		return errors.New("oss.bucket should not be empty")
	}
	return
}

func initBaiDu() (err error) {
	BaiDu.Host = viper.GetString("baidu.host")
	if BaiDu.Host == "" {
		return errors.New("baidu.host should not be empty")
	}
	BaiDu.ApiKey = viper.GetString("baidu.app_id")
	if BaiDu.ApiKey == "" {
		return errors.New("baidu.app_id should not be empty")
	}
	BaiDu.SecretKey = viper.GetString("baidu.app_key")
	if BaiDu.SecretKey == "" {
		return errors.New("baidu.app_key should not be empty")
	}

	return
}

func initSMS() (err error) {
	SMS.Account = viper.GetString("sms.account")
	SMS.Password = viper.GetString("sms.password")
	return
}

func initRedis() (err error) {
	Redis.Address = viper.GetString("redis.address")
	if Redis.Address == "" {
		return errors.New("redis.address should not be empty")
	}
	return
}

//func initWeiXinPlatform() (err error) {
//	WechatPlatform.Address = viper.GetString("wechatPlatform.address")
//	WechatPlatform.DefaultAppId = viper.GetString("wechatPlatform.defaultAppId")
//	WechatPlatform.ShareAuthUrl = viper.GetString("wechatPlatform.shareAuthUrl")
//	if WechatPlatform.Address == "" {
//		return errors.New("wechatPlatform:address should not be empty")
//	}
//	if WechatPlatform.DefaultAppId == "" {
//		return errors.New("wechatPlatform:defaultAppId should not be empty")
//	}
//	if WechatPlatform.ShareAuthUrl == "" {
//		return errors.New("wechatPlatform:shareAuthUrl should not be empty")
//	}
//	return
//}

//func initGRPC() (err error) {
//	addresses := map[string]*string{
//		"grpc.preprocess":             &GRPC.Preprocess,
//		"grpc.matting":                &GRPC.Matting,
//		"grpc.beauty":                 &GRPC.Beauty,
//		"grpc.openpose":               &GRPC.Openpose,
//		"grpc.segment":                &GRPC.Segment,
//		"grpc.chicona":                &GRPC.Chicona,
//		"grpc.fusion":                 &GRPC.Fusion,
//		"grpc.hair":                   &GRPC.Hair,
//		"grpc.union_gaussian_predict": &GRPC.UnionGaussianPredict,
//	}
//	for k, v := range addresses {
//		a := viper.GetString(k)
//		if err = checkReachable(k, a); err != nil {
//			return
//		}
//		*v = a
//	}
//	return
//}
//
//func checkReachable(name, address string) error {
//	if address == "" {
//		return errors.Errorf("%s should not be empty", name)
//	}
//	conn, err := grpc.Dial(address, grpc.WithInsecure())
//	if err != nil {
//		return errors.Wrapf(err, "%s address is unreachable, address: %s", name, address)
//	}
//	defer conn.Close()
//	return nil
//}
