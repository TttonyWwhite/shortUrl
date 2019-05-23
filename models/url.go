package models

import (
	"errors"
	"time"
)



type Url struct {
	ID 					int64
	LongUrl				string
	ShortURL			string
	CreateTime			time.Time
	LastContactTime		time.Time // 最后一次被使用的时间，用于考察热点信息
	ContactTimes 		int64 	  // 总共被使用的次数，用于考察热点信息
}



/**
	根据id查询长网址
 */
func GetLongUrl(id int) (longUrl string, err error) {
	var result Url

	if DB.First(&result, id).RecordNotFound() {
		return "", errors.New("record not found")
	} else {
		return result.LongUrl, nil
	}
}

func GetShortUrl(longUrl string) (shortUrl string, err error) {

	var result Url

	if DB.Where("long_url = ?", longUrl).First(&result).RecordNotFound() {
		return "", errors.New("record not found")
	} else {
		return result.ShortURL, nil
	}
}

/**
	获得数据库中的最后一条记录的id，返回这个id+1，即下一条记录的id
	用来计算shortUrl
 */
func GetCount() (count int) {
	var url Url
	DB.Last(&url)
	return int(url.ID)
}

func InsertRecord(shortUrl string, longUrl string) {

	newUrl := Url{
		LongUrl: longUrl,
		ShortURL: shortUrl,
		CreateTime: time.Now(),
		LastContactTime: time.Now(),
	}

	DB.Create(&newUrl)
}



