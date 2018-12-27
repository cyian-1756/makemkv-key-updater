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
	fmt.Println("[*] Getting webpage")
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
	fmt.Println("[*] Extracting key")
	key := doc.Find(".codebox pre code").Text()
	newFile := ""
	fileContainsAppKey := false
	// If the file doesn't exist we skip reading it and set it to it's default values
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		newFile = "#\n# MakeMKV settings file, written by MakeMKV v1.10.8 linux(x64-release)\n#\napp_BackupDecrypted = \"1\"\n"
	} else {
		dat, err := ioutil.ReadFile(getHomeDir() + "/.MakeMKV/settings.conf")
		if err != nil {
			log.Fatal(err)
		}
		f := strings.Split(string(dat), "\n")
		for _, element := range f {
			if strings.HasPrefix(element, "app_Key") {
				fileContainsAppKey = true
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

			} else if element != "\n" {
				newFile += element + "\n"
			}
		}
	}
	// If fileContainsAppKey is false then the settings file did not contain a app_Key key and we append it
	if !fileContainsAppKey {
		newFile += "app_Key = \"" + key + "\""
	}
	fmt.Println("[*] Writing file")
	ioutil.WriteFile(getHomeDir()+"/.MakeMKV/settings.conf", []byte(newFile), 0644)
}

func main() {
	ExampleScrape()
}
