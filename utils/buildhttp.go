package utils

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
)

func BuildHttp(url string, cookieId string, method string, body io.Reader) *goquery.Document {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Cookie", cookieId)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("目标相应错误: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
