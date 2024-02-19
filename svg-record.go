package main

import (
	svg "github.com/ajstarks/svgo"
	"io"
	"time"
)

type SVGAnimation struct {
	Size        int
	PicFilePath string
	Duration    time.Duration
	BgColor     string
	FgColor     string
	Rocking     int
}

func svgAnimatedRecord(w io.Writer, config SVGAnimation) {
	size := config.Size
	recordSize := size - config.Rocking*2
	duration := config.Duration.Seconds()

	rock := int(config.Rocking)

	canvas := svg.New(w)

	canvas.Start(size, size)

	canvas.Def()
	canvas.Pattern("record_label", 0, 0, 1, 1, "object")
	canvas.Image(0, 0, recordSize/4*2, recordSize/4*2, imageFileToBase64(config.PicFilePath), `preserveAspectRatio="xMidYMid meet" id="pic"`)
	canvas.PatternEnd()
	canvas.DefEnd()

	canvas.Gid(`record_group`)
	canvas.Circle(size/2, size/2, recordSize/2, `id="record" fill="#111111"`)

	for i := 0; i < recordSize/2; i += 8 {
		canvas.Circle(size/2, size/2, i, `fill="none" stroke="#222222" stroke-width="0.5"`)
	}

	canvas.Circle(size/2, size/2, recordSize/4, `id="record_center" fill="url(#record_label)"`)
	canvas.Gend()
	canvas.AnimateRotate(`#record_group`, 0, size/2, size/2, 360, size/2, size/2+rock, duration, -1)
	canvas.End()
}
