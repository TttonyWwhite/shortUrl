package main

import (
	_ "./models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengan/shortUrl/models"
	"log"
	"net/http"
	"strings"
)

const(
	prefix = "localhost:8081/short/" // 因为是在localhost实现，所以前缀还不够短.
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
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

func Shorter(longUrl string) (shortUrl string) {
	shortUrl, err :=  models.GetShortUrl(longUrl)
	if err != nil {
		// 数据库中没有记录，说明是新的长地址，需要计算得到结果
		id := models.GetCount()
		str := encode(id)
		shortUrl = prefix + str
		models.InsertRecord(shortUrl, longUrl)
	}

	return shortUrl

}

func Longer(shortUrl string) (longUrl string) {
	shortPart := strings.Split(shortUrl, "/")[2]
	id := decode(shortPart)
	longUrl, err := models.GetLongUrl(id)

	if err != nil {
		log.Fatalf("Error: %s", err)
		return
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
