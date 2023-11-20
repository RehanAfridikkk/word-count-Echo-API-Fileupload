package cmd

import (
	"fmt"

	"mime/multipart"
	"time"

	"github.com/RehanAfridikkk/word-count-Echo-API/pkg"
)

func ProcessFile(fileHeader *multipart.FileHeader, routines int) (pkg.CountsResult, int, time.Duration ,error) {
	start := time.Now()
	results := make(chan pkg.CountsResult, routines)

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return pkg.CountsResult{}, 0, 0, err
	}
	defer file.Close()

	// Calculate the size of the file
	fileSize := int64(fileHeader.Size)

	chunkSize := fileSize / int64(routines)

	for i := 0; i < routines; i++ {
		chunk := make([]byte, chunkSize)
		_, err := file.Read(chunk)
		if err != nil {
			return pkg.CountsResult{}, 0, 0, err
		}

		go pkg.Counts(chunk, results)
	}

	totalCounts := pkg.CountsResult{}

	for i := 0; i < routines; i++ {
		result := <-results
		totalCounts.LineCount += result.LineCount
		totalCounts.WordsCount += result.WordsCount
		totalCounts.VowelsCount += result.VowelsCount
		totalCounts.PunctuationCount += result.PunctuationCount
	}

	runTime := time.Since(start)
	fmt.Println("Number of lines:", totalCounts.LineCount)
	fmt.Println("Number of words:", totalCounts.WordsCount)
	fmt.Println("Number of vowels:", totalCounts.VowelsCount)
	fmt.Println("Number of punctuation:", totalCounts.PunctuationCount)
	fmt.Println("Run Time:", runTime)

	return totalCounts, routines, runTime, nil
}
