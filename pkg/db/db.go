package db

import (
	"backend/config"
	"backend/pkg/logger"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Database interface {
	Instance() *gorm.DB
	Close() error
	Ping()
}

type MySQL struct {
	DB *gorm.DB
	l  logger.Logger
}

func NewMySQL(cfg config.MySQL, l logger.Logger) (*MySQL, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Name)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	return &MySQL{DB: db, l: l}, nil
}

func (p *MySQL) Instance() *gorm.DB {
	return p.DB
}

func (p *MySQL) Close() error {
	if p.DB == nil {
		return errors.New("db connection is already closed")
	}
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (p *MySQL) Ping() {
	sqlDB, err := p.DB.DB()
	if err != nil {
		text := fmt.Sprintf("mysql ping error: %s", err.Error())
		p.l.Error(text, err)
	}

	err = sqlDB.Ping()
	if err != nil {
		text := fmt.Sprintf("mysql ping error: %s", err.Error())
		p.l.Error(text, err)
	}
}
