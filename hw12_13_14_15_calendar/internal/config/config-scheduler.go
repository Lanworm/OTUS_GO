package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/validation"
	"gopkg.in/yaml.v3"
)

type SchedulerConfig struct {
	Interval string
	Logger   LoggerConf
	Rabbit   RabbitConf
	Database DatabaseConf `validate:"required_if=Storage.Place database"`
}

type RabbitConf struct {
	Dsn      string
	Producer ProducerSettings
	Consumer ConsumerSettings
}

type ProducerSettings struct {
	Exchange   Exchange
	Queue      Queue
	RoutingKey string
}

type Exchange struct {
	Name    string
	Type    string
	Durable bool
}

type Queue struct {
	Name    string
	Durable bool
}

type ConsumerSettings struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
}

func NewSchedulerConfig(configFile string) (*SchedulerConfig, error) {
	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	c := &SchedulerConfig{}
	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		log.Fatalf("parse congig file: %v", err)
	}

	err = validation.Validate(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
