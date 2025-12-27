package service

// Config 配置
type Config struct {
	Application Application `json:"application"`
	Cloud       Cloud       `json:"cloud"`
	Redis       Redis       `json:"redis"`
	DB          DB          `json:"db"`
}

// Application 应用配置
type Application struct {
	Name    string         `json:"name"`
	IP      string         `json:"ip"`
	Port    int            `json:"port"`
	Env     string         `json:"env"`
	Version string         `json:"version"`
	Debug   PprofDebugAuth `json:"debug"`
}

// Cloud 微服务配置
type Cloud struct {
	Nacos     Nacos     `json:"nacos"`
	Zookeeper Zookeeper `json:"zookeeper"`
}

// DB 数据库配置
type DB struct {
	Enable   bool              `json:"enable"`
	Type     string            `json:"type"`
	Ip       string            `json:"ip"`
	Port     int               `json:"port"`
	DBName   string            `json:"dbname"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Params   map[string]string `json:"params"`
}

// Nacos nacos注册中心配置
type Nacos struct {
	Enable    bool   `json:"enable"`
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
	Namespace string `json:"namespace"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Group     string `json:"group"`
}

// Zookeeper zookeeper注册中心配置
type Zookeeper struct {
	Enable bool   `json:"enable"`
	Ip     string `json:"ip"`
	Port   uint64 `json:"port"`
}

// Redis redis缓存配置
type Redis struct {
	Enable bool   `json:"enable"`
	Ip     string `json:"ip"`
	Port   uint64 `json:"port"`
	DB     int    `json:"db"`
}

type PprofDebugAuth struct {
	Enable   bool   `json:"enable"`
	Username string `json:"username"`
	Password string `json:"password"`
}
