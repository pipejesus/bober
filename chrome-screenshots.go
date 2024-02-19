package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/disintegration/imaging"
	"log"
	"os"
)

func takeScreenshots(animLength int, fps int) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte

	_ = chromedp.Run(ctx,
		chromedp.Navigate(`http://localhost:8080/`))

	currTime := 0.0
	framesToGrab := animLength * fps

	if err := os.MkdirAll("frames", 0o755); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < framesToGrab; i++ {
		fmt.Println("taking screenshot", i)
		currTime = float64(animLength) / float64(framesToGrab) * float64(i)
		err := chromedp.Run(ctx,
			chromedp.Evaluate(setAnimationTime(currTime), nil),
			chromedp.Screenshot("svg", &buf, chromedp.NodeVisible),
		)

		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(fmt.Sprintf("frames/frame-%03d.png", i), buf, 0o644); err != nil {
			log.Fatal(err)
		}

		img, err := imaging.Open(fmt.Sprintf("frames/frame-%03d.png", i))
		if err != nil {
			log.Fatal(err)
		}
		err = imaging.Save(img, fmt.Sprintf("frames/frame-%03d.jpg", i))
		if err != nil {
			log.Fatal(err)
		}

	}
}
