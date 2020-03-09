package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"video_collector/model"

	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	model.Database(os.Getenv("MYSQL_DSN"))
	model.Getlist()
	model.TypeCount()
	pagecount, recordcount := model.Getpagecount()
	//pagecount = 10
	beforcount := 0
	aftercount := 0
	model.DB.Table("Videolists").Count(&beforcount)
	log.Println("pagecount = ", pagecount, "---recordcount: ", recordcount)
	t1 := time.Now()
	//PS, _ := strconv.Atoi(os.Getenv("POOL_SIZE"))
	//pool := pool.New(2, pagecount)
	//pool

	//
	for i := 0; i < pagecount; i++ {

		//pool.Run(func() {

		// p.Submit(func() {
		log.Println("采集到第:", i+1, "页,共", pagecount, "页")
		model.Getvideos(i)
		// 	wg.Done()
		// })

		//})

	}
	//	pool.Wait()

	isOK := time.Since(t1).String()
	model.DB.Table("Videolists").Count(&aftercount)
	aftcount := aftercount - beforcount
	errcount := recordcount - aftercount

	log.Println("采集完成，用时： ", isOK, " 更新条数： ", aftcount, " 失败条数：", errcount)
	msgURL := os.Getenv("MSG_URL") + "?text=用时" + isOK + "更新" + strconv.Itoa(aftcount) + "失败" + strconv.Itoa(errcount)
	req, _ := http.NewRequest("POST", msgURL, nil)
	res, _ := http.DefaultClient.Do(req)
	res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
