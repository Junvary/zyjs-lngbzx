package utils

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func GetClassList(doc *goquery.Document) map[string]string {
	videoList := doc.Find("p.fr")
	idList := make(map[string]string)
	for _, v := range videoList.Nodes {
		startIndex := strings.Index(v.Attr[2].Val, "(")
		endIndex := strings.Index(v.Attr[2].Val, ")")
		tempSting := v.Attr[2].Val[startIndex+1 : endIndex-1]
		splitIndex := strings.Index(tempSting, ",'")
		idList[tempSting[:splitIndex]] = tempSting[splitIndex+2:]
	}
	//fmt.Println(idList)
	return idList
}
