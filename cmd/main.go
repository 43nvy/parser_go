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
const layoutHHMMSS = "15:04:05.000"

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
		return
	}

	var wg sync.WaitGroup

	for index, file := range xmlFiles {
		fmt.Printf("Началась обработка файла: %s\n", file)
		wg.Add(1)
		go parseXMLtoCSV(file, index+1, &wg)
	}

	wg.Wait()
}

func parseXMLtoCSV(filename string, index int, wg *sync.WaitGroup) {
	defer wg.Done()

	mapSlice, err := service.ReadXMLFile(filename)
	if err != nil {
		panic(err)
	}

	csvFilename := fmt.Sprintf("result%d_%s.csv", index, time.Now().Format(layoutHHMMSS))

	err = service.ToSCV(csvFilename, mapSlice)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Файл '%s' преобразован в '%s'\n", filename, csvFilename)
}
