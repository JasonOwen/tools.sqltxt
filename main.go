package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var inputFile, delimiterString, regexString, columnHeadersCSV, queryString, loadParser, loadSQL, saveParser, saveSQL, outputMode, outputFile, tmpFilePrep string
var delimiterMethod, regexMethod, firstLineColumnHeaders, verboseMade, printPresets, boolExport, boolSilent, boolKeep bool
var presets presetsObject
var userDataDirectory = "~/.local/share/sqltxt.conf"

func init() {
	//Read in Parameters
	flag.BoolVar(&delimiterMethod, "d", true, "Use Delimiter Method")
	flag.BoolVar(&regexMethod, "r", false, "Use Regular Expression Extraction Method")
	flag.StringVar(&inputFile, "i", "", "Input text file, prepend <tablename>;/filepath, e.g. tbl2;testfile.csv")
	flag.StringVar(&delimiterString, "ds", "[\\s,\\t]", "Delimiter Seperation String/Character")
	flag.StringVar(&regexString, "rs", "", "Regular Expression Extraction String")
	flag.BoolVar(&firstLineColumnHeaders, "f", false, "Use first line as Column Headers")
	flag.StringVar(&columnHeadersCSV, "c", "", "Column Headers CSV")
	flag.StringVar(&queryString, "q", "SELECT * FROM tbl", "Query SQL Statement (table name [tbl])")
	flag.StringVar(&loadParser, "lp", "", "Load Parser Rule")
	flag.StringVar(&loadSQL, "lsql", "", "Load SQL Query")
	flag.StringVar(&saveParser, "sp", "", "Save Parser Rule")
	flag.StringVar(&saveSQL, "ssql", "", "Save SQL Query")
	flag.BoolVar(&verboseMade, "v", false, "Verbose Messaging")
	flag.BoolVar(&printPresets, "p", false, "Print out preset options")
	flag.BoolVar(&boolSilent, "s", false, "Silent Mode, do not print output (except errors)")
	flag.BoolVar(&boolExport, "x", false, "Export Output")
	flag.StringVar(&outputFile, "xfile", "", "Export Filename")
	flag.StringVar(&outputMode, "xmode", "csv", "Export type (csv)")
	flag.StringVar(&tmpFilePrep, "tmpfile", "", "Temp filename and path (optional)")
	flag.BoolVar(&boolKeep, "keep", false, "Setting will not dispose of the temporary sql database file")

	flag.Parse()

	//Stage User Data
	presets.Queries = make(map[string]string)
	presets.Parser = make(map[string]parserObject)
	userDataDirectory = fmt.Sprintf("%s/.local/share/sqltxt.conf", getUserHomeDir())
}

func main() {
	//Load Prset if selected
	readInUserData("/etc/sqltxt.conf")
	readInUserData(userDataDirectory)
	loadPresetData()

	if printPresets {
		printPresetsDisplay()
	} else {
		// Check if any piped in data or file contents if specified
		tableName := getTableNameFromFeedIn(inputFile)
		fileName := getFileNameFromFeedIn(inputFile)
		readInString := getFeedInString(fileName)

		//Split readInString into an arrayed map
		blankDataTable := loadDataBlock(readInString)

		//Setup Table Columns
		DataTable := setupTableColumns(blankDataTable)

		//Create Temporary DB
		tmpFile, err := ioutil.TempFile(os.TempDir(), "tmpdb.*")
		if err != nil {
			fmt.Println("Error creating temporary database")
			return
		}
		if tmpFilePrep != "" {
			os.Remove(tmpFile.Name())
			tmpFile, err = os.OpenFile(tmpFilePrep, os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Println("Error creating temporary database")
				return
			}
		}

		//Load data into a temporary database
		successfulCreation, tmpDB := loadDatatableIntoSQL(DataTable, tmpFile.Name(), tableName)

		//If successful query the data
		if successfulCreation {
			resultTable := queryDB(queryString, tmpDB)

			//Print Table
			if !boolSilent {
				printTable(resultTable)
			}

			//Process export request
			if boolExport {
				exportToFile(outputFile, resultTable)
			}
		}

		//Save Parser and Query if requested
		savePresetData()

		//Clean up
		if !boolKeep {
			defer os.Remove(tmpFile.Name())
		} else {
			if !boolSilent {
				fmt.Println(tmpFile.Name())
			}
		}
		tmpDB.Close()
	}
}

func printPresetsDisplay() {
	fmt.Println("Preset SQL Statements:")
	for sqlPresentName, sqlStatements := range presets.Queries {
		fmt.Println(fmt.Sprintf("%s:\t%s", sqlPresentName, sqlStatements))
	}
	fmt.Println("\nPreset Parsers:")
	for parserName, parser := range presets.Parser {
		fmt.Println(fmt.Sprintf("%s:", parserName))
		fmt.Println(fmt.Sprintf("%s: %s\n", parser.ParseMethod, parser.ParseString))
	}
}
