package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getTableNameFromFeedIn(fileInName string) string {
	if strings.Contains(fileInName, ";") {
		tmpDBName := strings.Split(fileInName, ";")[0]
		if len(tmpDBName) > 0 {
			return tmpDBName
		}
	}
	return "tbl"
}

func getFileNameFromFeedIn(fileInName string) string {
	if strings.Contains(fileInName, ";") {
		return strings.Split(fileInName, ";")[1]
	}
	return fileInName
}

func getFeedInString(fileInName string) string {
	readInString := ""

	if fileInName == "" {
		info, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}
		if info.Mode()&os.ModeCharDevice == 0 || info.Size() > 0 {
			readInStringBytes, err := ioutil.ReadAll(os.Stdin)
			readInString = string(readInStringBytes)
			if err != nil {
				panic(err)
			}
		}
	} else {
		readInStringBytes, err := ioutil.ReadFile(fileInName)
		if err != nil {
			panic(err)
		}
		readInString = string(readInStringBytes)
	}

	return readInString
}

func reverseStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseStrings(input[1:]), input[0])
}

func typeof(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case float64:
		return "float64"
	//... etc
	default:
		return fmt.Sprintf("%T", v)
	}
}

func debugMessage(message, verboseMessage string, variables ...interface{}) {
	if message != "" && !boolSilent {
		fmt.Println(message)
	}
	if verboseMade {
		if variables != nil {
			fmt.Println(verboseMessage)
		} else {
			fmt.Println(fmt.Sprintf(verboseMessage, variables))
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
