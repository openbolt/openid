package main

// Arguments:
// - Path to testcases.csv
// - Path to testtemplate
// - Path to output
import (
	"encoding/csv"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type Data struct {
	URI     string
	Case    string
	RexBody string
	RexPath string
	Status  int
}

func main() {
	t := template.Must(template.New("test.tmpl").ParseFiles(os.Args[2]))

	// Read testcases
	fd, _ := os.Open(os.Args[1])
	defer fd.Close()
	c := csv.NewReader(fd)
	c.Comma = ','
	table, _ := c.ReadAll()

	reqfields := numReqfields(table[0])
	header := table[0]

	// Process rows
	datas := []Data{}
	for _, row := range table[1:len(table)] {
		dat := Data{}

		// Generate Request URI
		params := url.Values{}
		for k, v := range row[0:reqfields] {
			if v != "" {
				params.Add(header[k], v)
			}
		}
		dat.URI = params.Encode()

		// Fill other variables
		for k, v := range row[reqfields:len(header)] {
			if v != "" {
				switch header[k] {
				case "-Case":
					dat.Case = v
				case "-RexBody":
					dat.RexBody = v
				case "-RexPath":
					dat.RexPath = v
				case "-Status":
					dat.Status, _ = strconv.Atoi(v)
				}
			}
		}
		datas = append(datas, dat)
	}

	// Render template
	fdo, _ := os.Open(os.Args[3])
	defer fdo.Close()
	if err := t.Execute(fdo, datas); err != nil {
		log.Println("executing template: ", err)
	}

}

// BUG: Returns wrong number
func numReqfields(row []string) int {
	f := 0
	for _, cell := range row {
		if !strings.HasPrefix(cell, "-") {
			f++
		} else {
			break
		}
	}
	return f
}
