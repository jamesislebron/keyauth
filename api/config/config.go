package config

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"openauth/api/logger"
	"openauth/api/logger/logrus"
)

var (
	db       *sql.DB
	oalogger logger.OpenAuthLogger
)

// Configer use to get conf
type Configer interface {
	GetConf() (*Config, error)
}

// Config is service conf
type Config struct {
	APP   *appConf   `toml:"app"`
	MySQL *mysqlConf `toml:"mysql"`
	Log   *logConf   `toml:"log"`
}

type appConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type mysqlConf struct {
	Host        string `toml:"host"`
	Port        string `toml:"port"`
	User        string `toml:"user"`
	Pass        string `toml:"pass"`
	DB          string `toml:"db"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxLifeTime int    `toml:"max_life_time"`
}

type logConf struct {
	Name     string `toml:"name"`
	Level    string `toml:"level"`
	FilePath string `toml:"path"`
}

// Validate use to check the service config
func (c *Config) Validate() error {
	if err := c.validateAPP(); err != nil {
		return err
	}

	if err := c.validateMySQL(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateAPP() error {
	if c.APP == nil {
		c.APP = &appConf{}
	}

	if c.APP.Host == "" {
		c.APP.Host = "0.0.0.0"
	}
	if c.APP.Port == "" {
		c.APP.Port = "8080"
	}

	return nil
}

func (c *Config) validateMySQL() error {
	if c.MySQL == nil {
		c.MySQL = &mysqlConf{}
	}

	if c.MySQL.Host == "" {
		c.MySQL.Host = "127.0.0.1"
	}
	if c.MySQL.Port == "" {
		c.MySQL.Port = "3306"
	}

	if c.MySQL.User == "" || c.MySQL.Pass == "" || c.MySQL.DB == "" {
		return errors.New("mysql user or pass or db isn't config")
	}

	return nil
}

// GetDBConn use to get mysql database connection
func (c *Config) GetDBConn() (*sql.DB, error) {
	var (
		err  error
		once sync.Once
	)

	once.Do(func() {
		err = c.initDBConn()
	})

	if err != nil {
		return nil, err
	}

	return db, nil

}

// GetLogger use to get logger instance
func (c *Config) GetLogger() (logger.OpenAuthLogger, error) {
	var (
		err  error
		once sync.Once
	)

	opts := logger.Opts{Name: c.Log.Name, Level: c.Log.Level, FilePath: c.Log.FilePath}
	once.Do(func() {
		oalogger, err = logrus.NewLogrusLogger(&opts)
	})

	if err != nil {
		return nil, err
	}

	return oalogger, nil

}

func (c *Config) initDBConn() error {

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", c.MySQL.User, c.MySQL.Pass, c.MySQL.Host, c.MySQL.Port, c.MySQL.DB)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(c.MySQL.MaxOpenConn)
	db.SetMaxIdleConns(c.MySQL.MaxIdleConn)
	db.SetConnMaxLifetime(time.Minute * time.Duration(c.MySQL.MaxLifeTime))
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}

	return nil
}