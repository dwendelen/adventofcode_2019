package util

import (
	"bufio"
	"io"
	"log"
	"os"
)

func ReadLines(file string, fn func(string) error) {
	ReadDelimiter(file, '\n', fn)
}

func ReadCommaSeparated(file string, fn func(string) error) {
	ReadDelimiter(file, ',', fn)
}

func ReadDelimiter(file string, delimiter byte, fn func(string) error) {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal("Could not read " + file, err)
	}

	reader := bufio.NewReader(input)
	for true {
		item, readErr := reader.ReadString(delimiter)
		if readErr != nil && readErr != io.EOF {
			log.Fatal("reading standard input:", readErr)
		}

		if len(item) != 0 {
			fnErr := fn(removeDelimiter(item, delimiter))
			if fnErr != nil {
				log.Fatal("Error handling "+item, readErr)
			}
		}

		if readErr == io.EOF {
			return
		}
	}
}

func removeDelimiter(item string, delimiter byte) string {
	if len(item) == 0 {
		return item
	}

	if item[len(item) - 1] == delimiter {
		return item[:len(item) - 1]
	} else {
		return item
	}
}