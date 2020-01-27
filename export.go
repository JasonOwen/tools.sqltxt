package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func exportToFile(filename string, dt dataTable) {
	switch strings.ToLower(outputMode) {
	case "csv":
		exportToCSVFile(filename, dt)
		break
	default:
		break
	}
}

func exportToCSVFile(filename string, dt dataTable) {
	if len(filename) > 0 {

		//Get lines ready to write
		var outputLines []string
		RowTxt := ""
		for colInd, columnName := range dt.Columns {
			if colInd > 0 {
				RowTxt = fmt.Sprintf("%s,%s", RowTxt, columnName)
			} else {
				RowTxt = columnName
			}
		}
		outputLines = append(outputLines, RowTxt)

		for rowNum := 0; rowNum < dt.Rows; rowNum++ {
			RowTxt = ""
			for colInd, columnName := range dt.Columns {
				CellValue := strings.Replace(dt.RowData[columnName].value[rowNum], `"`, `""`, -1)
				exclusionCharacters, _ := regexp.Match(`[,//n"]`, []byte(CellValue))
				if exclusionCharacters {
					CellValue = fmt.Sprintf("\"%s\"", CellValue)
				}
				if colInd > 0 {
					RowTxt = fmt.Sprintf("%s,%s", RowTxt, CellValue)
				} else {
					RowTxt = CellValue
				}
			}
			outputLines = append(outputLines, RowTxt)
		}

		//Write lines out to file
		fileOut, err := os.Create(filename)
		if err != nil {
			debugMessage("", "Error creating output CSV file: %s", filename)
			fileOut.Close()
			return
		}
		for _, line := range outputLines {
			fmt.Fprintln(fileOut, line)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		err = fileOut.Close()
		if err != nil {
			debugMessage("Error saving CSV output", "Error saving CSV output: %v", err)
			return
		}

		debugMessage("Output Saved", "Output CSV saved to %s", filename)
	}
}
