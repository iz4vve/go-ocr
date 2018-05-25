/*
 * Copyright (C) 2018 Pietro Mascolo - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the Apache 2.0 license (included in the project,
 * and available at: https://www.apache.org/licenses/LICENSE-2.0).
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/otiai10/gosseract"
	"github.com/schollz/progressbar"

	docopt "github.com/docopt/docopt-go"
)

func main() {
	usage := `Desc.
	Usage:
	  gocr convert DIRECTORY [--target=TARGET]
	  gocr FILE [--target=TARGET]
	  gocr -h | --help
	Arguments:
		DIRECTORY         	Directory containing the pdf files to be converted
		FILE			  	Single pdf file to be converted
	Options:
	  -h --help                     	Show this screen.
	  --target=TARGET					Target directory for results.`

	arguments, _ := docopt.ParseArgs(usage, nil, "1.0")

	targetPath, _ := arguments["DIRECTORY"].(string)
	targetFile, _ := arguments["FILE"].(string)
	target, _ := arguments["--target"].(string)
	convert := arguments["convert"].(bool)

	if !strings.HasSuffix(targetFile, ".png") && len(targetFile) > 1 {
		fmt.Println("OCR only works on .png files")
	}

	targetPath = strings.TrimRight(targetPath, string(os.PathSeparator))
	targetDir := path.Join(targetPath, "output")
	_, targetFileName := filepath.Split(targetFile)

	if target != "" && target != "." {
		targetDir = target
	}

	// client stuff
	client := gosseract.NewClient()
	defer client.Close()

	// single file
	if !convert {
		text, err := ocr(client, targetFile)
		handleOcrResults(targetDir, targetFileName, text, err)
	}

	// whole directory
	// fmt.Println("Directory conversion is not supported yet")
	files := getFiles(targetPath)
	fmt.Printf("Saving images to: %s\n", targetDir)
	fmt.Printf("Processing %d files\n", len(files))

	bar := progressbar.New(len(files))
	for _, file := range files {
		bar.Add(1)

		text, err := ocr(client, file)
		handleOcrResults(targetDir, targetFileName, text, err)
	}
	fmt.Println()
}

// ocr performs the actual recognition
func ocr(client *gosseract.Client, imgPath string) (string, error) {
	client.SetImage(imgPath)
	text, err := client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}

// getFiles retrieves all files in a directory
// the ocr requires files to be in png format
// therefore a filter will exclude everything else
func getFiles(filePath string) []string {

	var paths []string
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println(filePath, err)
	}

	for _, f := range files {
		if path.Ext(f.Name()) == ".png" {
			paths = append(paths, path.Join(filePath, f.Name()))
		}
	}
	return paths
}

// saveresults saves the results of the
func saveResults(targetPath, text string) error {

	dir, _ := filepath.Split(targetPath)
	os.MkdirAll(dir, 0770)

	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}

// handleOcrResults takes care of reporting the results or errors
// occurred in the OCR process
func handleOcrResults(targetDir, targetFileName, text string, err error) {
	if err != nil {
		log.Fatal(err)
	}
	err = saveResults(path.Join(targetDir, targetFileName+".txt"), text)

	if err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}
