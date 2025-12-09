package service

type Config struct {
	Application Application `json:"application"`
	Cloud       Cloud       `json:"cloud"`
	Redis       Redis       `json:"redis"`
	DB          DB          `json:"db"`
}

type Application struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

type Cloud struct {
	Nacos     Nacos     `json:"nacos"`
	Zookeeper Zookeeper `json:"zookeeper"`
}

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

type Nacos struct {
	Enable    bool   `json:"enable"`
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
	Namespace string `json:"namespace"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Group     string `json:"group"`
}

type Zookeeper struct {
	Enable bool   `json:"enable"`
	Ip     string `json:"ip"`
	Port   uint64 `json:"port"`
}

type Redis struct {
	Enable bool   `json:"enable"`
	Ip     string `json:"ip"`
	Port   uint64 `json:"port"`
	DB     int    `json:"db"`
}
