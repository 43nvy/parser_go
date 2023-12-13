package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

func ToSCV(filename string, mapSlice []map[string]string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла CSV: %v", err)
	}
	defer file.Close()

	// Кодировка зависит от ОС
	var encoder *encoding.Encoder
	var writer *csv.Writer

	switch runtime.GOOS {
	case "windows":
		encoder = charmap.Windows1251.NewEncoder()
		writer = csv.NewWriter(encoder.Writer(file))
		writer.Comma = ';'
	default:
		encoder = unicode.UTF8.NewEncoder()
		writer = csv.NewWriter(encoder.Writer(file))
		writer.Comma = ','
	}

	for index, dataMap := range mapSlice {
		if index < 2 {
			continue
		}

		data := []string{
			dataMap[GoodsNumeric],
			dataMap[GoodsDescripton],
			dataMap[Brutto],
			dataMap[GoodCost],
			mapSlice[0][ContractCode],
			mapSlice[0][ContractRate],
			mapSlice[1][DeliveryCode] + " " + mapSlice[1][DeliveryPlace],
			mapSlice[0]["DocumentCode"],
			dataMap[DocumentNum] + " от " + dataMap[DocumentDate],
			dataMap[Tamozhen],
			dataMap[Manufacturer],
			dataMap[Model],
			dataMap[TradeMark],
			dataMap[Quantity],
			dataMap[EdIzmer]}

		err := writer.Write(data)
		if err != nil {
			return fmt.Errorf("ошибка при записи data в CSV: %v", err)
		}
	}

	writer.Flush()

	err = writer.Error()
	if err != nil {
		return fmt.Errorf("неизвестная ошибка при записи в CSV: %v", err)
	}

	return nil
}

func CreateCSVFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}

	// Кодировка зависит от ОС
	var encoder *encoding.Encoder
	var writer *csv.Writer

	switch runtime.GOOS {
	case "windows":
		encoder = charmap.Windows1251.NewEncoder()
		writer = csv.NewWriter(encoder.Writer(file))
		writer.Comma = ';'
	default:
		encoder = unicode.UTF8.NewEncoder()
		writer = csv.NewWriter(encoder.Writer(file))
		writer.Comma = ','
	}

	header := []string{
		"Номер",
		"Название",
		"Вес брутто(кг)",
		"Цена товара",
		"Валюта",
		"Курс",
		"Условия поставки",
		"Номер ГТД",
		"Номер инвойса",
		"Таможенная стоимость",
		"Производитель",
		"Модель",
		"Торговая марка",
		"Количество",
		"Единица измерения"}

	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("ошибка при записи header в CSV: %v", err)
	}

	writer.Flush()

	file.Close()

	return nil
}
