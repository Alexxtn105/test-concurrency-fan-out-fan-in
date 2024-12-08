package main

import (
	"image"
	"strings"
	imageprocessing "test-concurrency-fan-in-fan-out/image_processing"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

func loadImage(paths []string) []Job {
	var jobs []Job

	for _, p := range paths {
		job := Job{
			InputPath: p,
			OutPath:   strings.Replace(p, "images/", "images/output/", 1),
		}

		job.Image = imageprocessing.ReadImage(p)
		jobs = append(jobs, job)
	}

	return jobs
}
func resize(jobs *[]Job) {
	for index := range *jobs {
		(*jobs)[index].Image = imageprocessing.Resize((*jobs)[index].Image)
	}
}

func saveImages(jobs *[]Job) {
	for _, job := range *jobs {
		imageprocessing.WriteImage(job.OutPath, job.Image)
	}
}

func main() {
	imagePaths := []string{
		"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	jobs := loadImage(imagePaths)
	resize(&jobs)
	saveImages(&jobs)
}
