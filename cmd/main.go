package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
	"xmlparser/service"
)

// Шаблон для форматирования времени
const layoutHHMMSS = "15_04_05"

func main() {
	xmlFolder := flag.String("dir", "./data", "Путь до папки, где находятся файлы XML")
	flag.Parse()

	_, err := os.Stat(*xmlFolder)
	if os.IsNotExist(err) {
		err := os.MkdirAll(*xmlFolder, os.ModePerm)
		fmt.Println("Папка data создана, поместите в нее файлы формата xml и перезапустите программу.")

		if err != nil {
			fmt.Println("Возникла ошибка при создании папки data, попробуйте создать её вручную.")
			os.Exit(1)
			return
		}
	} else if err != nil {
		fmt.Println("Возникла ошибка при создании папки data, попробуйте создать её вручную.")
		os.Exit(1)
		return
	}

	xmlFiles, err := service.FindXMLFiles(*xmlFolder)

	if err != nil {
		fmt.Printf("Ошибка поиска файлов: %v\n", err)
		os.Exit(1)
		return
	}

	csvFilename := fmt.Sprintf("result_%s.csv", time.Now().Format(layoutHHMMSS))
	err = service.CreateCSVFile(csvFilename)
	if err != nil {
		fmt.Printf("Возникла ошибка: %v\n", err)
		os.Exit(1)
		return
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, file := range xmlFiles {
		fmt.Printf("Началась обработка файла: %s\n", file)
		wg.Add(1)
		go parseXMLtoCSV(file, csvFilename, &wg, &mutex)
	}

	wg.Wait()
}

func parseXMLtoCSV(filename string, csvFilename string, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	service.ReadXMLFile(filename)
	mapsData, err := service.ReadXMLFile(filename)
	if err != nil {
		fmt.Printf("ошибка чтения в XML файла: %v\n", err)
	}

	mutex.Lock()
	err = service.ToSCV(csvFilename, mapsData)
	if err != nil {
		fmt.Printf("ошибка записи в CSV файл: %v\n", err)
	}
	mutex.Unlock()
	fmt.Printf("Файл '%s' преобразован в '%s'\n", filename, csvFilename)
}
