package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Video            *Video    `json:"video" valid:"-"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdateAt         time.Time `json:"updated_at" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewJob(output, status string, video *Video) (*Job, error) {
	job := Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}

	if err := job.prepare(); err != nil {
		return nil, err
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *Job) Validate() error {
	_, err := govalidator.ValidateStruct(j)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) prepare() error {
	j.ID = uuid.NewV4().String()
	j.CreatedAt = time.Now()
	j.UpdateAt = time.Now()

	return nil
}
