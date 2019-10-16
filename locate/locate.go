package locate

import (
	"encoding/json"
	"fmt"
	"go-object-server/commons"
	"go-object-server/mq"
	"go-object-server/objects"
	"go-object-server/util"
	"log"
)

const (
	//发送到的exchange
	exchangeName string = "dataServers"
)

//判断资源是否存在
func Locate(name string) bool {
	isExists, err := util.PathExists(name)

	if err != nil {
		log.Fatalln("Locate error : ", err)
	}

	return isExists
}

type locateMsg struct {
	name string `json:"name"`
}

func StartLocate() {
	//新建RabbitMQ实体
	q := mq.New(commons.GetConfigIns().GetMqUrl())
	defer q.Close()

	q.Name = "locate"
	q.Bind(exchangeName)
	c := q.Consume()

	for msg := range c {
		mqMsg := &locateMsg{}
		err := json.Unmarshal(msg.Body, mqMsg)

		if err != nil {
			log.Fatalln("StartLocate json.Unmarshal error : ", err)
		}
		path := fmt.Sprintf("%s/objects/%s", objects.StorageRoot, mqMsg.name)
		if Locate(path) {
			q.Send(msg.ReplyTo, util.Ip)
		}
	}
}
