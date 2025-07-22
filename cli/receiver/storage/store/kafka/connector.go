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
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Connector struct {
	connection *kafka.Writer
	config     map[string]string
}

func (c *Connector) Init(cfg map[string]string) error {
	if cfg == nil {
		return fmt.Errorf("Не корректная ссылка на конфигурацию")
	}

	c.config = cfg
	topic := c.config["topic"]
	addrStr := fmt.Sprintf("%s:%s", c.config["host"], c.config["port"])

	w := &kafka.Writer{
		Addr:         kafka.TCP(addrStr),
		Topic:        topic,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
		BatchSize:    1,
	}

	c.connection = w

	return nil
}

func (c *Connector) Save(msg interface{ ToBytes() ([]byte, error) }) error {
	if msg == nil {
		return fmt.Errorf("Не корректная ссылка на пакет")
	}

	innerPkg, err := msg.ToBytes()
	if err != nil {
		return fmt.Errorf("Ошибка сериализации  пакета: %v", err)
	}

	err = c.connection.WriteMessages(context.Background(), kafka.Message{Value: []byte(innerPkg)})
	if err != nil {
		return fmt.Errorf("Ошибка отправки сырого пакета в Kafka: %v", err)
	}
	
	return nil
}

func (c *Connector) Close() error {
	if c != nil && c.connection != nil {
		if err := c.connection.Close(); err != nil {
			return err
		}
	}
	return nil
}
