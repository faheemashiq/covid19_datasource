package curlib

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Data struct{
    Cumulative_Test_positive int
    Cumulative_tests_performed int
    Date string
    Discharged int
    Expired int
    Region string
    Still_admitted int
}

func Load(path string) []Data {
	table := make([]Data, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
		d := Data{
			Cumulative_Test_positive: row[0],
			Cumulative_tests_performed:    row[1],
			Date:    row[2],
			Discharged:  row[3],
			Expired: row[4]
			Region: row[5]
			Still_admitted: row[6]
		}
		table = append(table, cd)
	}
	return table
}

func Find(table []Data, filter string) []Data {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]Data, 0)
	filter = strings.ToUpper(filter)
	for _, dat := range table {
		if dat.Region == filter ||
			dat.Date == filter ||
			//strings.Contains(strings.ToUpper(dat.Country), filter) ||
			//strings.Contains(strings.ToUpper(cur.Name), filter) {
			result = append(result, dat)
		}
	}
	return result
}