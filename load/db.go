package load

import (
	"gotv/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=Root@1998 dbname=gotv port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		log.Fatalln("数据库连接失败: " + err.Error())
	}
	// 根据*grom.DB对象获得*sql.DB的通用数据库接口
	sqlDb, _ := db.DB()
	// 设置最大连接数
	sqlDb.SetMaxIdleConns(10)
	// 设置最大的空闲连接数
	sqlDb.SetMaxOpenConns(10)

	db.Migrator().AutoMigrate(&model.Comment{}, &model.LikeRecordModel{}, &model.DislikeRecord{})
	return db
}
