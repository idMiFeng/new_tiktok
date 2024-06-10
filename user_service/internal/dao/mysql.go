package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"user_service/config"
	"user_service/global"
	"user_service/internal/model"
	"user_service/logger"
)

// 初始化MySQL连接
func InitMysql(cfg *config.MySQLConfig) (err error) {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	global.Db = db

	// 额外的连接配置
	sqlDB, err := db.DB() // database/sql.DB
	if err != nil {
		return
	}

	// 以下配置要配合 my.conf 进行配置
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	migration()
	return
}
func migration() {
	// 自动迁移
	err := global.Db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&model.User{}, &model.UserCount{})
	if err != nil {
		logger.Log.Error(err)
	}

}
