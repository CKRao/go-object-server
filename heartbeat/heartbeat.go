package heartbeat

import (
	"go-object-server/commons"
	"go-object-server/mq"
	"go-object-server/util"
	"time"
)

const (
	//发送到的exchange
	exchangeName string = "apiServers"
	//间隔时间
	sendDuration time.Duration = 5 * time.Second
)

//数据服务心跳发送逻辑
//每5秒发送一次消息到apiServers exchange
//apiServers exchange 需在RabbitMQ服务端提前创建
func StartHeartBeat() {
	//新建RabbitMQ实体
	q := mq.New(commons.GetConfigIns().GetMqUrl())

	defer q.Close()

	for {
		q.Publish(exchangeName, util.Ip)
		time.Sleep(sendDuration)
	}
}
