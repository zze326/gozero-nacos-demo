package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
 * @Author: zze
 * @Date: 2022/5/20 13:23
 * @Desc: 数据库配置结构体
 */

type Mysql struct {
	User         string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

func (conf *Mysql) GetDatasource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DatabaseName,
	)
}

func (conf *Mysql) InitGorm(dstModels ...interface{}) (*gorm.DB, error) {
	Orm, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conf.GetDatasource(), // DSN data source name
		DefaultStringSize:         256,                  // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                 // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                 // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                 // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		// 禁止创建物理外键
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := Orm.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	err = Orm.AutoMigrate(dstModels...)
	if err != nil {
		return nil, err
	}
	return Orm, nil
}
