package main

import (
	"fmt"
	"net/http"
	"./logger"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"os"
	"runtime"
	"time"
)

const (
	Host=[]string{
		"https://www.ifa.plus",
		"https://b1.ifa.plus",
		"https://b2.ifa.plus",
		"https://b3.ifa.plus",
		"https://b4.ifa.plus",
		"https://b5.ifa.plus",
		"https://b6.ifa.plus",
		"https://b7.ifa.plus",
		"https://b8.ifa.plus",
		"https://b9.ifa.plus",
		"https://b10.ifa.plus",
		"https://b11.ifa.plus",
		"https://b12.ifa.plus"
	}

	// Host    = "https://github.com"
	PageUrl = "https://github.com/paritytech/parity/releases"
)

type MinerRsJson struct {
	Balance  float64  `json:"balance"`
	Hashrate []int64  `json:"hashrate"`
	Time     []string `json:"time"`
}

type HashRate struct {
	Time     time.Time
	HashRate int64
}

type WaletRsJons struct {
	Balance float64 `json:"balance"`
}

func HttpReq(host string, url string) (*http.Response, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Set("Accept-Languag", "e:zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)

	return client.Do(req)
}
func getBestVersion() (bestVersion string, err error) {
	defer func() {
		if rev := recover(); rev != nil {
			logger.Error("getBestVersion :", rev)
		}
	}()

	return getPage()
}

func getPage() (page string, err error) {

	var doc *goquery.Document
	var res *http.Response
	res, err = HttpReq(Host, PageUrl)
	if err != nil {
		logger.Error("getPage err:", err)
		return
	}
	for i := 0; i < 10; i++ {
		doc, err = goquery.NewDocumentFromResponse(res)

		if err != nil {
			logger.Error("NewDocumentFromResponse", err)
			continue
		}
		// div_release_header := doc.Find("div.release-header")
		div_release_header := doc.Find("div.repository-content")
		// fmt.Println("div_release_header:", div_release_header.Text())
		h1 := div_release_header.Find("h3")
		_, err := h1.Html()
		if err != nil {
			logger.Error(err)
		}
		// fmt.Println("h1:", txt)
		a := h1.Find("a")
		// fmt.Println("a:", a.Text())
		page = a.Text()
		break
		// if page, ok := a.Attr("href"); ok {
		// 	fmt.Println("got %v.", page)
		// 	break
		// }
	}

	return
}

type log_conf struct {
	Dir     string
	Name    string
	Console bool
	Num     int32
	Size    int64
	Level   string
}
type Configuration struct {
	Log       log_conf
	Exec_time string
}

var configuration Configuration

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, _ := os.Open("./conf.json")
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	logger.Info("Exec_time：", configuration.Exec_time)
	log_init()
}

func log_init() {
	logger.SetConsole(configuration.Log.Console)
	logger.SetRollingFile(configuration.Log.Dir, configuration.Log.Name, configuration.Log.Num, configuration.Log.Size, logger.KB)
	//ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF
	logger.SetLevel(logger.ERROR)
	if configuration.Log.Level == "info" {
		logger.SetLevel(logger.INFO)
	} else if configuration.Log.Level == "error" {
		logger.SetLevel(logger.ERROR)
	}
}
func start() {
	// i := 0
	// for ; i < 10; i++ {
	// 	ret, err := getBestVersion()
	// 	if err != nil {
	// 		logger.Error("获取parity版本错误:", err)
	// 	} else {
	// 		logger.Info("parity 最近版本：" + ret)
	// 		break
	// 	}
	// }
	// if i >= 10 {
	// 	logger.Error("获取parity版本错误:retry more than 10 times")
	// }
}
func crawlall(){

}
func main() {
	// testSendTransactionWithNonce()
	crawlall()
	c := cron.New()

	Interval := configuration.Exec_time
	logger.Info("config is :", Interval)
	c.AddFunc(Interval, func() {
		start()
	})

	//启动定时任务
	c.Start()
	ch := make(chan int, 1)
	<-ch
}
