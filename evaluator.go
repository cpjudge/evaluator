package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/udhos/equalfile"
	// "io"
	// "bytes"
)

// EvaluateSubmission : Evaluates the submission
func EvaluateSubmission(submissionID string, questionID string) int32 {
	// Path refers to path inside the cpjudge_webserver container because
	// directory is already mounted.
	base := "/media/vaibhav/Coding/go/src/github.com/cpjudge/cpjudge_webserver"
	submissionsPath := base + "/submissions/" + submissionID
	expectedOutputPath := base + "/questions/testcases/" +
		questionID +
		"/output/"

	expectedOutputFiles, err := ioutil.ReadDir(expectedOutputPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range expectedOutputFiles {
		errorFilePath := submissionsPath + "/error/" + f.Name()
		compilationErrorFilePath := submissionsPath + "/compilation_error.err"
		outputFilePath := submissionsPath + "/output/" + f.Name()
		expectedOutputFilePath := expectedOutputPath + f.Name()

		if fileInfo, err := os.Stat(compilationErrorFilePath); err == nil {
			if fileInfo.Size() != 0 {
				return 3 //Compilation Error
			}
		}
		if fileInfo, err := os.Stat(errorFilePath); err == nil {
			if fileInfo.Size() != 0 {
				return 4 //Runtime Error
			}
		}
		cmp := equalfile.New(nil, equalfile.Options{Debug: false})
		if equal, _ := cmp.CompareFile(outputFilePath, expectedOutputFilePath); !equal {
			return 1 // Wrong Answer
		}
	}
	return 0 // Correct Answer
}
