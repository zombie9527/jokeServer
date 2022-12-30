package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type JuheRequest struct {
	Date string `json:"date"`
	Key  string `json:"key"`
}

type Response struct {
	Result ResultData `json:"result"`
}

type ResultData struct {
	Data []JokeList `json:"data"`
}
type JokeList struct {
	Content  string `json:"content"`
	Unixtime int64  `json:"unixtime"`
}

var jokeList []string
var unixtime int64 = 1514708824
var unixtimeReal int64
var page int = 1
var current = 0

const (
	JuheKey  = "juhekey"
	pageSize = 20
)

func getJokes() {
	url := "http://v.juhe.cn/joke/content/list.php?sort=asc&time=" + strconv.FormatInt(unixtime, 10) + "&key=" + JuheKey + "&page=" + strconv.Itoa(page) + "&pagesize=" + strconv.Itoa(pageSize)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 一次性读取
	var res Response
	bs, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(bs), &res)
	if err != nil {
		return
	}
	list := res.Result.Data
	jokeList = jokeList[:0]
	for _, item := range list {
		unixtimeReal = item.Unixtime
		jokeList = append(jokeList, item.Content)
	}
	page++

}

func getJoke(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	arg := values.Get("type")
	var body string
	if arg == "1" {
		if jokeList == nil {
			getJokes()

		}
		if current >= 20 || len(jokeList)-2 < current {
			getJokes()
			if len(jokeList) < pageSize {
				unixtime = unixtimeReal
				page = 1
			}
			current = 0
		}

		body = jokeList[current]
		current++
	}

	w.Write([]byte(body))

}

func main() {
	http.HandleFunc("/j", getJoke)
	http.ListenAndServe(":8090", nil)
}
