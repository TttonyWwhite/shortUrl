package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	inputfile, inputerr := os.Open("url.txt")
	var arr []string
	if inputerr != nil {
		log.Fatalln("File open failed")
		return
	}
	defer inputfile.Close()

	inputReader := bufio.NewReader(inputfile)
	for {
		url, err := inputReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		url = "http://localhost:8081/long?longUrl=" + url
		arr = append(arr, url)
	}

	outputFile, outputError := os.OpenFile("searchUrls.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		log.Fatalln("File open failed")
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	for i := 0; i < len(arr); i++ {
		outputWriter.WriteString(arr[i])
	}

	outputWriter.Flush()
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
