package commons

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const YamlConfigPath string = "./application-config.yaml"

func NewConfig() *ConfigModel {
	return getConfig()
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
	Protocol string `yaml:"protocol"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
}

//获取配置文件
func getConfig() *ConfigModel {
	c := &ConfigModel{}
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
func (c *ConfigModel) GetMqPath() string {
	//amqp://guest:guest@localhost:5672/
	mqServer := c.RabbitMqServer

	path := mqServer.Protocol + "://" + mqServer.Username + ":" + mqServer.Password + "@" + mqServer.Ip + ":" + mqServer.Port + "/"

	log.Println("获取RabbitMq url : {}", path)

	return path
}
