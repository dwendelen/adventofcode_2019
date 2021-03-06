package util

import (
	"bufio"
	"io"
	"log"
	"os"
)

func ReadBytes(file string, fn func(data byte)) {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal("Could not read", file, err)
	}
	buffer := make([]byte, 4096)
	for {
		nbRead, err := input.Read(buffer)
		if nbRead == 0 && err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("Could not read", file, err)
		}
		for i := 0; i < nbRead; i++ {
			fn(buffer[i])
		}
	}
}

func ReadLines(file string, fn func(string) error) {
	ReadDelimiter(file, '\n', fn)
}

func ReadCommaSeparated(file string, fn func(string) error) {
	ReadDelimiter(file, ',', fn)
}

func ReadDelimiter(file string, delimiter byte, fn func(string) error) {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal("Could not read "+file, err)
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

	if item[len(item)-1] == delimiter {
		return item[:len(item)-1]
	} else {
		return item
	}
}

func ReadAllLines(file string) []string {
	res := make([]string, 0)
	ReadLines(file, func(in string) error {
		res = append(res, in)
		return nil
	})
	return res
}

func Abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func Abs64(a int64) int64 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func Min(a int64, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max(a int64, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Lcm(a int64, b int64) int64 {
	return (a / Gcd(a, b)) * b
}

func Gcd(a int64, b int64) int64 {
	for {
		if a == b {
			return a
		}

		if a > b {
			a = a - b
		} else {
			b = b - a
		}
	}
}

func Ceil64(a int64, b int64) int64 {
	return (a + b - 1) / b
}
