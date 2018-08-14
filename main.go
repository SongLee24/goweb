package main

import (
	"os"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/ini.v1"
	"fmt"
)

type Req struct {
	Version  string `json:"version"`
	Status   string `json:"status"`
	Receiver string `json:"receiver"`
	Alerts []struct {
		Annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		EndsAt       string `json:"endsAt"`
		GeneratorURL string `json:"generatorURL"`
		Labels       struct {
			Alertname string `json:"alertname"`
			Instance  string `json:"instance"`
			Team      string `json:"team"`
		} `json:"labels"`
		StartsAt string `json:"startsAt"`
		Status   string `json:"status"`
	} `json:"alerts"`
	CommonAnnotations []interface{} `json:"-"`
	CommonLabels      struct {
		Alertname string `json:"alertname"`
		Team      string `json:"team"`
	} `json:"commonLabels"`
	ExternalURL string `json:"externalURL"`
	GroupKey    string `json:"groupKey"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
}

func getConfig(section string, key string) (string) {
	config, err := ini.Load("conf.ini")
	if err != nil {
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	return config.Section(section).Key(key).String()
}

func init() {

	logFilePath := getConfig("DEFAULT", "LogPath")
	fp, _ := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND,0755)
	log.SetFormatter(&log.TextFormatter{})    // 日志格式化为JSON而不是默认的ASCII
	log.SetOutput(fp)                         // 重定向输出到文件

	level := getConfig("DEFAULT", "LogLevel")
	logLevel, _ := log.ParseLevel(level)
	log.SetLevel(logLevel)
}

func send(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		input, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Read failed:", err)
		}
		defer r.Body.Close()

		req := &Req{}
		err = json.Unmarshal(input, req)
		if err != nil {
			log.Error("json format error:", err)
		}

		for _, v := range req.Alerts {
			log.Info(v.Annotations.Description)
		}
	} else {
		log.Error("ONly support Post")
	}
}

func main() {
	port := getConfig("DEFAULT", "ListenPort")
	http.HandleFunc("/sendAlarm", send) //设置访问的路由
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
