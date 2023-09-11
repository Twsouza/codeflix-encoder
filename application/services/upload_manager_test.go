package services_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Twsouza/codeflix-encoder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("buckettest")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "buckettest"
	videoUpload.VideoPath = fmt.Sprintf("%s/%s", os.Getenv("localStoragePath"), video.ID)

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)
	result := <-doneUpload
	require.Equal(t, result, "upload completed")

	err = videoService.Finish()
	require.Nil(t, err)
}
