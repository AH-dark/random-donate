package database

import (
	"github.com/AH-dark/random-donate/pkg/conf"
	"github.com/AH-dark/random-donate/pkg/utils"
)

func needMigration() bool {
	var str Setting
	err := DB.Where(&Setting{
		Name: "db_version",
		Type: "system",
	}).First(str).Error

	if err != nil {
		return true
	}

	if str.Value != conf.DbVersion {
		return true
	}

	return false
}

func migration() {
	if !needMigration() {
		utils.Log().Info("数据库版本匹配，跳过初始化")
		return
	}

	utils.Log().Info("准备进行数据库初始化")

	// 自动迁移模式
	if conf.DatabaseConfig.Type == "mysql" {
		DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	err := DB.AutoMigrate(&Setting{})
	if err != nil {
		utils.Log().Panic("数据库初始化时错误，", err.Error())
		return
	}

	addDefaultSettings()

	utils.Log().Info("数据库初始化结束")
}

func addDefaultSettings() {
	for _, value := range defaultSettings {
		DB.Where(&Setting{
			Name: value.Name,
		}).Create(&value)
	}
}
