package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

type dataRow struct {
	value map[int]string
}

type dataTable struct {
	Columns []string
	RowData map[string]dataRow
	Rows    int
}

type presetsObject struct {
	Queries map[string]string       `json:"Queries"`
	Parser  map[string]parserObject `json:"Parsers"`
}

type parserObject struct {
	ParseMethod string
	ParseString string
}

func readInUserData(address string) {
	//Check if file exists first
	debugMessage("", "Attempting to read Config File - %s", address)
	if fileExists(address) {
		//If exists attempt to open and rad in config file.
		configFile, err := ioutil.ReadFile(address)

		if err != nil {
			debugMessage("", "Error Reading Config File - %s", address)
		} else {
			if (string(configFile)) != "" {
				//If file loads, configure settings into presets object
				m := presetsObject{}
				debugMessage("", "Attempting to understand Config File - %s", address)
				if err := json.Unmarshal(configFile, &m); err != nil {
					debugMessage("Error understanding userdata file.", "Error Unmarshalling JSON File - %s - %v", address, err)
				}
				loadInParserObject(m)
			}
		}
	}
}

func loadInParserObject(p presetsObject) {
	for qn, qv := range p.Queries {
		presets.Queries[qn] = qv
	}

	for pn, pv := range p.Parser {
		if _, ok := presets.Parser[pn]; !ok {
			pEmpty := parserObject{}
			presets.Parser[pn] = pEmpty
		}
		presets.Parser[pn] = pv
	}
}

func loadPresetData() {
	if loadParser != "" || loadSQL != "" {
		if loadParser != "" {
			if _, ok := presets.Parser[loadParser]; ok {
				if (presets.Parser[loadParser].ParseMethod) == "delimiter" {
					delimiterMethod = true
					delimiterString = presets.Parser[loadParser].ParseString
				} else {
					regexMethod = true
					regexString = presets.Parser[loadParser].ParseString
				}
			}
		}
		if loadSQL != "" {
			if _, ok := presets.Queries[loadSQL]; ok {
				queryString = presets.Queries[loadSQL]
			}
		}
	}
}

func savePresetData() {
	if saveParser != "" || saveSQL != "" {

		if saveParser != "" {
			saveObj := parserObject{}
			if delimiterMethod {
				saveObj.ParseMethod = "delimieter"
				saveObj.ParseString = delimiterString
			} else {
				saveObj.ParseMethod = "regex"
				saveObj.ParseString = regexString
			}
			presets.Parser[saveParser] = saveObj

			debugMessage("", "Saving Parser Preset - %s", saveParser)
		}
		if saveSQL != "" {
			presets.Queries[saveSQL] = queryString
			debugMessage("", "Saving SQL Preset - %s", saveSQL)
		}

		if saveObjByte, err := json.Marshal(presets); err != nil {
			fmt.Println("Error Marshalling Presets")
			panic(err)
		} else {
			saveObjByte, _ = json.Marshal(presets)
			ioutil.WriteFile(userDataDirectory, saveObjByte, 0644)
			debugMessage("", "Save Successful", nil)
		}
	}
}

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		debugMessage("", "Cannot identify user home directory")
	}

	return usr.HomeDir
}
