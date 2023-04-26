package initialize

import (
	"chatgpt-web-go/src/global"
	model "chatgpt-web-go/src/model/api/user"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitGormMysql() {
	var (
		err                          error
		dbName, user, password, host string
	)
	dbName = global.Cfg.Database.Name
	sec := global.Cfg.Database
	user = sec.User
	password = sec.Password
	host = sec.Host
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)
	global.Gdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
	}
	global.Gdb.AutoMigrate(&model.User{})
	dbConfig, _ := global.Gdb.DB()
	dbConfig.SetMaxIdleConns(10)
	dbConfig.SetMaxOpenConns(100)
	registerCallbacks()
}
func registerCallbacks() {
	global.Gdb.Callback().Create().Before("gorm:before_create").Register("update_timestamp_before_create", updateTimeStampForCreateCallback)
	global.Gdb.Callback().Update().Before("gorm:before_update").Register("update_timestamp_before_create", updateTimeStampForUpdateCallback)
}
func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil {
		if idValue, ok := db.Statement.ReflectValue.FieldByName("ID").Interface().(uint64); ok && idValue == 0 {
			node, _ := snowflake.NewNode(1)
			db.Statement.SetColumn("ID", node.Generate().Int64())
		}
		db.Statement.SetColumn("UpdateTime", time.Now().Format("2006-01-02 15:04:05"))
		db.Statement.SetColumn("CreateTime", time.Now().Format("2006-01-02 15:04:05"))

	}
}
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil {
		db.Statement.SetColumn("UpdateTime", time.Now().Format("2006-01-02 15:04:05"))
	}
}
