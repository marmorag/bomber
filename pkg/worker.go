package pkg

import (
	"fmt"
	"net/http"
	"time"
)

func Worker(id int, jobs <-chan Job, results chan<- Job) {
	for job := range jobs {
		fmt.Println(fmt.Sprintf("worker %d : started job %d", id, job.Id))

		start := time.Now()
		response, _ := http.Get(job.Url)

		job.Time = time.Since(start).Seconds()
		job.Response = response.StatusCode

		results <- job
	}
}