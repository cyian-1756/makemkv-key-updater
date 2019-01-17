package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getConfigPath() (string, error) {
	makemkvFolder := ""
	if runtime.GOOS == "linux" {
		makemkvFolder = getHomeDir() + "/.MakeMKV/"
	} else if runtime.GOOS == "windows" {
		makemkvFolder = os.Getenv("ProgramFiles(x86)") + "/MakeMKV/"
	}
	// If we can't find the install dir we just error out
	// TODO let the user pass in the install dir
	if _, err := os.Stat(makemkvFolder); os.IsNotExist(err) {
		return "", errors.New("Can't find makemkv install dir!")
	}
	return makemkvFolder + "settings.conf", nil
}

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

// Handle windows werid new lines
func getNewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func writeKeyToDisk(key string) {
	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	newLine := getNewLine()
	// Set up a var for the file we're going to write to disk
	newFile := ""
	fileContainsAppKey := false
	// If the file doesn't exist we skip reading it and set it to it's default values
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("[!] " + configPath + " does not exist")
		fmt.Println("[*] Creating file using default values")
		newFile = fmt.Sprintf("#%s# MakeMKV settings file, written by MakeMKV v1.10.8 linux(x64-release)%s#%sapp_BackupDecrypted = \"1\"%s", newLine, newLine, newLine, newLine)
	} else {
		dat, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
		f := strings.Split(string(dat), newLine)
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

			} else if element != newLine {
				// Add every line we find to the file new
				newFile += element + newLine
			}
		}
	}
	// If fileContainsAppKey is false then the settings file did not contain a app_Key key and we append it
	if !fileContainsAppKey {
		newFile += "app_Key = \"" + key + "\""
	}
	fmt.Println("[*] Writing file")
	ioutil.WriteFile(configPath, []byte(newFile), 0644)
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
	writeKeyToDisk(key)
}

func main() {
	scrape()
}
