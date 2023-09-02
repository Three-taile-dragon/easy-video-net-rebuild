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
	MC    *MysqlConfig
	JC    *JwtConfig
	AC    *AesConfig
	Host  *HostConfig
	Live  *LiveConfig
	Vod   *VodConfig
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

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Db       string
}

type JwtConfig struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}

type AesConfig struct {
	AesKey string
}

type HostConfig struct {
	TencentOssHost string
	LocalHost      string
}

type LiveConfig struct {
	IP           string
	ClientIP     string
	Rtmp         int
	Flv          int
	Hls          int
	Api          int
	Agreement    string
	CliAgreement string
}

type VodConfig struct {
	Appid                 int64
	Key                   string
	LicenseUrl            string
	AudioVideoType        string
	RawAdaptiveDefinition int64
	PsignExpire           int64
}

func InitConfig() *Config {
	//初始化viper
	conf := &Config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("/Initial/config")
	conf.viper.AddConfigPath(workDir + "/config")
	//读入配置
	err := conf.viper.ReadInConfig()
	if err != nil {
		zap.L().Error("viper配置读入失败,err: " + err.Error())
		log.Fatalf("viper配置读入失败,err: %v \n ", err)
	}
	conf.InitZapLog()
	conf.ReadServerConfig()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	conf.ReadRedisConfig()
	conf.ReadMysqlConfig()
	conf.ReadJwtConfig()
	conf.ReadAesConfig()
	conf.ReadHostConfig()
	conf.ReadLiveConfig()
	conf.ReadVodConfig()
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

// ReadMysqlConfig 读取mysql配置
func (c *Config) ReadMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.username"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetInt("mysql.port"),
		Db:       c.viper.GetString("mysql.db"),
	}

	c.MC = mc
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

// ReadAesConfig 读取Aes配置
func (c *Config) ReadAesConfig() {
	ac := &AesConfig{
		AesKey: c.viper.GetString("aes.key"),
	}
	c.AC = ac
}

// ReadHostConfig 读取腾讯云oss配置
func (c *Config) ReadHostConfig() {
	hostConfig := &HostConfig{}
	hostConfig.TencentOssHost = c.viper.GetString("host.tencentOss.host")
	hostConfig.LocalHost = c.viper.GetString("host.local.host")
	c.Host = hostConfig
}

// ReadLiveConfig 读取直播配置
func (c *Config) ReadLiveConfig() {
	liveConfig := &LiveConfig{}
	liveConfig.IP = c.viper.GetString("live.ip")
	liveConfig.ClientIP = c.viper.GetString("live.clientIP")
	liveConfig.Rtmp = c.viper.GetInt("live.rtmp")
	liveConfig.Flv = c.viper.GetInt("live.flv")
	liveConfig.Hls = c.viper.GetInt("live.hls")
	liveConfig.Api = c.viper.GetInt("live.api")
	liveConfig.Agreement = c.viper.GetString("live.agreement")
	liveConfig.CliAgreement = c.viper.GetString("live.cliAgreement")
	c.Live = liveConfig
}

// ReadVodConfig 读取腾讯云vod配置
func (c *Config) ReadVodConfig() {
	vodConfig := &VodConfig{}
	vodConfig.Appid = c.viper.GetInt64("vod.appid")
	vodConfig.Key = c.viper.GetString("vod.key")
	vodConfig.LicenseUrl = c.viper.GetString("vod.licenseUrl")
	vodConfig.AudioVideoType = c.viper.GetString("vod.audioVideoType")
	vodConfig.RawAdaptiveDefinition = c.viper.GetInt64("vod.rawAdaptiveDefinition")
	vodConfig.PsignExpire = c.viper.GetInt64("vod.psignExpire")
	c.Vod = vodConfig
}
