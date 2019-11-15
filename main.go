package main

import (
	"fmt"
	"github.com/marmorag/bomber/pkg"
	"github.com/marmorag/optresolver/pkg/optresolver"
	"os"
	"strconv"
)

var args map[string]string
var err error

func init() {
	resolver := optresolver.OptionResolver{
		Help:    `========== Bomber ==========`,
	}

	resolver.AddOption(optresolver.Option{
		Short:    "r",
		Long:     "request",
		Required: false,
		Type:     optresolver.ValueType,
		Default:  "200",
		Help:     "The number of request to send",
	})

	resolver.AddOption(optresolver.Option{
		Short:    "c",
		Long:     "concurrent",
		Required: false,
		Type:     optresolver.ValueType,
		Default:  "10",
		Help:     "The number of concurrent request to be send",
	})

	resolver.AddOption(optresolver.Option{
		Short:    "h",
		Long:     "host",
		Required: true,
		Type:     optresolver.ValueType,
		Help:     "The host to be targeted",
	})

	args, err = resolver.Parse(os.Args)

	if err != nil {
		fmt.Println(err)
	}
}


func main() {

	host := args["host"]
	requestNum, _ := strconv.Atoi(args["request"])
	workerNum, _ := strconv.Atoi(args["concurrent"])

	jobs := make(chan pkg.Job, requestNum)
	results := make(chan pkg.Job, requestNum)

	for w := 1; w <= workerNum; w++ {
		go pkg.Worker(w, jobs, results)
	}

	for j := 1; j <= requestNum; j++ {
		jobs <- pkg.Job{
			Id:  j,
			Url: host,
		}
	}

	close(jobs)

	for a := 1; a <= requestNum; a++ {
		fmt.Println(<-results)
	}
}