package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//获取每个银行的swift code，保存到数据库
func getBankSwiftCodes(bankURL string) {
	listCount := getSwiftCodeListCount(bankURL)
	for i := 1; i <= listCount; i++ {
		go getBankSwiftCodesOf(bankURL, i)
	}
}

func getBankSwiftCodesOf(bankURL string, index int) {
	lastDotIdx := strings.LastIndex(bankURL, ".")
	curPage := bankURL[:lastDotIdx] + "-" + fmt.Sprintf("%d", index) + ".html"
	retry := 0
	for {
		resp, err := http.Get(curPage)
		if err != nil {
			retry++
			if retry <= 3 {
				//log.Println(err)
				continue
			} else {
				//log.Println(err)
				return
			}
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			//log.Println(err)
			return
		}
		if doc == nil {
			return
		}
		doc.Find("table.magt10").First().
			Find("tr").
			FilterFunction(func(i int, s *goquery.Selection) bool {
				return i != 0
			}).
			Each(func(i int, tr *goquery.Selection) {
				swiftCodeURL, exist := tr.Find("td").Eq(4).Find("a").Attr("href")
				if !exist {
					return
				}
				go getSwiftCode(swiftCodeURL)
			})
	}
}

func getSwiftCode(swiftCodeURL string) {
	retry := 0
	for {
		resp, err := http.Get(swiftCodeURL)
		if err != nil {
			retry++
			if retry <= 3 {
				//log.Println(err)
				continue
			} else {
				//log.Println(err)
				return
			}
		}
		defer resp.Body.Close()
		msg := ""
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil || doc == nil {
			return
		}
		doc.Find("table.magt10").First().
			Find("tr").
			Each(func(i int, tr *goquery.Selection) {
				v := strings.Trim(tr.Find("td").Eq(1).Text(), " ")
				if i == 0 {
					msg += "[SwiftCode:" + v + ", "
				} else if i == 1 {
					msg += "Country:" + v + ", "
				} else if i == 2 {
					msg += "Bank:" + v + ", "
				} else if i == 3 {
					msg += "Branch:" + v + ", "
				} else if i == 4 {
					msg += "City:" + v + ", "
				} else if i == 5 {
					msg += "ZipCode:" + v + ","
				} else if i == 6 {
					msg += "Address:" + v + "]"
					log.Println(msg)
				}
			})
	}
}

//获取每个银行的列表页数
func getSwiftCodeListCount(bankURL string) int {
	retry := 0
	for {
		resp, err := http.Get(bankURL)
		if err != nil {
			retry++
			if retry <= 3 {
				//log.Println(err)
				continue
			} else {
				//log.Println(err)
				return 0
			}
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			//log.Println(err)
			return 0
		}
		return getListCount(doc)
	}
}
