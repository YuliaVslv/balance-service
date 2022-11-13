package report

import (
	"encoding/csv"
	"fmt"
	"os"
)

func CreateReportFile(data [][]string, year uint32, month uint32) (string, error) {
	var strMonth string
	if month >= 10 {
		strMonth = fmt.Sprint(month)
	} else {
		strMonth = fmt.Sprint("0", month)
	}
	filename := fmt.Sprint("report-files/", year, "-", strMonth, ".csv")
	file, err := os.Create(filename)

	if err != nil {
		return "", err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		return filename, err
	}
	return filename, nil
}
