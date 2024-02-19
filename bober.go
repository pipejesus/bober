package main

import (
	"errors"
	"fmt"
	"github.com/icza/mjpeg"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var svgConfig SVGAnimation

func serveSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	svgAnimatedRecord(w, svgConfig)
}

func saveSVGToFile(filePath string) {
	file, _ := os.Create(filePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	svgAnimatedRecord(file, svgConfig)
}

func getInputFileNames() (picPath string, wavPath string, outPath string, err error) {
	if len(os.Args) < 4 {
		return "", "", "", errors.New("Usage: bober picture.jpg sound.wav out.mp4")
	}

	picPath = os.Args[1]
	wavPath = os.Args[2]
	outPath = os.Args[3]
	err = nil
	return
}

func main() {
	picPath, audioPath, outPath, err := getInputFileNames()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	svgConfig = SVGAnimation{
		Size:        300,
		PicFilePath: picPath,
		Duration:    time.Second * 5,
		BgColor:     "#111111",
		FgColor:     "#222222",
		Rocking:     6,
	}

	saveSVGToFile("test.svg")

	http.HandleFunc("/", serveSVG)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Failed to start the HTTP server")
			os.Exit(0)
		}
	}()

	animLength := 5
	fps := 30
	framesToGrab := animLength * fps

	takeScreenshots(animLength, fps)

	aw, _ := mjpeg.New("test.avi", int32(svgConfig.Size), int32(svgConfig.Size), int32(fps))

	audioLength, err := getWavLength(audioPath)
	if err != nil {
		log.Fatal(err)
	}

	movieLength := int(audioLength.Seconds())
	for j := 0; j < movieLength; j = j + animLength {
		for i := 0; i < framesToGrab; i++ {
			data, _ := os.ReadFile(fmt.Sprintf("frames/frame-%03d.jpg", i))
			_ = aw.AddFrame(data)
		}
	}

	_ = aw.Close()

	cmd := exec.Command("ffmpeg", "-y", "-i", "test.avi", "-i", audioPath, "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", outPath)
	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
