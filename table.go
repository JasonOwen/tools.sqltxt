package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

func printTable(dt dataTable) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	RowTxt := ""
	DividerTxt := ""
	for _, columnName := range dt.Columns {
		RowTxt = fmt.Sprintf("%s%s\t", RowTxt, columnName)
		DividerTxt = fmt.Sprintf("%s%s\t", DividerTxt, `-----`)
	}
	fmt.Fprintln(w, RowTxt)
	fmt.Fprintln(w, DividerTxt)
	for rowNum := 0; rowNum < dt.Rows; rowNum++ {
		RowTxt = ""
		for _, columnName := range dt.Columns {
			RowTxt = fmt.Sprintf("%s%s\t", RowTxt, dt.RowData[columnName].value[rowNum])
		}
		fmt.Fprintln(w, RowTxt)
	}

	w.Flush()
}

func setupTableColumns(dtin map[int]dataRow) dataTable {
	var dt dataTable
	startingRow := 0

	//Setup Column Headers
	if columnHeadersCSV != "" {
		for _, colCSVItem := range strings.Split(columnHeadersCSV, ",") {
			dt.Columns = append(dt.Columns, colCSVItem)
		}
	} else {
		if firstLineColumnHeaders {
			for rowID, row := range dtin {
				if rowID > 0 || delimiterMethod {
					dt.Columns = append(dt.Columns, row.value[0])
				}
			}

			//Adjust starting row to read in for data as first line is assumed to be the column headings
			startingRow = 1
		} else if regexMethod {
			reLine := regexp.MustCompile(regexString)
			groupNames := reLine.SubexpNames()

			for grpIndex, grpString := range groupNames {
				if grpIndex > 0 {
					name := grpString
					if name == "" {
						name = fmt.Sprintf("COL_%d", grpIndex)
					}
					dt.Columns = append(dt.Columns, name)
				}
			}
		} else {
			for rowID := range dtin {
				if rowID > 0 || delimiterMethod {
					dt.Columns = append(dt.Columns, fmt.Sprintf("COL_%d", rowID))
				}
			}
		}
	}

	//Load Data Rows Into Table
	//Get Maximum Row Count
	dt.Rows = 0
	for _, col := range dtin {
		if dt.Rows < len(col.value) {
			dt.Rows = len(col.value)
		}
		dt.RowData = make(map[string]dataRow)
	}

	for rowID := startingRow; rowID < dt.Rows; rowID++ {
		for colID, colName := range dt.Columns {
			if _, ok := dt.RowData[colName]; !ok {
				dtBlank := dataRow{}
				dtBlank.value = make(map[int]string)
				dt.RowData[colName] = dtBlank
			}

			colIDAdjusted := colID
			if regexMethod {
				colIDAdjusted = colID + 1
			}
			if _, ok := dtin[colIDAdjusted].value[rowID]; ok {
				dt.RowData[colName].value[rowID] = dtin[colIDAdjusted].value[rowID]
			} else {
				dt.RowData[colName].value[rowID] = ""
			}
		}
	}

	return dt
}

func loadDataBlock(readInString string) map[int]dataRow {
	dt := make(map[int]dataRow)

	//Split input string in by line
	dataRows := strings.Split(readInString, "\n")
	rowNum := 0
	for i := 0; i < len(dataRows); i++ {
		if strings.Trim(dataRows[i], " ") != "" {

			//delimiter is true by default, if regex also specified assume that it is regex method
			if regexMethod {
				//Read block in via regex method
				reLine := regexp.MustCompile(regexString)
				dataRowLine := reLine.FindAllStringSubmatch(dataRows[i], -1)

				for _, matchset := range dataRowLine {
					for j, v := range matchset {
						if j > 0 {
							//If Column doesn't exist prime it
							if _, ok := dt[j]; !ok {
								dtBlank := dataRow{}
								dtBlank.value = make(map[int]string)
								dt[j] = dtBlank
							}
							dt[j].value[rowNum] = v
						}
					}
					rowNum++
				}
			} else {
				// Read block in via delimiter method
				//dataRowLine := strings.Split(dataRows[i], delimiterString)
				dataRowLine := regSplit(dataRows[i], delimiterString)
				for j := 0; j < len(dataRowLine); j++ {
					if _, ok := dt[j]; !ok {
						dtBlank := dataRow{}
						dtBlank.value = make(map[int]string)
						dt[j] = dtBlank
					}
					dt[j].value[rowNum] = dataRowLine[j]
				}
				rowNum++
			}
		}
	}
	return dt
}

func regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}
