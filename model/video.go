package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/beevik/etree"
)

//Videolist 视频模型
type Videolist struct {
	ID       int64 `gorm:"primary_key;auto_increment"`
	Name     string
	Last     int64
	Type     string
	Pic      string
	Lang     string
	Area     string
	Year     string
	Actor    string
	Director string
	Dd       string `gorm:"type:longtext;"`
	Des      string `gorm:"type:text;"`
}

//VideoClass VideoClass
type VideoClass struct {
	ID    int64
	Type  string
	Count int64
}

//Getlist 获取视频分类
func Getlist() {
	url := os.Getenv("WEB_URL")
	//log.Println(url)

	res, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

	}
	// 读取资源数据 body: []byte
	body, err := ioutil.ReadAll(res.Body)

	// 关闭资源流
	res.Body.Close()
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		panic(err)
	}
	rss := doc.SelectElement("rss")

	class := rss.SelectElement("class")
	for _, ty := range class.SelectElements("ty") {
		id, _ := strconv.ParseInt(ty.SelectAttrValue("id", "unknown"), 10, 64)
		name := ty.Text()
		VideoClass := VideoClass{
			ID:    id,
			Type:  name,
			Count: 0,
		}

		err := DB.Where("type=?", name).Find(&VideoClass).Error
		if err != nil {
			saveerr := DB.Save(&VideoClass).Error
			if saveerr == nil {
				log.Println(VideoClass.Type + "  SUCC")

			} else {
				log.Println(string(VideoClass.Type), saveerr)
			}
		}

	}

}

//Getvideos 获取视频详情
func Getvideos(page int) {
	url := os.Getenv("WEB_URL") + "?ac=videolist&ct=1" + "&pg=" + strconv.Itoa(page+1)

	res, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

	}
	// 读取资源数据 body: []byte
	body, err := ioutil.ReadAll(res.Body)

	// 关闭资源流
	res.Body.Close()

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		panic(err)
	}
	rss := doc.SelectElement("rss")
	root := rss.SelectElement("list")
	for _, video := range root.SelectElements("video") {
		name := video.SelectElement("name")
		last := video.SelectElement("last")
		lasttime, _ := time.ParseInLocation("2006-01-02 15:04:05", last.Text(), time.Local)
		typew := video.SelectElement("type")
		pic := video.SelectElement("pic")
		lang := video.SelectElement("lang")
		area := video.SelectElement("area")
		year := video.SelectElement("year")
		actor := video.SelectElement("actor")
		director := video.SelectElement("director")
		dl := video.SelectElement("dl")
		dd := dl.SelectElement("dd")
		des := video.SelectElement("des")

		Videolist := Videolist{
			Name:     name.Text(),
			Last:     lasttime.Unix(),
			Type:     typew.Text(),
			Pic:      pic.Text(),
			Lang:     lang.Text(),
			Area:     area.Text(),
			Year:     year.Text(),
			Actor:    actor.Text(),
			Director: director.Text(),
			Dd:       dd.Text(),
			Des:      des.Text(),
		}
		videoname := Videolist.Name

		err := DB.Where("name=?", videoname).Find(&Videolist).Error
		if err != nil {
			saveerr := DB.Save(&Videolist).Error
			if saveerr == nil {
				log.Println(Videolist.Name + "   SUCC")

			} else {
				log.Println(string(Videolist.Name), saveerr)
			}
		}

	}

}

//Getpagecount 获取总页码
func Getpagecount() (pagecount, recordcount int) {
	weburl := os.Getenv("WEB_URL") + "?ac=videolist&ct=1"
	url := weburl
	//Getlist()
	res, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

	}
	// 读取资源数据 body: []byte
	body, err := ioutil.ReadAll(res.Body)

	// 关闭资源流
	res.Body.Close()

	doc1 := etree.NewDocument()
	if err := doc1.ReadFromBytes(body); err != nil {
		panic(err)
	}
	rss1 := doc1.SelectElement("rss")
	root1 := rss1.SelectElement("list")
	pagecount, _ = strconv.Atoi(root1.SelectAttrValue("pagecount", "unknown"))
	recordcount, _ = strconv.Atoi(root1.SelectAttrValue("recordcount", "unknown"))
	body = nil
	return pagecount, recordcount
}
