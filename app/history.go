package main

import (
	"bufio"
	"fmt"
	"os"
)

type History struct {
	id  int
	cmd string
}

var SessionHistory []string
var lastAppendedIdx int

func fetchAllHistory(fileName string) ([]History, error) {

	historyFile, err := os.Open(fileName)
	if err != nil {

		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("history: error reading history: %v", err)
	}
	defer historyFile.Close()

	scanner := bufio.NewScanner(historyFile)
	var histories []History

	for scanner.Scan() {
		histories = append(histories, History{
			id:  len(histories) + 1,
			cmd: scanner.Text(),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return histories, nil
}

func fetchLastNHistory(n int) ([]History, error) {
	histories, err := fetchAllHistory(historyFilePath)
	if err != nil {
		return nil, err
	}

	if n > len(histories) {
		n = len(histories)
	}
	return histories[len(histories)-n:], nil
}

func getSessionHistory() []History {
	var histories []History
	for i, cmd := range SessionHistory {
		histories = append(histories, History{
			id:  i + 1,
			cmd: cmd,
		})
	}
	return histories
}

func loadHistoryOnStartup() error {
	histfile := os.Getenv("HISTFILE")
	if histfile == "" {
		return nil
	}
	// extF, err := os.Open(histfile)
	// if err != nil {
	// 	return err
	// }

	// inMemF, err := os.OpenFile(historyFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return err
	// }

	// defer extF.Close()
	// defer inMemF.Close()

	// reader := bufio.NewScanner(extF)
	// writer := bufio.NewWriter(inMemF)
	// for reader.Scan() {
	// 	if reader.Text() == "" {
	// 		continue
	// 	}
	// 	writer.WriteString(reader.Text() + "\n")
	// }
	// writer.Flush()
	// return nil
	return writeToFile(histfile, historyFilePath)
}

func saveHistoryToFile() error {
	histFile := os.Getenv("HISTFILE")

	if histFile == "" {
		return nil
	}

	// inMemF, err := os.Open(historyFilePath)
	// if err != nil {
	// 	return err
	// }
	// defer inMemF.Close()

	// extF, err := os.OpenFile(histFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer extF.Close()

	// reader := bufio.NewScanner(inMemF)
	// writer := bufio.NewWriter(extF)

	// for reader.Scan() {
	// 	line := reader.Text()
	// 	if line != "" {
	// 		fmt.Fprintln(writer, line)
	// 	}
	// }
	// writer.Flush()
	// return nil
	return writeToFile(historyFilePath, histFile)
}

func writeToFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(dstFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer src.Close()
	defer dst.Close()

	srcScanner := bufio.NewScanner(src)
	dstWriter := bufio.NewWriter(dst)

	for srcScanner.Scan() {
		if srcScanner.Text() == "" {
			continue
		}
		dstWriter.WriteString(srcScanner.Text() + "\n")
	}
	dstWriter.Flush()
	return nil
}
