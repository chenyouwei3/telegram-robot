package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

func InitMysql(dns string) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	// 设置最大空闲连接数为 20
	sqlDB.SetMaxIdleConns(20)
	// 设置最大打开连接数为 100
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大生命周期为 30 秒
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	DB = db
}
