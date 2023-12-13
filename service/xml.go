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
	DocumentCode = "<!-- -->"

	ContainerDelivery = "catESAD_cu:CUESADDeliveryTerms"
	DeliveryPlace     = "cat_ru:DeliveryPlace"
	DeliveryCode      = "cat_ru:DeliveryTermsStringCode"

	ContainerContract = "ESADout_CUMainContractTerms"
	ContractCode      = "catESAD_cu:ContractCurrencyCode"
	ContractRate      = "catESAD_cu:ContractCurrencyRate"

	ContainerDocument = "ESADout_CUPresentedDocument"
	DocumentName      = "cat_ru:PrDocumentName" // Инвойс надо
	DocumentNum       = "cat_ru:PrDocumentNumber"
	DocumentDate      = "cat_ru:PrDocumentDate"

	ContainerGoods  = "ESADout_CUGoods"
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

var delivery = []string{ContainerDelivery, DeliveryPlace, DeliveryCode}
var contract = []string{ContainerContract, ContractCode, ContractRate}
var document = []string{ContainerDocument, DocumentName, DocumentNum, DocumentDate}
var goods = []string{ContainerGoods, GoodsNumeric, GoodsDescripton, Brutto, GoodCost, Tamozhen, Manufacturer, Model, TradeMark, Quantity, EdIzmer}

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

	documentCode, err := extractDocumentCode(xmlContent)
	if err != nil {
		documentCode = "Не найден"
	}

	deliveryData, err := onceContainerData(xmlContent, delivery)
	if err != nil {
		return nil, err
	}

	contractData, err := onceContainerData(xmlContent, contract)
	if err != nil {
		return nil, err
	}

	contractData["DocumentCode"] = documentCode

	goodsMaps, err := goodsContainerData(xmlContent, goods)
	if err != nil {
		return nil, err
	}

	goodsMaps = append([]map[string]string{deliveryData}, goodsMaps...)
	goodsMaps = append([]map[string]string{contractData}, goodsMaps...)

	return goodsMaps, nil
}

func goodsContainerData(xmlContent string, goodsTags []string) ([]map[string]string, error) {
	goodsContainers, err := extractAllData(xmlContent, goodsTags[0])
	if err != nil {
		return nil, err
	}

	var goodsMaps []map[string]string

	for _, goodsContainer := range goodsContainers {
		goodsContainerMap := make(map[string]string)

		for _, tag := range goodsTags[1:] {
			value, err := extractData(goodsContainer, tag)
			if err != nil {
				return nil, err
			}

			goodsContainerMap[tag] = value
		}

		documentsContainers, err := extractAllData(goodsContainer, document[0])
		if err != nil {
			return nil, err
		}

		documentMap := make(map[string]string)

		for _, documentContainer := range documentsContainers {
			for _, tag := range document[1:] {
				value, err := extractData(documentContainer, tag)
				if err != nil {
					return nil, err
				}

				documentMap[tag] = value
			}

			if documentMap[DocumentName] == "ИНВОЙС (СЧЕТ-ФАКТУРА) К ДОГОВОРУ" {
				goodsContainerMap[DocumentNum] = documentMap[DocumentNum]
				goodsContainerMap[DocumentDate] = documentMap[DocumentDate]
			}
		}

		goodsMaps = append(goodsMaps, goodsContainerMap)
	}

	return goodsMaps, nil
}

func onceContainerData(xmlContent string, containerTags []string) (map[string]string, error) {
	containerContent, err := extractData(xmlContent, containerTags[0])
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	containerMap := make(map[string]string)

	for index, tag := range containerTags {
		if index == 0 {
			continue
		}

		value, err := extractData(containerContent, tag)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		containerMap[tag] = value
	}

	return containerMap, nil
}

func extractDocumentCode(xmlContent string) (string, error) {
	startTag := "<!--"
	endTag := "-->"

	startIndex := strings.Index(xmlContent, startTag)
	if startIndex == -1 {
		return "", fmt.Errorf("tag not found")
	}

	endIndex := strings.Index(xmlContent[startIndex:], endTag) + startIndex
	data := xmlContent[startIndex+len(startTag) : endIndex]

	data = strings.TrimPrefix(data, "ND=")

	return data, nil
}

func extractData(xmlContent, tag string) (string, error) {
	startTag := fmt.Sprintf("<%s>", tag)
	endTag := fmt.Sprintf("</%s>", tag)

	startIndex := strings.Index(xmlContent, startTag)
	if startIndex == -1 {
		return "", fmt.Errorf("tag '%s' not found", tag)
	}

	endIndex := strings.Index(xmlContent[startIndex:], endTag) + startIndex
	data := xmlContent[startIndex+len(startTag) : endIndex]

	return data, nil
}

func extractAllData(xmlContent, tag string) ([]string, error) {
	startTag := fmt.Sprintf("<%s>", tag)
	endTag := fmt.Sprintf("</%s>", tag)

	var allData []string

	startIndex := 0
	for {
		index := strings.Index(xmlContent[startIndex:], startTag)
		if index == -1 {
			break
		}

		startIndex += index
		endIndex := strings.Index(xmlContent[startIndex:], endTag) + startIndex
		if endIndex == -1 {
			return nil, fmt.Errorf("closing tag '</%s>' not found for tag '%s'", tag, tag)
		}

		data := xmlContent[startIndex+len(startTag) : endIndex]
		allData = append(allData, data)

		// Передвигаем индекс для следующего поиска
		startIndex = endIndex + len(endTag)
	}

	return allData, nil
}
