package repositories

import (
	"fmt"

	"github.com/Twsouza/codeflix-encoder/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDb struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDb {
	return &VideoRepositoryDb{Db: db}
}

// Insert inserts a new video into the database and returns the inserted video or an error if it fails.
func (vr VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	err := vr.Db.Create(video).Error
	if err != nil {
		return nil, err
	}

	return video, nil
}

// Find retrieves a video from the database by its ID.
// It returns a pointer to the video and an error, if any.
func (vr VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video
	vr.Db.Preload("Jobs").First(&video, "id = ?", id)
	if video.ID == "" {
		return nil, fmt.Errorf("video does not exist")
	}

	return &video, nil
}
