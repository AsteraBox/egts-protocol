package kafkawrt

/*
Плагин для работы с Kafka.

Раздел настроек, которые должны отвечать в конфиге для подключения хранилища:

host = "localhost"
port = "5672"
topic = "egts-messages"
partition = 0
*/

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"strconv"
	"time"
)

type Connector struct {
	connection *kafka.Conn
	config     map[string]string
}

func (c *Connector) Init(cfg map[string]string) error {
	var (
		err error
	)
	if cfg == nil {
		return fmt.Errorf("Не корректная ссылка на конфигурацию")
	}

	c.config = cfg
	topic := c.config["topic"]
	st := c.config["partition"]
	partition, err := strconv.Atoi(st)
	addrStr := fmt.Sprintf("%s:%s", c.config["host"], c.config["port"])

	if c.connection, err = kafka.DialLeader(context.Background(), "tcp", addrStr, topic, partition); err != nil {
		return fmt.Errorf("Ошибка установки соединеия Kafka: %v", err)
	}

	return err
}

func (c *Connector) Save(msg interface{ ToBytes() ([]byte, error) }) error {
	if msg == nil {
		return fmt.Errorf("Не корректная ссылка на пакет")
	}

	innerPkg, err := msg.ToBytes()
	if err != nil {
		return fmt.Errorf("Ошибка сериализации  пакета: %v", err)
	}

	c.connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = c.connection.WriteMessages(
		kafka.Message{Value: []byte(innerPkg)},
	)
	if err != nil {
		fmt.Errorf("Ошибка отправки сырого пакета в Kafka: %v", err)
	}
	return nil
}

func (c *Connector) Close() error {
	var err error
	if c != nil {
		if c.connection != nil {
			if err = c.connection.Close(); err != nil {
				return err
			}
		}
	}
	return err
}
