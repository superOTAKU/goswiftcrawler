package main

import (
	"fmt"
	"time"
)

//用go语言重写爬虫
func main() {
	fmt.Println("start crawler!")
	//1.获取所有的bankUrl
	urls := getBankUrls()
	//2. 获取每个bank的所有swift code
	for url := urls.Front(); url != nil; url = url.Next() {
		urlStr, _ := url.Value.(string)
		go getBankSwiftCodes(urlStr)
	}
	time.Sleep(time.Duration(1) * time.Hour)
}
