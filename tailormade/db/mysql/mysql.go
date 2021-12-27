package mysql

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

// Config 配置
type Config struct {
	Host     string
	DB       string
	User     string
	Password string
	Log      bool

	MaxIdleConns int
	MaxOpenConns int

	dsn string
}

const (
	// DSN_DEFAULT utf-8
	DSN_DEFAULT = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"

	// DSN_UTF8MB4 utf-8 mb4
	DSN_UTF8MB4 = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

func (c *Config) String() string {
	if c.dsn == "" {
		c.dsn = DSN_DEFAULT
	}
	return c.string(c.dsn)
}

func (c *Config) SetDSN(dsn string) {
	c.dsn = dsn
}

func (c *Config) string(format string) string {
	return fmt.Sprintf(format, c.User, c.Password, c.Host, c.DB)
}

// New 创建数据库连接
func New(config Config, log logr.Logger) (*gorm.DB, error) {
	var logger lg.Interface
	if config.Log {
		logger = newLogger(log, lg.Config{
			SlowThreshold: time.Second,
			LogLevel:      lg.Info,
			Colorful:      true,
		})
	}
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       config.String(),
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}),
		&gorm.Config{
			Logger: logger,
		})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 20
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	return db, nil
}