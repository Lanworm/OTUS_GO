package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/validation"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConf
	Server   ServerConf
	Storage  StorageConf
	Database DatabaseConf `validate:"required_if=Storage.Place database"`
}

type ServerConf struct {
	HTTP ServerHTTPConf
	GRPC ServerGRPCConf
}

type ServerHTTPConf struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Protocol string
	Timeout  time.Duration
}

func (s *ServerHTTPConf) GetFullAddress() string {
	address := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))

	return address
}

type ServerGRPCConf struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Protocol string
	Timeout  time.Duration
}

func (s *ServerGRPCConf) GetFullAddress() string {
	address := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))

	return address
}

type StorageConf struct {
	Place string `validate:"required,oneof=memory database"`
}

func (sc *StorageConf) InDatabase() bool {
	return sc.Place == "database"
}

type LoggerConf struct {
	Level string `validate:"required,oneof=DEBUG INFO WARNING ERROR"`
}

type DatabaseConf struct {
	User     string
	Password string
	Database string
	Host     string
	Port     string
}

func (dc *DatabaseConf) GetDsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.Database,
	)
}

func NewConfig(configFile string) (*Config, error) {
	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	c := &Config{}
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
