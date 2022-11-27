package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"zyjs-lngbzx/utils"
)

var cookieId string
var classType int
var page int

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type SaveViewRequest struct {
	json string
}

type StartWatch struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

type Pages struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

func Welcome() {
	fmt.Println("========================================================================")
	fmt.Println("欢迎使用辽宁-专业技术人员继续教育-观看软件！请先登录，完成选课，并拿到 Cookie ！")
	fmt.Println("Cookie的获取可通过打开F12并随意找到此网站的一个请求 --> Headers --> Request Headers --> Cookie")
	fmt.Println("本软件仅供学习交流使用！")
	fmt.Println("========================================================================")

	fmt.Println("请输入Cookie值(以 JSESSIONID= 开头)")
	fmt.Scanln(&cookieId)
	fmt.Println("请选择观看类型：1.必修课，2.选修课")
	fmt.Scanln(&classType)
	fmt.Println("请输入总页码数（没有页码输入1）：")
	fmt.Scanln(&page)
}

func main() {
	Welcome()
	mainUrl := "https://zyjs.lngbzx.gov.cn/study/yearplan/gostudy/" + strconv.Itoa(classType)
	for i := 1; i <= page; i++ {
		fmt.Println("**********************************************")
		fmt.Println("开始观看第" + strconv.Itoa(i) + " 页")
		//var pageData Pages
		pageData := url.Values{}
		pageData.Set("pageNo", strconv.Itoa(i))
		pageData.Set("pageSize", strconv.Itoa(12))
		pageDataString := pageData.Encode()
		pageDataByte := []byte(pageDataString)
		docMain := utils.BuildHttp(mainUrl, cookieId, "POST", bytes.NewReader(pageDataByte), true)
		idList := utils.GetClassList(docMain)
		for k, v := range idList {
			doc := utils.BuildHttp("https://zyjs.lngbzx.gov.cn/study/resource/info/"+k+"/"+v, cookieId, "GET", nil, false)
			GetClassDetail(doc)
		}
	}
	var end interface{}
	fmt.Println("感谢使用，按任意键退出，再见！")
	fmt.Scanf("%s", &end)
}

func GetClassDetail(doc *goquery.Document) {
	input := doc.Find("input")
	var value string
	var title string
	var length string
	var id string
	for _, v := range input.Nodes {
		//fmt.Println(v.Attr)
		for _, a := range v.Attr {
			//fmt.Println(a, reflect.TypeOf(a))
			if a.Key == "name" && a.Val == "title" {
				title = v.Attr[2].Val
			}
			if a.Key == "id" && a.Val == "tmsource" {
				value = v.Attr[2].Val
			}
			if a.Key == "name" && a.Val == "length" {
				length = v.Attr[2].Val
			}
			if a.Key == "name" && a.Val == "id" {
				id = v.Attr[2].Val
			}
		}
	}
	fmt.Println("开始观看：" + title)
	WatchVideo(value, title, length, id)
}

func WatchVideo(value, title, length, id string) {
	videoUrl := "https://zyjs.lngbzx.gov.cn/study/resource/saveview"
	videoData := url.Values{}
	videoData.Set("json", "{'cid':"+id+",'source': "+value+",'position': '','percent': '0'}")
	videoDataString := videoData.Encode()
	videoDataByte := []byte(videoDataString)
	client := http.Client{}
	req, err := http.NewRequest("POST", videoUrl, bytes.NewReader(videoDataByte))
	if err != nil {
		log.Fatal("观看视频：" + title + " 失败！")
	}
	req.Header.Add("Cookie", cookieId)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, err := client.Do(req)
	defer res.Body.Close()
	//fmt.Println(res, res.Body)
	body, err := io.ReadAll(res.Body)
	var sw StartWatch
	_ = json.Unmarshal(body, &sw)
	//fmt.Println(body, sw)
	randomId := sw.Id
	EndWatch(strconv.Itoa(randomId), id, title, length)
}

func EndWatch(randomId string, id string, title string, length string) {
	endUrl := "https://zyjs.lngbzx.gov.cn/study/resource/saveview"
	videoData := url.Values{}
	videoData.Set("json", "{'cid':"+id+",'historyId': "+randomId+",'position': "+length+",'len': "+length+",'position':"+"'822.07' }")
	videoDataString := videoData.Encode()
	videoDataByte := []byte(videoDataString)
	client := http.Client{}
	req, err := http.NewRequest("POST", endUrl, bytes.NewReader(videoDataByte))
	req.Header.Add("Cookie", cookieId)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Fatal("观看视频：" + title + " 失败！")
	} else {
		fmt.Println("观看 " + title + " 完成！")
		fmt.Println("=============================")
	}
}
