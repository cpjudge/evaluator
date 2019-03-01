// Get Results
// Check output and send back results
// Queue

package main

import (
	"io/ioutil"
	"log"
	"os"
	// "io"
	// "bytes"
	"github.com/udhos/equalfile"
)

func EvaluateSubmission(submissionPath string, expectedOutputPath string) int32 {

	expectedOutputFiles, err := ioutil.ReadDir(expectedOutputPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range expectedOutputFiles {
		errorFilePath := submissionPath + "error/" + f.Name()
		outputFilePath := submissionPath + "output/" + f.Name()
		expectedOutputFilePath := expectedOutputPath + f.Name()

		if fileInfo, _ := os.Stat(errorFilePath); fileInfo.Size() != 0 {
			return 4 //Runtime Error
		}
		cmp := equalfile.New(nil, equalfile.Options{Debug: false})
		if equal, _ := cmp.CompareFile(outputFilePath, expectedOutputFilePath); !equal {
			return 1 // Wrong Answer
		}
		// if !deepCompare(outputFilePath, expectedOutputFilePath)	{
		// 	return 1
		// }
	}

	log.Println("CORRECT")
	return 0 // Correct Answer
}
