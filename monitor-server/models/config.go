package models

import (
	"sync"
	"log"
	"github.com/toolkits/file"
	"encoding/json"
)

type LogConfig struct {
	Enable  bool  `json:"enable"`
	Level   string  `json:"level"`
	File    string  `json:"file"`
	Stdout  bool  `json:"stdout"`
}

type LdapConfig struct {
	Enable  bool  `json:"enable"`
	Server  string  `json:"server"`
	Port   int   `json:"port"`
	BindDN  string  `json:"bindDN"`
	BaseDN  string  `json:"baseDN"`
	Filter  string  `json:"filter"`
	Attributes  []string  `json:"attributes"`
}

type HttpConfig struct {
	Port    string   `json:"port"`
	Swagger  bool  `json:"swagger"`
	Cross  bool  `json:"cross"`
	Alive   int64    `json:"alive"`
	Ldap   *LdapConfig  `json:"ldap"`
	Log    *LogConfig   `json:"log"`
}

type StoreConfig struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Server  string  `json:"server"`
	Port  int     `json:"port"`
	User  string  `json:"user"`
	Pwd   string  `json:"pwd"`
	DataBase  string  `json:"database"`
	MaxOpen  int  `json:"maxOpen"`
	MaxIdle  int  `json:"maxIdle"`
	Timeout  int  `json:"timeout"`
}

type DataSourceConfig struct {
	Env  string  `json:"env"`
	Servers  []*DatasourceServers  `json:"servers"`
	DivideTime  int64  `json:"divide_time"`
	WaitTime    int    `json:"wait_time"`
}

type DependenceConfig struct {
	Name  string  `json:"name"`
	Server  string  `json:"server"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Expire    int     `json:"expire"`
}

type DatasourceServers struct {
	Id  int  `json:"id"`
	Type  string  `json:"type"`
	Env  string  `json:"env"`
	Host  string  `json:"host"`
	Token  string  `json:"token"`
}

type PrometheusConfig struct {
	ConfigPath  string  `json:"config_path"`
	ConfigReload  string  `json:"config_reload"`
}

type GlobalConfig struct {
	Http  *HttpConfig  `json:"http"`
	Store  []*StoreConfig  `json:"store"`
	Datasource  DataSourceConfig  `json:"datasource"`
	LimitIp  []string  `json:"limitIp"`
	Dependence  []*DependenceConfig  `json:"dependence"`
	Prometheus  PrometheusConfig  `json:"prometheus"`
	TagBlacklist  []string  `json:"tag_blacklist"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func InitConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}