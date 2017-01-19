package comparer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

func Compare(original, translation string) error {
	// Create Folders Diff
	originalDir, err := readDir(original)
	if err != nil {
		return err
	}
	for _, f := range originalDir {
		if f.IsDir() {
			err = Compare(filepath.Join(original, f.Name()), filepath.Join(translation, f.Name()))
		} else {
			err = readFiles(filepath.Join(original, f.Name()), filepath.Join(translation, f.Name()))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func readDir(path string) ([]os.FileInfo, error) {
	err := os.Chdir(path)
	if err != nil {
		return nil, err
	}
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	file, err := fi.Readdir(0)
	if err != nil {
		return nil, err
	}
	return file, nil

}

func readFiles(orgF, trltF string) error {
	err := os.Chdir(filepath.Dir(orgF))
	if err != nil {
		return err
	}
	fName := strings.Split(orgF, "/")
	fileName := fName[len(fName)-1]
	orgTags, err := readFile(fileName, filepath.Dir(orgF))
	if err != nil {
		return err
	}
	fName = strings.Split(trltF, "/")
	fileName = fName[len(fName)-1]
	trltTags, err := readFile(fileName, filepath.Dir(trltF))
	if err != nil {
		// create empty file with orgTags
		return nil
	}
	if trltTags == nil {
		return nil
	}
	missingTags := findMissing(orgTags, trltTags)
	if missingTags == nil {
		return nil
	}
	// create missingTags in file
	f, err := os.OpenFile(trltF, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, t := range missingTags {
		if _, err = f.WriteString(fmt.Sprintf("<%s>Add your translation here</%s>\n", t, t)); err != nil {
			return err
		}
	}
	return nil
}

// Errors: cant read file
// tags: wrong file format slice => nil
func readFile(file, path string) ([]string, error) {
	if file[len(file)-3:] != "xml" {
		return nil, nil
	}
	inFile, err := os.Open(filepath.Join(path, file))
	if err != nil {
		return nil, err
	}
	defer inFile.Close()
	tags := []string{}
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		indexStart := strings.Index(line, "<")
		indexEnd := strings.Index(line, ">")
		if (indexStart != -1) && (indexEnd != -1) {
			tags = append(tags, line[indexStart:indexEnd+1])
		}
	}
	return tags, nil
}

// More info: https://gist.github.com/ArxdSilva/7392013cbba7a7090cbcd120b7f5ca31
func findMissing(fileFolderA, fileFolderB []string) []string {
	sort.Strings(fileFolderA)
	sort.Strings(fileFolderB)
	if reflect.DeepEqual(fileFolderA, fileFolderB) {
		return nil
	}
	for i := len(fileFolderA) - 1; i >= 0; i-- {
		for _, vD := range fileFolderB {
			if fileFolderA[i] == vD {
				fileFolderA = append(fileFolderA[:i], fileFolderA[i+1:]...)
				break
			}
		}
	}
	return fileFolderA
}
