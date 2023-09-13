package services

import (
	"encoding/json"
	"os"
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

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, workerID int) {
	for message := range messageChannel {
		err := utils.IsJson(string(message.Body))
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		jobService.VideoService.Video.ID = uuid.NewV4().String()

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, &message, err)
			continue
		}

		jobService.Job.Video = jobService.VideoService.Video
		jobService.Job.OutputBucketPath = os.Getenv("outputBucketName")
		jobService.Job.ID = uuid.NewV4().String()
		jobService.Job.Status = "STARTING"
		jobService.Job.CreatedAt = time.Now()

		_, err = jobService.JobRepository.Insert(jobService.Job)
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
