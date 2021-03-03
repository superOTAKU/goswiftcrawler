package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//BaseURL 爬虫的起点
const baseURL string = "https://www.swiftcodelist.com/"

//获取所有的bankurls
//使用管道控制go routine数量
func getBankUrls() *list.List {
	listCount := getBankListPageCount()
	log.Printf("bankListCount is %d, next step is get all bank urls!\n", listCount)
	listChan := make(chan *list.List, listCount)
	for i := 1; i <= listCount; i++ {
		go getBankUrlsOf(i, listChan)
	}
	bankUrls := list.New()
	for i := 0; i < listCount; i++ {
		bankURLList := <-listChan
		bankUrls.PushBackList(bankURLList)
	}
	log.Println("get All Bank Urls Finish!")
	return bankUrls
}

//index 当前列表页下标
//listChan 收集列表页，投递到这个chan
func getBankUrlsOf(index int, listChan chan *list.List) {
	log.Println("start get bankurls of page", index)
	retry := 0
	for {
		resp, err := http.Get(baseURL + "banks-" + fmt.Sprint(index) + ".html")
		if err != nil {
			retry++
			if retry <= 3 {
				log.Println(err)
				continue
			} else {
				log.Println(err)
				return
			}
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		urlList := list.New()
		doc.Find("table.magt10").First().
			Find("tr").
			FilterFunction(func(i int, s *goquery.Selection) bool {
				return i != 0
			}).
			Each(func(i int, tr *goquery.Selection) {
				url, _ := tr.Find("td").Eq(2).Find("a").Attr("href")
				urlList.PushBack(url)
			})
		listChan <- urlList
		log.Println("finish get bankurls of page", index)
		return
	}
}

func getBankListPageCount() int {
	retry := 0
	for {
		resp, err := http.Get(baseURL + "banks.html")
		if err != nil {
			retry++
			if retry <= 3 {
				log.Println(err)
				continue
			} else {
				log.Fatalln(err)
			}
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		return getListCount(doc)
	}
}
