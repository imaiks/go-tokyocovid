package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

var pdfTbl map[string]bool = map[string]bool{}

func makePdfTbl() {
	files, err := filepath.Glob("*.pdf")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		pdfTbl[f] = true
	}
}

func retrieve(debugOn bool) {

	var c *colly.Collector

	if debugOn {
		c = colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	} else {
		c = colly.NewCollector()
	}
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(e.Text, "患者の発生") {
			if debugOn {
				fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			}
			c.Visit(e.Request.AbsoluteURL(link))
		} else if strings.Contains(e.Text, "別紙") {
			if debugOn {
				fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			}
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		if debugOn {
			fmt.Printf("Visiting %s, Depth:%d\n", r.URL.String(), r.Depth)
		}
	})
	c.OnResponse(func(r *colly.Response) {
		if debugOn {
			fmt.Printf("Response %s\n", r.Request.URL.String())
		}
		if strings.Contains(r.Request.URL.String(), "pdf") {
			s := strings.Split(r.Request.URL.String(), "/")
			fn := s[len(s)-1]
			if pdfTbl[fn] {
				fmt.Println("already retrieved: ", fn)
				os.Exit(0)
			} else {
				fmt.Println("download: ", fn)
				r.Save(fn)
			}
		}
	})
	c.Visit("https://www.bousai.metro.tokyo.lg.jp/taisaku/saigai/1007261/index.html")
}

func main() {
	debugOn := flag.Bool("d", false, "debug on")
	flag.Parse()
	makePdfTbl()
	retrieve(*debugOn)
}
