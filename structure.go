package main

import (
	"html/template"
	"time"
)

//// Hazana output structures
type LoadtestReport struct {
	StartedAt     time.Time `json:"startedAt"`
	FinishedAt    time.Time `json:"finishedAt"`
	Configuration struct {
		Rps            int    `json:"rps"`
		AttackTimeSec  int    `json:"attackTimeSec"`
		RampupTimeSec  int    `json:"rampupTimeSec"`
		RampupStrategy string `json:"rampupStrategy"`
		MaxAttackers   int    `json:"maxAttackers"`
		OutputFilename string `json:"outputFilename"`
		Verbose        bool   `json:"verbose"`
		Metadata       struct {
			GrpcEndpoint string `json:"grpc_endpoint"`
			LogName      string `json:"log_name"`
			MetricType   string `json:"metric.type"`
			ProjectID    string `json:"project_id"`
			UseSsl       string `json:"use_ssl"`
		} `json:"metadata"`
		DoTimeoutSec int `json:"doTimeoutSec"`
	} `json:"configuration"`
	RunError string                   `json:"runError"`
	Metrics  map[string]MetricsResult `json:"metrics"`
	Failed   bool                     `json:"failed"`
	Output   struct {
	} `json:"output"`
}

type MetricsResult struct {
	Latencies struct {
		Total   int64 `json:"total"`
		Mean    int   `json:"mean"`
		Five0Th int   `json:"50th"`
		Nine5Th int   `json:"95th"`
		Nine9Th int   `json:"99th"`
		Max     int   `json:"max"`
	} `json:"latencies"`
	Earliest    time.Time `json:"earliest"`
	Latest      time.Time `json:"latest"`
	End         time.Time `json:"end"`
	Duration    int64     `json:"duration"`
	Wait        int       `json:"wait"`
	Requests    int       `json:"requests"`
	Rate        float64   `json:"rate"`
	Success     float64   `json:"success"`
	StatusCodes struct {
	} `json:"status_codes"`
	Errors interface{} `json:"errors"`
}

//// html Template structures
type NamedChart struct {
	Name    template.JS
	ID      template.JS
	Chart   ChartData
	Options ChartOptions `json:"options"`
}

//// Chart.js structure
type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	Label           string `json:"label"`
	BackgroundColor string `json:"backgroundColor"`
	BorderColor     string `json:"borderColor"`
	Data            []int  `json:"data"`
}

type ChartOptions struct {
	ShowLines bool          `json:"showLines"`
	Elements  OptionElement `json:"elements"`
}
type OptionElement struct {
	Tension int `json:"tension"`
}
