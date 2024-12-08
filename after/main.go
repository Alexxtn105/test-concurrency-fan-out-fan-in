package main

import (
	"image"
	"strings"
	imageprocessing "test-concurrency-fan-in-fan-out2/image_processing"
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

// resize Метод производит изменение размера картинки. Возвращает канал типа Job
func resize(jobs *[]Job) <-chan Job {
	// создаем буферизированный канал с именем out и длиной по количеству заданий
	out := make(chan Job, len(*jobs))

	// создаем горутину для каждой задачи
	for _, job := range *jobs {
		go func(job Job) {
			job.Image = imageprocessing.Resize(job.Image)
			// помещаем задачу в канал out
			out <- job
		}(job)
	}

	// возвращаем канал
	return out

	// Это было в непараллельном варианте
	//for index := range *jobs {
	//	(*jobs)[index].Image = imageprocessing.Resize((*jobs)[index].Image)
	//}
}

func saveImages(jobs *[]Job) {
	for _, job := range *jobs {
		imageprocessing.WriteImage(job.OutPath, job.Image)
	}
}

// collectJobs Собирает результаты из канала в слайс заданий
func collectJobs(input <-chan Job, imageCnt int) []Job {
	var resizedJobs []Job

	for i := 0; i < imageCnt; i++ {
		// получаем из канала
		job := <-input

		//пишем в слайс результата
		resizedJobs = append(resizedJobs, job)
	}
	return resizedJobs
}

func main() {
	imagePaths := []string{
		"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	jobs := loadImage(imagePaths)

	// Для этой функции будем делать разбиение (Fan-Out) на несколько горутин
	//resize(&jobs)   // <- так было в старом варианте
	out := resize(&jobs) // читаем из канала в переменную out (resize возвращает канал)

	// собираем воедино (Fan-in)
	resizedJobs := collectJobs(out, len(jobs))

	//saveImages(&jobs)
	saveImages(&resizedJobs)
}
