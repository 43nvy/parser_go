package service

import (
	"encoding/csv"
	"os"
)

func ToSCV(filename string, mapSlice []map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Номер",
		"Название",
		"Вес брутто(кг)",
		"Цена товара",
		"Таможенная стоимость",
		"Производитель",
		"Модель",
		"Торговая марка",
		"Колличество",
		"Единица измерения"}

	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, dataMap := range mapSlice {
		data := []string{
			dataMap[GoodsNumeric],
			dataMap[GoodsDescripton],
			dataMap[Brutto],
			dataMap[GoodCost],
			dataMap[Tamozhen],
			dataMap[Manufacturer],
			dataMap[Model],
			dataMap[TradeMark],
			dataMap[Quantity],
			dataMap[EdIzmer]}

		err = writer.Write(data)
		if err != nil {
			return err
		}
	}

	writer.Flush()

	err = writer.Error()
	if err != nil {
		return err
	}

	return nil
}
