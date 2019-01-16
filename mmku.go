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

func scrape() {
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

	// Set up a var for the file we're going to write to disk
	newFile := ""
	fileContainsAppKey := false
	// If the file doesn't exist we skip reading it and set it to it's default values
	if _, err := os.Stat(getHomeDir() + "/.MakeMKV/settings.conf"); os.IsNotExist(err) {
		fmt.Println("[!] " + getHomeDir() + "/.MakeMKV/settings.conf " + "does not exist")
		fmt.Println("[*] Creating file using default values")
		newFile = "#\n# MakeMKV settings file, written by MakeMKV v1.10.8 linux(x64-release)\n#\napp_BackupDecrypted = \"1\"\n"
	} else {
		dat, err := ioutil.ReadFile(getHomeDir() + "/.MakeMKV/settings.conf")
		if err != nil {
			log.Fatal(err)
		}
		f := strings.Split(string(dat), "\n")
		// Loop over every line in the file
		for _, element := range f {
			// Get the app key from the file and exit if it's the same as the current app key
			// if it's not then change the app key line.
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
				// Add every line we find to the file new
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
	scrape()
}
