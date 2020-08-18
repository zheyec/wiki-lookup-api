package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	// Cwd - current working directory
	Cwd, _ = os.Getwd()
)

// FileExist - checks if a file exists
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// RemoveWithDelay - remove a file after a certain period of time
func RemoveWithDelay(path string, dur time.Duration) {
	timer := time.NewTimer(dur)
	<-timer.C
	fmt.Println("Removing temp file: ", path)
	os.Remove(path)
}

// RemoveAllTemp - remove all files starting with "temp_" in a specific folder
func RemoveAllTemp(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "temp_") {
			err := os.Remove(path + f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
