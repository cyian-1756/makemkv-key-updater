package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func ExampleScrape() {
	// Request the HTML page.
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
	key := doc.Find(".codebox pre code").Text()
	newFile := ""
	dat, err := ioutil.ReadFile(getHomeDir() + "/.MakeMKV/settings.conf")
	f := strings.Split(string(dat), "\n")
	for _, element := range f {
		if strings.HasPrefix(element, "app_Key") {
			currentKey := strings.Replace(element, "app_Key", "", -1)
			currentKey = strings.Replace(currentKey, " ", "", -1)
			currentKey = strings.Replace(currentKey, "=", "", -1)
			currentKey = strings.Replace(currentKey, "\"", "", -1)
			if currentKey == key {
				fmt.Println("[*] Already using newest key")
				return
			}

			fmt.Println("[*] Updating key")
			newFile += "app_Key = \"" + key + "\""

		} else {
			newFile += element + "\n"
		}
	}
	fmt.Println("[*] Writing file")
	ioutil.WriteFile(getHomeDir()+"/.MakeMKV/settings.conf", []byte(newFile), 0644)
}

func main() {
	ExampleScrape()
}
