package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		log.Println("Missing path agrument.")
		log.Fatal("Usage: hazana-report-visualizer [PATH]...")
	}
	path := formatPath(args[0])
	outputFilePath := "./report.html"
	if len(args) > 1 {
		outputFilePath = args[1]
		if !strings.HasSuffix(outputFilePath, ".html") {
			outputFilePath = outputFilePath + ".html"
		}
	}

	reports, _ := readDirectory(path)
	keys, labels := extractKeysAndLabels(reports)

	meanChart := getMeanChart(labels, keys, reports)
	charts := []NamedChart{meanChart}

	for _, key := range keys {
		if key != "" {
			charts = append(charts, getFullChart(key, labels, reports))
		}
	}

	genWeb(charts, "Loadtest results", outputFilePath)
}

func getMeanChart(labels []string, keys []string, reports map[time.Time]LoadtestReport) NamedChart {
	data := ChartData{}
	data.Labels = labels
	for _, key := range keys {
		line := Dataset{}
		line.Label = key
		for _, datekey := range sortedKeys(reports) {
			mean := reports[datekey].Metrics[key].Latencies.Mean / 1000000
			line.Data = append(line.Data, mean)
		}
		line.BackgroundColor = getRandomColorString() //getColorString(255,0,0)
		data.Datasets = append(data.Datasets, line)
	}
	output := NamedChart{Name: "allMeans", ID: "allMeanResponsesID", Chart: data, Options: ChartOptions{ShowLines: true, Elements: OptionElement{Tension: 0}}}
	return output
}

func getRandomColorString() string {
	return getColorString(rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func getColorString(r int, g int, b int) string {
	if r > 255 || r < 0 || g > 255 || g < 0 || b > 255 || b < 0 {
		log.Fatal(fmt.Sprintf("ERROR color not valid value : RGB r:%v g:%v b:%v ", r, g, b))
	}
	return fmt.Sprintf("rgba(%v, %v, %v, 0.2)", r, g, b)
}

func getFullChart(key string, labels []string, reports map[time.Time]LoadtestReport) NamedChart {
	data := ChartData{}
	data.Labels = labels

	// mean
	lineMean := Dataset{}
	lineMean.Label = "Mean"
	for _, datekey := range sortedKeys(reports) {
		mean := reports[datekey].Metrics[key].Latencies.Mean / 1000000
		lineMean.Data = append(lineMean.Data, mean)
	}
	lineMean.BackgroundColor = "rgba(255, 0, 0, 0.2)"
	data.Datasets = append(data.Datasets, lineMean)

	// 50th
	line50th := Dataset{}
	line50th.Label = "50th"
	for _, datekey := range sortedKeys(reports) {
		mean := reports[datekey].Metrics[key].Latencies.Five0Th / 1000000
		line50th.Data = append(line50th.Data, mean)
	}
	line50th.BackgroundColor = "rgba(0, 255, 0, 0.2)"
	data.Datasets = append(data.Datasets, line50th)

	// 95th
	line95th := Dataset{}
	line95th.Label = "95th"
	for _, datekey := range sortedKeys(reports) {
		mean := reports[datekey].Metrics[key].Latencies.Nine5Th / 1000000
		line95th.Data = append(line95th.Data, mean)
	}
	line95th.BackgroundColor = "rgba(0, 255, 255, 0.2)"
	data.Datasets = append(data.Datasets, line95th)

	// 99th
	line99th := Dataset{}
	line99th.Label = "99th"
	for _, datekey := range sortedKeys(reports) {
		mean := reports[datekey].Metrics[key].Latencies.Nine9Th / 1000000
		line99th.Data = append(line99th.Data, mean)
	}
	line99th.BackgroundColor = "rgba(0, 0, 255, 0.2)"
	data.Datasets = append(data.Datasets, line99th)

	// Max
	lineMax := Dataset{}
	lineMax.Label = "Max"
	for _, datekey := range sortedKeys(reports) {
		mean := reports[datekey].Metrics[key].Latencies.Max / 1000000
		lineMax.Data = append(lineMax.Data, mean)
	}
	lineMax.BackgroundColor = "rgba(255, 0, 255, 0.2)"
	data.Datasets = append(data.Datasets, lineMax)

	output := NamedChart{Name: getName(key), ID: getID(key), Chart: data, Options: ChartOptions{ShowLines: true, Elements: OptionElement{Tension: 0}}}
	return output
}

func getName(key string) template.JS {
	return template.JS(strings.Replace(makeFirstLowerCase(key), "-", "", -1))
}
func getID(key string) template.JS {
	return template.JS(makeFirstLowerCase(fmt.Sprintf("%s%s", strings.Replace(key, "-", "", -1), "ID")))
}
