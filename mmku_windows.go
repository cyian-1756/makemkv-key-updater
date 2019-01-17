package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	registry "golang.org/x/sys/windows/registry"
)

func scrape() {
	// Request the HTML page.
	log.Println("[*] Getting webpage")
	res, err := http.Get("https://www.makemkv.com/forum/viewtopic.php?f=5&t=1053")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the new key
	log.Println("[*] Extracting key")
	key := doc.Find(".codebox pre code").Text()
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HKEY_CURRENT_USER\Software\MakeMKV`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()
	s, _, err := k.GetStringValue("app_Key")
	if err != nil {
		log.Fatal(err)
	}

	if s == key {
		log.Println("Already using latest key")
	} else {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HKEY_CURRENT_USER\Software\MakeMKV`, registry.WRITE)
		if err != nil {
			log.Fatal(err)
		}
		k.SetStringValue("app_key", key)
	}
}

func main() {
	scrape()
}
