package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Worker(id int, jobs <-chan Job, results chan<- Job) {
	for job := range jobs {
		fmt.Println(fmt.Sprintf("worker %d : started job %d", id, job.Id))

		start := time.Now()
		response, err := http.Get(job.Url)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, _ = ioutil.ReadAll(response.Body)
		err = response.Body.Close()

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		job.Time = time.Since(start).Seconds()
		job.Response = response.StatusCode

		results <- job
	}
}