package commons

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

const YamlConfigPath string = "./application-config.yaml"

var config *ConfigModel

var lock sync.Mutex

//单例模式
func GetConfigIns() *ConfigModel {
	lock.Lock()
	defer lock.Unlock()

	if config == nil {
		config = getConfig()
	}

	return config
}

//配置文件结构体
type ConfigModel struct {
	Server         *server         `yaml:"server"`
	RabbitMqServer *rabbitMqServer `yaml:"rabbit"`
}

//服务器配置
type server struct {
	Port string `yaml:"port"`
}

//RabbitMq配置
type rabbitMqServer struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
}

//获取配置文件
func getConfig() *ConfigModel {
	c := &ConfigModel{}
	//读取配置文件
	//todo:当前从程序运行的目录下读取，之后需优化为可自定义配置文件路径，通过启动参数设置
	yamlFile, err := ioutil.ReadFile(YamlConfigPath)

	if err != nil {
		log.Fatalf("GetConfig Error :  %s ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		log.Fatalf("Unmarshal Error :  %s ", err)
	}

	return c
}

//获取RabbitMq url
func (c *ConfigModel) GetMqUrl() string {
	// 连接RabbitMq url -> amqp://guest:guest@localhost:5672/

	//获取RabbitMq配置信息
	mqConfig := c.RabbitMqServer

	//拼接url
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", mqConfig.Username, mqConfig.Password, mqConfig.Ip, mqConfig.Port)

	log.Println("获取RabbitMq url : {}", url)

	return url
}
