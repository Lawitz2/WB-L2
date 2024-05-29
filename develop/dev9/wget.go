package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

var pageCounter = -1
var r bool
var wg sync.WaitGroup
var once = true

func copyHtml(mut *sync.Mutex, urlstr string, filenamebase string) {
	defer wg.Done()
	fmt.Printf("dialing %s\n", urlstr)
	pageCounter++
	filename := filenamebase + "-" + strconv.Itoa(pageCounter) + ".html"

	response, err := http.Get(urlstr)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		slog.Error(err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		slog.Error(err.Error())
	}

	if r {
		findLinks(mut, file.Name(), filenamebase)
	}
	changeEmbeds(file.Name(), urlstr)
}

func findLinks(mut *sync.Mutex, filename string, filenamebase string) {
	linksMap := make(map[string]struct{})
	file, err := os.Open(filename)
	if err != nil {
		slog.Error(err.Error())
	}
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		slog.Error(err.Error())
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if exists && strings.Contains(href, "http") {
			linksMap[href] = struct{}{}
		}
	})

	if r && once {
		for i := range linksMap {
			wg.Add(1)
			copyHtml(mut, i, filenamebase)
		}
		once = false
	}
}

func changeEmbeds(filename string, urlstr string) {
	box, err := url.Parse(urlstr)
	if err != nil {
		slog.Error(err.Error())
	}
	file, err := os.Open(filename)
	if err != nil {
		slog.Error(err.Error())
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		slog.Error(err.Error())
	}
	doc.Find("link").Each(func(i int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if exists && !strings.Contains(href, "//") {
			link.SetAttr("href", box.Scheme+"://"+box.Host+href)
		}
	})

	doc.Find("img").Each(func(i int, link *goquery.Selection) {
		srcset, exists := link.Attr("srcset")
		if exists {
			link.SetAttr("srcset", "https:"+srcset)
		}
	})
	doc.Find("a").Each(func(i int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if exists && !strings.Contains(href, "//") {
			link.Find("img").Each(func(i int, selection *goquery.Selection) {
				src, exists := link.Attr("src")
				if exists {
					link.SetAttr("href", src)
				}
			})
		}
	})
	file2, _ := os.Create(file.Name())
	goquery.Render(file2, doc.Selection)
	fmt.Printf("saved %s as %s\n", urlstr, file2.Name())
}

func main() {
	mut := &sync.Mutex{}

	rf := flag.Bool("r", false, "recursively collect all pages")
	flag.Parse()
	r = *rf

	urls := make([]string, 3)
	urls[0] = "https://scrapeme.live/shop/"
	urls[1] = "https://example.com"
	urls[2] = "https://www.iana.org/help/example-domains"
	filename := "wget"

	for _, lurl := range urls {
		wg.Add(1)
		go copyHtml(mut, lurl, filename)
	}

	wg.Wait()
}
