package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readFile(filename string) (report LoadtestReport, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR An error occured opening of file '%s', error: %v", filename, err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	started := false
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if !started && strings.HasPrefix(line, "{") {
			started = true
		}
		if started {
			if strings.HasPrefix(line, "}") && strings.HasSuffix(line, "---------") {
				lines = append(lines, "}")
				break
			}
			lines = append(lines, line)

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(fmt.Sprintf("ERROR An error occured during line scan for file '%s', error: %v", filename, err))
	}

	output := strings.Join(lines, "\n")

	report = LoadtestReport{}

	if err = json.Unmarshal([]byte(output), &report); err != nil {
		fmt.Println("An error occured during unmarshal:", err)
	}

	return
}

func readDirectory(path string) (reports map[time.Time]LoadtestReport, error error) {
	files, err := filepath.Glob(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR An error occured during scan for files a path '%s', error: %v", path, err))
	}

	log.Println("reading reports from files:")
	for _, filename := range files {
		log.Println(fmt.Sprintf("report: [%s]", filename))
	}

	reports = map[time.Time]LoadtestReport{}

	for _, file := range files {
		var report LoadtestReport
		report, readError := readFile(file)
		if readError != nil {
			log.Fatal(fmt.Sprintf("ERROR Could not read loadtest from file '%s', reason: %v", file, readError))
		}
		reports[report.StartedAt] = report
	}

	return
}
