package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	cookieparser "github.com/shaunxu/pc/cookie-parser"
)

type ResponseDataValueKeyResult struct {
	Name   string
	Rate   float64
	Weight float64
}

type ResponseDataValue struct {
	Name        string
	Key_results []ResponseDataValueKeyResult
	Rate        float64
}

type ResponseDataReference struct{}

type ResponseData struct {
	Value      []ResponseDataValue
	References ResponseDataReference
}

type Response struct {
	OID  string
	Code int
	Data ResponseData
}

func main() {
	hostPtr := flag.String("h", "at.pingcode.com", "host")
	periodPtr := flag.String("p", "", "period")
	userPtr := flag.String("u", "c01a95b2898c4f339ef80befbca2f037", "user")
	outputPtr := flag.String("o", "text", "output format")
	ratePtr := flag.Bool("r", false, "with rate column")
	cookieFilePath := flag.String("c", "", "cookie file path in 'Netscape Cookie File' format")
	flag.Parse()
	url := fmt.Sprintf("https://%s/api/teams/periods/%s/users/%s/objectives", *hostPtr, *periodPtr, *userPtr)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	if *cookieFilePath != "" {
		parser := cookieparser.NCFCookieParser{Path: *cookieFilePath}
		cookies, err := parser.Parse()
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("Cookie", parser.String(cookies))
	}

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	response := Response{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalln(err)
	}

	for i, value := range response.Data.Value {
		response.Data.Value[i].Rate = 0 // do not use the objective rate from api since its not correct
		for _, kr := range value.Key_results {
			response.Data.Value[i].Rate += kr.Rate * kr.Weight
		}
	}

	builder := new(strings.Builder)
	switch *outputPtr {
	case "markdown":
		builder.WriteString("| 目标 | 关键点 |")
		if *ratePtr {
			builder.WriteString(" 完成情况 |")
		}
		builder.WriteString("\n")

		builder.WriteString("| --- | --- |")
		if *ratePtr {
			builder.WriteString(" --- |")
		}
		builder.WriteString("\n")

		for _, value := range response.Data.Value {
			builder.WriteString(fmt.Sprintf("| %s ", value.Name))
			builder.WriteString("| <ul>")
			for _, kr := range value.Key_results {
				if *ratePtr {
					builder.WriteString(fmt.Sprintf("<li>%s（%.0f%%）</li>", kr.Name, kr.Rate*100))
				} else {
					builder.WriteString(fmt.Sprintf("<li>%s</li>", kr.Name))
				}
			}
			builder.WriteString("</ul> |")
			if *ratePtr {
				builder.WriteString(fmt.Sprintf(" %.0f%% |", value.Rate*100))
			}
			builder.WriteString("\n")
		}
	case "json":
	default:
		for _, value := range response.Data.Value {
			builder.WriteString(fmt.Sprintf("目标: %s （%.0f%%）\n", value.Name, value.Rate*100))
			for _, kr := range value.Key_results {
				builder.WriteString(fmt.Sprintf("\t关键结果: %s（%.0f%%）\n", kr.Name, kr.Rate*100))
			}
		}
	}
	fmt.Print(builder.String())

}
