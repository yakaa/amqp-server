package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"ampp-server/common/rabbitmq"
	"ampp-server/config"
	"ampp-server/handler"
	"ampp-server/model"
	"ampp-server/service"
	"github.com/streadway/amqp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/log4g"
)

var configFile = flag.String("f", "config/config.json", "Please set config file")

//
//func main() {
//	flag.Parse()
//	body, err := ioutil.ReadFile(*configFile)
//	if err != nil {
//		log.Fatalf("read file %s: %s", *configFile, err)
//	}
//	conf := new(config.Config)
//	if err := json.Unmarshal(body, conf); err != nil {
//		log.Fatalf("json.Unmarshal %s: %s", *configFile, err)
//	}
//	mpsPublisher, err := rabbitmq.NewPublisher(
//		conf.MpsRabbitMq.DataSource,
//		conf.MpsRabbitMq.QueueName,
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	i := 1
//	for {
//		operate := rabbitmq.InsertType
//		log4g.Info(mpsPublisher.Push(rabbitmq.Message{
//			DataBase: "mps",
//			Table:    "AuthItem",
//			Operate:  operate,
//			Data: map[string]interface{}{
//				"name":        "yk_" + strconv.Itoa(i),
//				"type":        "type_" + strconv.Itoa(i),
//				"description": "description" + strconv.Itoa(i),
//				"bizrule":     "bizrule_" + strconv.Itoa(i),
//				"data":        "data_" + strconv.Itoa(i),
//			},
//		}));
//		time.Sleep(1 * time.Second)
//		i++
//	}
//}
func main() {
	flag.Parse()
	body, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("read file %s: %s", *configFile, err)
	}
	conf := new(config.Config)
	if err := json.Unmarshal(body, conf); err != nil {
		log.Fatalf("json.Unmarshal %s: %s", *configFile, err)
	}
	log4g.Init(conf.Log4g)
	mpsMysqlEngine, err := xorm.NewEngine("mysql", conf.MpsMysql.DataSource)
	erpMysqlEngine, err := xorm.NewEngine("mysql", conf.ErpMysql.DataSource)
	romeoMysqlEngine, err := xorm.NewEngine("mysql", conf.RomeoMysql.DataSource)

	mysqlEngine, err := xorm.NewEngine("mysql", conf.AmqpMysql.DataSource)
	if err != nil {
		log.Fatalf("read file %s: %s", *configFile, err)
	}
	amqpDial, err := amqp.Dial(conf.RabbitMq.DataSource)
	if err != nil {
		log.Fatalf("create connect fail %+v", err)
	}

	mpsHandler := handler.NewMpsHandler(service.NewMpsService(
		model.NewBaseModel(mpsMysqlEngine),
		model.NewMessagesModel(mysqlEngine),
	))
	erpHandler := handler.NewErpHandler(service.NewErpService(
		model.NewBaseModel(erpMysqlEngine),
		model.NewMessagesModel(mysqlEngine),
	))
	romeoHandler := handler.NewRomeoHandler(service.NewRomeoService(
		model.NewBaseModel(romeoMysqlEngine),
		model.NewMessagesModel(mysqlEngine),
	))
	mpsConsumer := rabbitmq.BuildConsumer(amqpDial, conf.RabbitMq.MpsQueueName, mpsHandler.Consumer)
	erpConsumer := rabbitmq.BuildConsumer(amqpDial, conf.RabbitMq.ErpQueueName, erpHandler.Consumer)
	romeoConsumer := rabbitmq.BuildConsumer(amqpDial, conf.RabbitMq.RomeoQueueName, romeoHandler.Consumer)
	defer rabbitmq.Close(amqpDial)
	rabbitmq.RunConsumes(erpConsumer, romeoConsumer, mpsConsumer)

}
