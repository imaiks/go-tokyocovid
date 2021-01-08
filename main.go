package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		} else if strings.Contains(e.Text, "最新の本部報") {
			if debugOn {
				fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			}
			c.Visit(e.Request.AbsoluteURL(link))
		} else if strings.Contains(e.Text, "報ー第") {
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

	c.Visit("https://www.bousai.metro.tokyo.lg.jp/taisaku/saigai/1010035/index.html")
}

func csvOut(files []string, debugOn bool) {
	out, err := exec.Command("pdfgrep", "..*", files[0]).Output()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		s := strings.Trim(scanner.Text(), " ")
		//fmt.Printf("town before:[%s]\n", s)
		//fmt.Println("town before:[", len(s), "]")
		if len(s) > 0 && (strings.HasPrefix(s, "千代田") ||
			strings.HasPrefix(s, "世田谷") ||
			strings.HasPrefix(s, "江戸川") ||
			strings.HasPrefix(s, "小平") ||
			strings.HasPrefix(s, "多摩") ||
			strings.HasPrefix(s, "新島")) {
			//fmt.Println("town in:[", s, "]")
			if scanner.Scan() {
				s := strings.Split(strings.Trim(scanner.Text(), " "), " ")
				fmt.Println(s)
				fmt.Println("len: ", len(s))
				wi := 0
				for _, w := range s {
					if len(w) != 0 {
						wi++
						fmt.Printf("%d [%v]\n", wi, w)
					}
				}
			}
		}
	}

}

func main() {
	debugOn := flag.Bool("d", false, "debug on")
	csv := flag.Bool("c", false, "csv out")
	flag.Parse()
	makePdfTbl()
	if !*csv {
		retrieve(*debugOn)
	} else if *csv {
		csvOut(flag.Args(), *debugOn)
	}
}
