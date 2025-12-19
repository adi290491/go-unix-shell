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
