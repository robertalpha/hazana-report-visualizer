package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
)

func genWeb(namedCharts []NamedChart, service string, outputFilePath string) {
	const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>
        {{.Title}}
    </title>
    <style>
        body {
            background-color: white;
        }

        div.chart-container {
            margin-top: 50px;
        }
        div.chart-container h2 {
            color: maroon;
            text-align: center;
        }
    </style>
</head>
	<body>
{{range .Charts}}
        <div class="chart-container">
           <h2>{{.Name}}</h2>
           <div class="chart">
               <canvas id="{{.ID}}" width="400" height="200"></canvas>
           </div>
        </div>
{{end}}
        <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.8.0/Chart.min.js"></script>

<script>
{{range .Charts}}
    var {{.Name}} = document.getElementById('{{.ID}}').getContext('2d');
    new Chart({{.Name}}, {
        // The type of chart we want to create
        type: 'line', // also try bar or other graph types

        // The data for our dataset
        data: {{.Chart}},
        options: {{.Options}}
        });
{{end}}
</script>
</body>
</html>`

	check := func(err error) {
		if err != nil {
			log.Fatal(fmt.Sprintf("ERROR An error occured during generation of report, error: %v", err))
		}
	}
	t, err := template.New("webpage").Parse(tpl)
	check(err)

	data := struct {
		Title  string
		Charts []NamedChart
	}{
		Title:  service,
		Charts: namedCharts,
	}

	f, err := os.Create(outputFilePath)
	defer f.Close()
	check(err)

	w := bufio.NewWriter(f)

	err = t.Execute(w, data)
	w.Flush()
	check(err)
}
