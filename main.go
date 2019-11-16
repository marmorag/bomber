package main

import (
	"fmt"
	"github.com/marmorag/bomber/pkg"
	"github.com/marmorag/optresolver"
	"os"
	"strconv"
)

var args map[string]string
var err error

var target string
var requestNum int
var workerNum int

func init() {
	resolver := optresolver.OptionResolver{
		Description:    ` >> Bomber tool
Provide a simple way to do load testing against web server.`,
	}

	_ = resolver.AddOption(optresolver.Option{
		Short:    "r",
		Long:     "request",
		Required: false,
		Type:     optresolver.ValueType,
		Default:  "200",
		Help:     "The number of request to send",
	})

	_ = resolver.AddOption(optresolver.Option{
		Short:    "c",
		Long:     "concurrent",
		Required: false,
		Type:     optresolver.ValueType,
		Default:  "10",
		Help:     "The number of concurrent request to be send",
	})

	_ = resolver.AddOption(optresolver.Option{
		Short:    "t",
		Long:     "target",
		Required: true,
		Type:     optresolver.ValueType,
		Help:     "The host to be targeted",
	})

	args, err = resolver.Resolve(os.Args)

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	target = args["target"]
	requestNum, _ = strconv.Atoi(args["request"])
	workerNum, _ = strconv.Atoi(args["concurrent"])

	jobs := make(chan pkg.Job, requestNum)
	results := make(chan pkg.Job, requestNum)

	var jobResults []pkg.Job

	fmt.Println(fmt.Sprintf("Spawning workers : %d", workerNum))
	for w := 1; w <= workerNum; w++ {
		go pkg.Worker(w, jobs, results)
	}

	fmt.Printf("Starting job enqueing...")
	for j := 1; j <= requestNum; j++ {
		jobs <- pkg.Job{
			Id:  j,
			Url: target,
		}
	}
	fmt.Println("Done.")

	close(jobs)

	fmt.Printf("Ready to receive results.")
	for a := 1; a <= requestNum; a++ {
		jobResults = append(jobResults, <-results)
	}
	fmt.Println("All results received.")

	processStats(jobResults)
}

func processStats(jobs []pkg.Job) {
	var totalTime float64
	totalSuccess := 0
	maxTime := jobs[0].Time
	minTime := jobs[0].Time

	for _, value := range jobs {
		totalTime += value.Time

		if value.Response == 200 {
			totalSuccess++
		}

		if value.Time > maxTime{
			maxTime = value.Time
		}

		if value.Time < minTime {
			minTime = value.Time
		}
	}

	avgTime := totalTime / float64(len(jobs))
	fmt.Println(fmt.Sprintf("OK Request : %d/%d ; avg %f ; min %f ; max %f", totalSuccess, len(jobs), avgTime, minTime, maxTime))
}