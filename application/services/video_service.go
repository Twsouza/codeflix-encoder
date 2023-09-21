package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/Twsouza/codeflix-encoder/application/repositories"
	"github.com/Twsouza/codeflix-encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.mp4", os.Getenv("localStoragePath"), v.Video.ID))
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("video %v has been stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("localStoragePath"), v.Video.ID), os.ModePerm)
	if err != nil {
		return err
	}

	source := fmt.Sprintf("%s/%s.mp4", os.Getenv("localStoragePath"), v.Video.ID)
	target := fmt.Sprintf("%s/%s.frag", os.Getenv("localStoragePath"), v.Video.ID)

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, fmt.Sprintf("%s/%s.frag", os.Getenv("localStoragePath"), v.Video.ID))
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, fmt.Sprintf("%s/%s", os.Getenv("localStoragePath"), v.Video.ID))
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	cmd := exec.Command("mp4dash", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(fmt.Sprintf("%s/%s.mp4", os.Getenv("localStoragePath"), v.Video.ID))
	if err != nil {
		return err
	}

	err = os.Remove(fmt.Sprintf("%s/%s.frag", os.Getenv("localStoragePath"), v.Video.ID))
	if err != nil {
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("localStoragePath"), v.Video.ID))
	if err != nil {
		return err
	}

	log.Println("files have been removed")
	return nil
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("======> Output: %s\n", string(out))
	}
}
