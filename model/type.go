package model

import (
	"log"
)

//TypeCount 统计片数
func TypeCount() {
	var VideoClasses []VideoClass
	DB.Find(&VideoClasses)

	count := 0

	for i := 0; i < len(VideoClasses); i++ {
		count = 0
		Update := VideoClasses[i]
		DB.Table("Videolists").Where("type = ?", Update.Type).Count(&count)
		log.Println("type = ", Update.Type, "  count= ", count)
		DB.Table("video_classes").Where("type = ?", Update.Type).Update("count", count)
	}
}
