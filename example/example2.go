package main

import (
	"fmt"
	"time"
	"github.com/stuarthicks/errpool"
)

func main() {
	g := errpool.Group{}
	g.StartWorkers(4)

	jobs := getWork()

	for _, job := range jobs{
		g.Run(job.DoTask)
	}

	errors := g.Wait()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}

func getWork() []*Job{
	numberOfJobs := 30
	var jobs []*Job

	for i := 0; i < numberOfJobs; i++ {
		i := i
		jobs = append(jobs, &Job{id: i})
	}
	return  jobs
}

type Job struct {
	id int
}

func(j *Job) DoTask() error {
	fmt.Printf("Start-Work:%d\n", j.id)
	// emulate time taken to do some work.
	time.Sleep(2 * time.Second)
	if j.id%2 == 1 {
		return fmt.Errorf("Error we dont like odd numbers:%d: ", j.id)
	}
	return nil
}
