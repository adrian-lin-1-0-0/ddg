package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	engine = "https://api.duckduckgo.com/?"
	query  = "q=%s&"
	format = "format=json&"
)

var (
	q string
	n int
)

var usageStr = `
Usage: duckduckgo instant answer [options]
Common Options:
    -h, --help                       Show this message
    -q,                              query
`

func usage() {
	fmt.Printf("%s\n", usageStr)
}

func init() {
	flag.StringVar(&q, "q", "", "query")
	flag.Usage = usage
	flag.Parse()
	if q == "" {
		usage()
		os.Exit(0)
	}
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

type Topic struct {
	Text   string `json:"Text"`
	Result string `json:"Result"`
}

type RelatedTopic struct {
	Text   string  `json:"Text"`
	Result string  `json:"Result"`
	Topics []Topic `json:"Topics"`
}

type Response struct {
	RelatedTopics []RelatedTopic `json:"RelatedTopics"`
	Abstract      string         `json:"Abstract"`
}

func printLn(idx int, restult, text string) {
	fmt.Printf("\033[33m%d.\033[0m \033[32m%s\033[0m \n%s\n\n", idx, findTitle(restult), text)
}

func printWiki(abstract string) {
	fmt.Printf("\033[33mWiki\033[0m\n%s\n\n", abstract)
}

func findTitle(str string) string {
	if start := strings.Index(str, ">"); start >= 0 {
		if end := strings.Index(str[start:], "<"); end >= 0 {
			return str[start+1 : start+end]
		}
	}
	return ""
}

func main() {
	apiUrl := engine + fmt.Sprintf(query, q) + format
	res, err := http.Get(apiUrl)
	checkErr(err)
	defer res.Body.Close()
	rsByte, err := ioutil.ReadAll(res.Body)
	checkErr(err)
	response := &Response{}
	err = json.Unmarshal(rsByte, &response)
	checkErr(err)
	idx := 1
	if len(response.RelatedTopics) == 0 {
		fmt.Println("no result")
		return
	}
	printWiki(response.Abstract)
	for _, rtc := range response.RelatedTopics {
		if len(rtc.Topics) == 0 {
			if rtc.Text != "" {
				printLn(idx, rtc.Result, rtc.Text)
				idx++
			}
		} else {
			for _, topic := range rtc.Topics {
				if topic.Text != "" {
					printLn(idx, topic.Result, topic.Text)
					idx++
				}
			}
		}
	}
}
