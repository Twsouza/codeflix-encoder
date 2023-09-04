package repositories

import (
	"fmt"

	"github.com/Twsouza/codeflix-encoder/domain"
	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(Job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	Db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepositoryDb {
	return &JobRepositoryDb{Db: db}
}

func (jr *JobRepositoryDb) Insert(Job *domain.Job) (*domain.Job, error) {
	err := jr.Db.Create(Job).Error
	if err != nil {
		return nil, err
	}

	return Job, nil
}

func (jr *JobRepositoryDb) Find(id string) (*domain.Job, error) {
	var Job domain.Job
	jr.Db.Preload("Video").First(&Job, "id = ?", id)
	if Job.ID == "" {
		return nil, fmt.Errorf("job does not exist")
	}

	return &Job, nil
}

func (jr *JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := jr.Db.Save(job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}
