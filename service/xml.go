package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Искомые поля
const (
	Container       = "ESADout_CUGoods"
	GoodsNumeric    = "catESAD_cu:GoodsNumeric"
	GoodsDescripton = "catESAD_cu:GoodsDescription"
	Brutto          = "catESAD_cu:GrossWeightQuantity"
	GoodCost        = "catESAD_cu:InvoicedCost"
	Tamozhen        = "catESAD_cu:CustomsCost"
	Manufacturer    = "catESAD_cu:Manufacturer"
	Model           = "catESAD_cu:GoodsModel"
	TradeMark       = "catESAD_cu:TradeMark"
	Quantity        = "catESAD_cu:GoodsQuantity"
	EdIzmer         = "catESAD_cu:MeasureUnitQualifierName"
)

func FindXMLFiles(dir string) ([]string, error) {
	var xmlFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))

		if !info.IsDir() && ext == ".xml" {
			xmlFiles = append(xmlFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return xmlFiles, nil
}

func ReadXMLFile(filename string) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	xmlContent := string(byteValue)

	var containers []string

	startContainer := fmt.Sprintf("<%s>", Container)
	endContainer := fmt.Sprintf("</%s>", Container)

	for startIndex := strings.Index(xmlContent, startContainer); startIndex != -1; startIndex = strings.Index(xmlContent, startContainer) {
		endIndex := strings.Index(xmlContent, endContainer)
		if endIndex == -1 {
			break
		}

		containerData := xmlContent[startIndex+len(startContainer) : endIndex]
		containers = append(containers, containerData)

		xmlContent = xmlContent[endIndex+len(endContainer):]
	}

	constants := []string{GoodsNumeric, GoodsDescripton, Brutto, GoodCost, Tamozhen, Manufacturer, Model, TradeMark, Quantity, EdIzmer}
	type XMLMap map[string]string
	XMLMapSlice := []XMLMap{}

	for _, container := range containers {
		dataMap := XMLMap{}

		for _, tag := range constants {
			startTag := fmt.Sprintf("<%s>", tag)
			endTag := fmt.Sprintf("</%s>", tag)

			startIndex := strings.Index(container, startTag)
			endIndex := strings.Index(container, endTag)

			if startIndex != -1 && endIndex != -1 {
				value := container[startIndex+len(startTag) : endIndex]
				dataMap[tag] = value
			}
		}

		XMLMapSlice = append(XMLMapSlice, dataMap)
	}

	var result []map[string]string
	for _, xmlMap := range XMLMapSlice {
		result = append(result, map[string]string(xmlMap))
	}

	return result, nil
}
