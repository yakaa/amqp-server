package config

import "github.com/yakaa/log4g"

//amqp.Dial("amqp://guest:guest@localhost:5672/")
type (
	Config struct {
		Log4g      log4g.Config
		RabbitMq   RabbitMq
		MpsMysql   Mysql
		AmqpMysql  Mysql
		ErpMysql   Mysql
		RomeoMysql Mysql
		RsaCert    RsaCert
	}

	RabbitMq struct {
		DataSource string
		QueueName  string
	}
	RsaCert struct {
		PublicKeyPath  string
		PrivateKeyPath string
	}
	Mysql struct {
		DataSource string
	}
)
