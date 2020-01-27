package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func loadDatatableIntoSQL(dt dataTable, tempDBfile string) (bool, *sql.DB) {
	//Load database file
	database, err := sql.Open("sqlite3", tempDBfile)

	if err != nil {
		fmt.Println(fmt.Sprintf("sql.Open : Error : %s\n", err))
		return false, database
	}

	//Create table
	strColumnsCreate := ""
	strColumnsUpdate := ""
	for _, strColName := range dt.Columns {
		strColumnsCreate = fmt.Sprintf("%s, %s TEXT", strColumnsCreate, strColName)
		if len(strColumnsUpdate) > 0 {
			strColumnsUpdate += ", "
		}
		strColumnsUpdate = fmt.Sprintf("%s%s", strColumnsUpdate, strColName)
	}
	sqlCreateTable := fmt.Sprintf(`CREATE TABLE tbl (ID INTEGER PRIMARY KEY AUTOINCREMENT%s);`, strColumnsCreate)
	// fmt.Println(sqlCreateTable)
	_, err = database.Exec(sqlCreateTable)
	if err != nil {
		fmt.Println(fmt.Sprintf("sql.Create : Error : %s\n", err))
		return false, database
	}

	//Populate from datarow
	strInsertValues := ""
	for i := 0; i < dt.Rows; i++ {
		strAddLine := "("
		for _, columns := range dt.Columns {
			if len(strAddLine) > 1 {
				strAddLine += ", "
			}
			strAddLine += "'" + strings.Replace(dt.RowData[columns].value[i], "'", "\"", -1) + "'"
		}
		strAddLine += ")"
		if len(strInsertValues) > 1 {
			strInsertValues += ", "
		}
		strInsertValues += strAddLine
	}
	if strings.Trim(strInsertValues, ` `) != "" {
		strInsertQuery := fmt.Sprintf("INSERT INTO tbl (%s) VALUES %s;", strColumnsUpdate, strInsertValues)
		_, err = database.Exec(strInsertQuery)
		if err != nil {
			fmt.Println(fmt.Sprintf("sql.UpdateRows Error : %s\n%s", err, strInsertQuery))
			return false, database
		}
	}

	return true, database
}

func queryDB(sql string, db *sql.DB) dataTable {
	var dt dataTable
	if sql == "" {
		return dt
	}

	dt.RowData = make(map[string]dataRow)
	rows, err := db.Query(sql)
	if err != nil {
		debugMessage("Query Failed", "sql.Query : Error : %s\n", err)
		return dt
	}

	dt.Columns, err = rows.Columns()

	rowNum := 0
	for rows.Next() {
		columns := make([]interface{}, len(dt.Columns))
		columnPointers := make([]interface{}, len(dt.Columns))
		for i := 0; i < len(dt.Columns); i++ {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			debugMessage("", "sql.Query : RowError : %s\n", err)
		}

		//Load return data into datatable
		value := ""
		for i := 0; i < len(dt.Columns); i++ {
			if typeof(columns[i]) == "int64" {
				value = fmt.Sprintf("%d", columns[i])
			} else {
				value = fmt.Sprintf("%s", columns[i])
			}
			//If Column doesn't exist prime it
			if _, ok := dt.RowData[dt.Columns[i]]; !ok {
				dtBlank := dataRow{}
				dtBlank.value = make(map[int]string)
				dt.RowData[dt.Columns[i]] = dtBlank
			}
			dt.RowData[dt.Columns[i]].value[rowNum] = value
		}
		rowNum++
	}
	dt.Rows = rowNum

	return dt
}
