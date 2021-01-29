package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func infectCol(col string) (nCol []string) {
	s := strings.Fields(col)
	// fmt.Println("town in:", col, " len:", len(s))
	for _, w := range s {
		if _, err := strconv.Atoi(w); err == nil {
			nCol = append(nCol, w)
		}
	}
	return nCol
}

func scatterTable(file string, l int, w *csv.Writer, debugOn bool) {
	out, err := exec.Command("pdfgrep", "..*", file).Output()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	date := fmt.Sprintf("%s-%s-%s", file[:4], file[4:6], file[6:8])
	var nInfected = []string{date}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		s := strings.Trim(scanner.Text(), " ")
		if len(s) > 0 && (strings.HasPrefix(s, "千代田") ||
			strings.HasPrefix(s, "世田谷") ||
			strings.HasPrefix(s, "江戸川") ||
			strings.HasPrefix(s, "小平") ||
			strings.HasPrefix(s, "多摩") ||
			strings.HasPrefix(s, "新島")) {
			region := strings.Fields(s)
			if debugOn {
				fmt.Println("town in:", region)
			}
			if scanner.Scan() {
				col := scanner.Text()
				if len(col) == 0 {
					scanner.Scan()
					col = scanner.Text()
				}
				nInfected = append(nInfected, infectCol(scanner.Text())...)
			}
		}
	}
	if len(nInfected) < 62 {
		return
	}
	if debugOn {
		fmt.Println("infected len:", len(nInfected))
	}
	if err := w.Write(nInfected[:(l + 1)]); err != nil {
		log.Fatalln("error write:", err)
	}

}

func filesToCSV(files []string, w io.Writer, debugOn bool) {
	cw := csv.NewWriter(w)
	header := []string{"Date"}
	header = append(header, TokyoJISCodes()...)
	cw.Write(header)
	if debugOn {
		fmt.Println("tokyo jis len:", len(TokyoJISCodes()))
	}
	for _, f := range files {
		scatterTable(f, len(TokyoJISCodes()), cw, debugOn)
	}
	cw.Flush()
}

func main() {
	debugOn := flag.Bool("d", false, "debug on")
	csv := flag.Bool("c", false, "csv out")
	flag.Parse()
	makePdfTbl()
	if !*csv {
		retrieve(*debugOn)
	} else if *csv {
		files := flag.Args()
		filesToCSV(files, os.Stdout, *debugOn)
	}
}
