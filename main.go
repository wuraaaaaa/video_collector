package main

import (
	"video_collector/model"

	"log"
	"os"
	"time"

	"github.com/bricdu/pool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	model.Database(os.Getenv("MYSQL_DSN"))
	model.Getlist()
	model.TypeCount()
	pagecount, recordcount := model.Getpagecount()
	beforcount := 0
	aftercount := 0
	model.DB.Table("Videolists").Count(&beforcount)
	log.Println("pagecount = ", pagecount, "---recordcount: ", recordcount)
	t1 := time.Now()

	pool := pool.New(8, pagecount)
	for i := 0; i <= pagecount; i++ {
		pool.Run(func() {

			log.Println("采集到第:", i, "页,共", pagecount, "页")

			model.Getvideos(i)
		})

	}
	pool.Wait()

	isOK := time.Since(t1)
	model.DB.Table("Videolists").Count(&aftercount)
	aftcount := aftercount - beforcount
	errcount := recordcount - aftercount

	log.Println("采集完成，用时： ", isOK, "---更新条数： ", aftcount, "---失败条数：", errcount)

}
