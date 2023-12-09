package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"unicode"
)

type Operation int

const (
	LINES Operation = iota
	WORDS
	BYTES
	CHARS
)

func (operation Operation) String() string {
	var operationsFlagsInOrder = []string{"-l", "-w", "-c", "-m"}
	return operationsFlagsInOrder[operation]
}

var OperationsOrder = []Operation{LINES, WORDS, BYTES, CHARS}

type Settings struct {
	Filenames  []string
	Operations []Operation
}

func (s *Settings) getOperations() {
	for _, operation := range OperationsOrder {
		if slices.Contains(os.Args[1:], operation.String()) && !slices.Contains(s.Operations, operation) {
			s.Operations = append(s.Operations, operation)
		}
	}
	if len(s.Operations) == 0 {
		s.Operations = append(s.Operations, LINES)
		s.Operations = append(s.Operations, WORDS)
		s.Operations = append(s.Operations, BYTES)
	}
}

func (s *Settings) getFiles() {
	for _, argument := range os.Args[1:] {
		if !strings.HasPrefix(argument, "-") {
			s.Filenames = append(s.Filenames, argument)
		}
	}
}

func (s *Settings) ParseArguments() {
	s.getFiles()
	s.getOperations()
}

func main() {
	settings := &Settings{}
	settings.ParseArguments()

	if len(settings.Filenames) == 0 {
		results, err := GetStat(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccwc: %s\n", err.Error())
		}
		for _, operation := range settings.Operations {
			fmt.Printf("%d\t", results[operation])
		}
		fmt.Println()
		os.Exit(0)
	}

	for _, filename := range settings.Filenames {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccwc: %s\n", err.Error())
			continue
		}
		defer file.Close()
		results, err := GetStat(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccwc: %s\n", err.Error())
			continue
		}
		for _, operation := range settings.Operations {
			fmt.Printf("%d\t", results[operation])
		}
		fmt.Printf("%s\n", filename)
	}
}

func GetStat(reader io.Reader) (map[Operation]int, error) {
	bytesReader := bufio.NewReader(reader)
	results := map[Operation]int{BYTES: 0, WORDS: 0, LINES: 0, CHARS: 0}

	cursorInWord := false

	for {
		char, size, err := bytesReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		results[BYTES] += size
		results[CHARS] += 1
		if char == '\n' {
			results[LINES]++
		}

		charIsPartOfWord := !unicode.IsSpace(char)
		if cursorInWord && !charIsPartOfWord {
			results[WORDS]++
			cursorInWord = false
		} else if !cursorInWord && charIsPartOfWord {
			cursorInWord = true
		}
	}
	if cursorInWord {
		results[WORDS]++
	}
	return results, nil
}
