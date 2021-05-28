package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.sajari.com/docconv"
)

type tableData struct {
	Operation      string
	CryptoAmount   string
	InvestedAmount string
	ReceivedAmount string
	Price          string
	Date           string
}

func (td *tableData) ToSlice() []string {
	data := make([]string, 0)
	data = append(data, td.Operation)
	data = append(data, td.CryptoAmount)
	data = append(data, td.InvestedAmount)
	data = append(data, td.ReceivedAmount)
	data = append(data, td.Price)
	data = append(data, td.Date)
	return data
}

func main() {
	sum := 0

	// create csv output file
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("creating file", err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// write header
	if err := csvWriter.Write([]string{"Operacion", "Cantidad", "Monto invertido", "Monto recibido", "Cotizacion", "Fecha"}); err != nil {
		log.Fatal("writing to file", err)
	}

	err = filepath.Walk("files", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || info.Name() == ".gitkeep" {
			return nil
		}
		if err != nil {
			return err
		}
		sum++
		log.Printf("reading file %s\n", path)
		data := &tableData{}

		res, err := docconv.ConvertPath(path)
		if err != nil {
			panic(err)
		}

		lines := strings.Split(res.Body, "\n")
		spplitedValues := strings.Split(lines[8], " ")

		operation := spplitedValues[0]
		data.Operation = operation + " " + spplitedValues[1]
		data.CryptoAmount = spplitedValues[3]

		var amount, price string
		if len(spplitedValues) == 8 {
			price = spplitedValues[4] + " " + spplitedValues[5]
			amount = spplitedValues[6] + " " + spplitedValues[7]
		} else if len(spplitedValues) == 6 {
			price = spplitedValues[4] + " " + spplitedValues[5]
			spplitedValues = strings.Split(lines[9], " ")
			amount = spplitedValues[0] + " " + spplitedValues[1]
		} else if len(spplitedValues) == 4 {
			spplitedValues = strings.Split(lines[9], " ")
			price = spplitedValues[0] + " " + spplitedValues[1]
			amount = spplitedValues[2] + " " + spplitedValues[3]
		}

		if operation == "Compra" {
			data.InvestedAmount = amount
			data.ReceivedAmount = ""
		} else if operation == "Venta" {
			data.InvestedAmount = ""
			data.ReceivedAmount = amount
		} else {
			return fmt.Errorf("unkown operation: %s", operation)
		}
		data.Price = price
		date := strings.Split(lines[7], " ")[0]
		data.Date = date

		// write data
		values := data.ToSlice()
		if err := csvWriter.Write(values); err != nil {
			log.Fatal("writing to file", err)
		}

		return nil
	})

	log.Printf("processed files: %d\n", sum)

	if err != nil {
		log.Fatal(err)
	}
}
