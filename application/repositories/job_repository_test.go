package repositories_test

import (
	"testing"
	"time"

	"github.com/Twsouza/codeflix-encoder/application/repositories"
	"github.com/Twsouza/codeflix-encoder/domain"
	"github.com/Twsouza/codeflix-encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	vr := repositories.NewVideoRepository(db)
	vr.Insert(video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	jr := repositories.NewJobRepository(db)
	_, err = jr.Insert(job)
	require.Nil(t, err)

	j, err := jr.Find(job.ID)
	require.Nil(t, err)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	vr := repositories.NewVideoRepository(db)
	vr.Insert(video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	jr := repositories.NewJobRepository(db)
	_, err = jr.Insert(job)
	require.Nil(t, err)

	job.Status = "Complete"
	_, err = jr.Update(job)
	require.Nil(t, err)

	j, err := jr.Find(job.ID)
	require.Nil(t, err)
	require.NotEmpty(t, j.ID)
	require.Equal(t, j.Status, job.Status)
}
