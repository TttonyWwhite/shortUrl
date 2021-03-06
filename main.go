package main

import (
	_ "./models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengan/shortUrl/models"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	prefix = "localhost:8081/short/" // 因为是在localhost实现，所以前缀还不够短.
)

var m *sync.Mutex

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	m = new(sync.Mutex)
	router.GET("/long", func(context *gin.Context) {
		longUrl := context.Query("longUrl")
		context.String(http.StatusOK, Shorter(longUrl))
	})

	router.GET("/short", func(context *gin.Context) {
		shortUrl := context.Query("shortUrl")
		longUrl := Longer(shortUrl)
		context.String(http.StatusOK, longUrl)
	})

	router.Run(":8081")

}

func WebRoot(context *gin.Context) {
	context.String(http.StatusOK, "hello, world")
}

// 将长url转换为短url
func Shorter(longUrl string) (shortUrl string) {
	// 先试着从redis缓存查询
	shortUrl, err := models.GetShortUrlFromRedis(longUrl)
	if err != nil {
		// redis缓存没有命中
		shortUrl, err = models.GetShortUrl(longUrl)
		if err != nil {
			// 数据库中也没有记录，说明是新的长地址，需要计算得到结果
			m.Lock()
			id := models.GetCount()
			str := encode(id)
			shortUrl = prefix + str
			// 将对应关系加入到mysql和redis
			models.InsertRecord(shortUrl, longUrl)
			models.AddLongToShort(longUrl, shortUrl)
			m.Unlock()
		} else {
			// mysql数据库中有记录，但是redis没有
			models.AddLongToShort(longUrl, shortUrl)
		}
	}

	return shortUrl
}

// 将短url还原为长url
func Longer(shortUrl string) (longUrl string) {
	// 先尝试从redis缓存查询
	longUrl, err := models.GetLongUrlFromRedis(shortUrl)
	if err != nil {
		// redis缓存没有命中，从mysql查询
		shortPart := strings.Split(shortUrl, "/")[2]
		id := decode(shortPart)
		longUrl, err = models.GetLongUrl(id)

		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}

		// 将对应关系写入redis
		models.AddShortToLong(shortUrl, longUrl)
	}

	return longUrl
}

/**
将一个62进制的数转换回十进制的id
*/
func decode(s string) int {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := len(dict)
	d := 0
	for _, ch := range s {
		for i, a := range dict {
			if a == ch {
				d = d*base + i
			}
		}
	}
	return d + 1
}

/**
将一个id转换成一个62进制的数(用string表示)并返回
*/
func encode(i int) string {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := len(dict)
	digits := []int{}
	for i > 0 {
		r := i % base
		digits = append([]int{r}, digits...) // 因为是从高位向低位转换，所以append函数要反过来写
		i = i / base
	}

	// 因为会有超过十进制表示范围的数，所以将这些两位数转换成字母表示形式
	r := []rune{}
	for _, d := range digits {
		r = append(r, dict[d])
	}
	return string(r)
}
