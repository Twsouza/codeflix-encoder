package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Twsouza/codeflix-encoder/application/repositories"
	"github.com/Twsouza/codeflix-encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) Start() error {
	err := j.changeJobStatus("DOWNLOADING")
	if err != nil {
		return err
	}

	err = j.VideoService.Download(os.Getenv("inputBucketName"))
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("FRAGMENTING")
	if err != nil {
		return err
	}

	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("ENCODING")
	if err != nil {
		return err
	}

	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return err
	}

	err = j.changeJobStatus("FINISHING")
	if err != nil {
		return err
	}

	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("COMPLETED")
	if err != nil {
		return err
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus("UPLOADING")
	if err != nil {
		return err
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = fmt.Sprintf("%s/%s", os.Getenv("localStoragePath"), j.Job.ID)
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(fmt.Errorf(uploadResult))
	}

	return nil
}

func (j *JobService) changeJobStatus(status string) error {
	var err error
	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(mainError error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = mainError.Error()
	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return mainError
}
