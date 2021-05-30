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
	CryptoCurrency string
	CryptoAmount   string
	InvestedAmount string
	ReceivedAmount string
	Price          string
	Date           string
}

func (td *tableData) ToSlice() []string {
	data := make([]string, 0)
	data = append(data, td.Operation)
	data = append(data, td.CryptoCurrency)
	data = append(data, td.CryptoAmount)
	data = append(data, td.InvestedAmount)
	data = append(data, td.ReceivedAmount)
	data = append(data, td.Price)
	data = append(data, td.Date)
	return data
}

func main() {
	log.Println("stating parser...")
	sum := 0

	// create csv output file
	file, err := os.Create("output/result.csv")
	if err != nil {
		log.Fatal("creating file", err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// write header
	if err := csvWriter.Write([]string{"Operacion", "Moneda", "Cantidad", "Monto invertido", "Monto recibido", "Cotizacion", "Fecha"}); err != nil {
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
		log.Printf("parsing file %s\n", path)
		data := &tableData{}

		res, err := docconv.ConvertPath(path)
		if err != nil {
			panic(err)
		}

		lines := strings.Split(res.Body, "\n")

		date := strings.Split(lines[26], " ")[0]
		data.Date = date

		data.CryptoCurrency = lines[27]
		operation := strings.Split(lines[28], " ")[0]
		data.Operation = operation
		data.CryptoAmount = strings.Split(lines[29], " ")[1]
		data.Price = lines[30]
		amount := lines[31]

		if operation == "Compra" {
			data.InvestedAmount = amount
			data.ReceivedAmount = ""
		} else if operation == "Venta" {
			data.InvestedAmount = ""
			data.ReceivedAmount = amount
		} else {
			return fmt.Errorf("unkown operation: %s", operation)
		}

		// write data
		values := data.ToSlice()
		if err := csvWriter.Write(values); err != nil {
			log.Fatal("writing to file", err)
		}

		return nil
	})

	log.Printf("processed files: %d\n", sum)
	log.Println("writing output/result.csv")

	if err != nil {
		log.Fatal(err)
	}
}
