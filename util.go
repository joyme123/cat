package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func parsePac(filename string) []string {
	var list []string

	if fileExist(filename) {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal()
		}
		list = strings.Split(string(data), "\n")

	}

	return list
}

func isInPac(url string, pac []string) bool {
	for _, v := range pac {
		matched, err := regexp.Match(v, []byte(url))
		if err == nil && matched {
			return true
		}
	}
	return false
}
