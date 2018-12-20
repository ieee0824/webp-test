package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/chai2010/webp"
)

func convertWebp(info os.FileInfo, semaphore chan bool) error {
	name := info.Name()
	f, err := os.Open(fmt.Sprintf("testdata/%s", name))
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	w, err := os.Create(fmt.Sprintf("webp/%s", strings.Split(name, ".")[0]+".webp"))
	if err != nil {
		return err
	}
	defer w.Close()

	return webp.Encode(w, img, &webp.Options{Lossless: true})
}

func convertPNG(info os.FileInfo, semaphore chan bool) error {
	name := info.Name()
	f, err := os.Open(fmt.Sprintf("testdata/%s", name))
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	w, err := os.Create(fmt.Sprintf("png/%s", strings.Split(name, ".")[0]+".png"))
	if err != nil {
		return err
	}
	defer w.Close()

	return png.Encode(w, img)
}

func convertJPEG(info os.FileInfo, semaphore chan bool) error {
	name := info.Name()
	f, err := os.Open(fmt.Sprintf("testdata/%s", name))
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	w, err := os.Create(fmt.Sprintf("jpeg/%s", strings.Split(name, ".")[0]+".jpg"))
	if err != nil {
		return err
	}
	defer w.Close()

	return jpeg.Encode(w, img, &jpeg.Options{
		Quality: 100,
	})
}

func main() {
	infos, err := ioutil.ReadDir("testdata")
	if err != nil {
		log.Fatalln(err)
	}

	semaphore := make(chan bool, runtime.NumCPU())

	var wg sync.WaitGroup
	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		wg.Add(1)
		wg.Add(1)
		wg.Add(1)
		go func(info os.FileInfo) {
			semaphore <- true
			defer func() {
				wg.Done()
				<-semaphore
			}()
			if err := convertWebp(info, semaphore); err != nil {
				log.Println(err)
			}
		}(info)

		go func(info os.FileInfo) {
			semaphore <- true
			defer func() {
				wg.Done()
				<-semaphore
			}()
			err := convertPNG(info, semaphore)
			if err != nil {
				log.Println(err)
			}
		}(info)
		go func(info os.FileInfo) {
			semaphore <- true
			defer func() {
				wg.Done()
				<-semaphore
			}()
			err := convertJPEG(info, semaphore)
			if err != nil {
				log.Println(err)
			}
		}(info)

	}
	wg.Wait()
}
