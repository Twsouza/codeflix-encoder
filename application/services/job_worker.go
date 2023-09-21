package services

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/Twsouza/codeflix-encoder/domain"
	"github.com/Twsouza/codeflix-encoder/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

var Mutex = &sync.Mutex{}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, workerID int) {
	for message := range messageChannel {
		Mutex.Lock()
		err := utils.IsJson(string(message.Body))
		Mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		Mutex.Lock()
		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		Mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		Mutex.Lock()
		jobService.VideoService.Video.ID = uuid.NewV4().String()
		Mutex.Unlock()

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		Mutex.Lock()
		err = jobService.VideoService.InsertVideo()
		Mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		jobService.Job = &domain.Job{}
		jobService.Job.Video = jobService.VideoService.Video
		jobService.Job.OutputBucketPath = os.Getenv("outputBucketName")
		jobService.Job.ID = uuid.NewV4().String()
		jobService.Job.Status = "STARTING"
		jobService.Job.CreatedAt = time.Now()

		Mutex.Lock()
		_, err = jobService.JobRepository.Insert(jobService.Job)
		Mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		err = jobService.Start()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		returnChannel <- returnJobResult(*jobService.Job, &message, nil)
	}
}

func returnJobResult(job domain.Job, message *amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: message,
		Error:   err,
	}

	return result
}
