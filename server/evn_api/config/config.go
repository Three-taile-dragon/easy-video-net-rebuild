package config

import (
	"dragonsss.cn/evn_common/logs"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

var C = InitConfig()

type Config struct {
	viper *viper.Viper
	SC    *ServerConfig
	GC    *GrpcConfig
	EC    *EtcdConfig
	JC    *JwtConfig
	Host  *HostConfig
	UP    *UploadConfig
}
type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name    string
	Addr    string
	Version string
	Weight  int64
}

type EtcdConfig struct {
	Addrs []string
}

type JwtConfig struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}

type HostConfig struct {
	TencentOssHost string
	LocalHost      string
}

type TencentConfig struct {
	SecretId        string `ini:"secretId"`
	SecretKey       string `ini:"secretKey"`
	Appid           string `ini:"appid"`
	Bucket          string `ini:"bucket"`
	Region          string `ini:"region"`
	DurationSeconds int    `ini:"durationSeconds"`
	Host            string `ini:"host"`
	TmpFileUrl      string
}

type LocalConfig struct {
	FileUrl    string
	TmpFileUrl string
}

type UploadConfig struct {
	*TencentConfig
	*LocalConfig
}

func InitConfig() *Config {
	//初始化viper
	conf := &Config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("/opt/lbk_background/evn_api/config")
	conf.viper.AddConfigPath(workDir + "/config")
	//读入配置
	err := conf.viper.ReadInConfig()
	if err != nil {
		zap.L().Error("viper配置读入失败,err: " + err.Error())
		log.Fatalf("viper配置读入失败,err: %v \n ", err)
	}
	conf.InitZapLog()
	conf.ReadUploadConfig()
	conf.ReadServerConfig()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	conf.ReadRedisConfig()
	conf.ReadJwtConfig()
	conf.ReadHostConfig()
	return conf
}

// ReadServerConfig 读取服务器地址配置
func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

// ReadGrpcConfig 读取grpc配置
func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.version")
	c.GC = gc
}

// ReadEtcdConfig 读入etcd配置
func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		zap.L().Error("etcd配置读取失败,err: " + err.Error())
		log.Fatalf("etcd配置读取失败,err: %v \n", err)
	}
	ec.Addrs = addrs
	c.EC = ec
}

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

// InitZapLog 初始化zap日志
func (c *Config) InitZapLog() {
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		zap.L().Error("zap日志服务初始化失败,err: " + err.Error())
		log.Fatalln(err)
	}
}

// ReadJwtConfig 读取 jwt配置
func (c *Config) ReadJwtConfig() {
	jc := &JwtConfig{
		AccessExp:     time.Duration(c.viper.GetInt64("jwt.accessExp")) * time.Minute,
		RefreshExp:    time.Duration(c.viper.GetInt64("jwt.refreshExp")) * time.Minute,
		AccessSecret:  c.viper.GetString("jwt.accessSecret"),
		RefreshSecret: c.viper.GetString("jwt.refreshSecret"),
	}
	c.JC = jc
}

// ReadHostConfig 读取腾讯云oss配置
func (c *Config) ReadHostConfig() {
	hostConfig := &HostConfig{}
	hostConfig.TencentOssHost = c.viper.GetString("host.tencentOss.host")
	hostConfig.LocalHost = c.viper.GetString("host.local.host")
	c.Host = hostConfig
}

func (c *Config) ReadUploadConfig() {
	tencentConfig := &TencentConfig{}
	tencentConfig.Region = c.viper.GetString("upload.tencentOss.region")
	tencentConfig.Bucket = c.viper.GetString("upload.tencentOss.bucket")
	tencentConfig.SecretId = c.viper.GetString("upload.tencentOss.secretId")
	tencentConfig.SecretKey = c.viper.GetString("upload.tencentOss.secretKey")
	tencentConfig.Appid = c.viper.GetString("upload.tencentOss.appid")
	tencentConfig.Host = c.viper.GetString("upload.tencentOss.host")
	tencentConfig.DurationSeconds = c.viper.GetInt("upload.tencentOss.durationSeconds")
	tencentConfig.TmpFileUrl = c.viper.GetString("upload.tencentOss.tmpFileUrl")

	localConfig := &LocalConfig{}
	localConfig.FileUrl = c.viper.GetString("upload.local.fileUrl")
	localConfig.TmpFileUrl = c.viper.GetString("upload.local.tmpFileUrl")

	upConfig := &UploadConfig{
		tencentConfig,
		localConfig,
	}
	c.UP = upConfig
}
