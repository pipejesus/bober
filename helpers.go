package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func setAnimationTime(time float64) string {
	jsAction := `
		(() => {
			const svg = document.querySelector('svg');
			svg.pauseAnimations();
			svg.setCurrentTime(%current_time%);
		})();
	`

	jsAction = strings.Replace(jsAction, "%current_time%", fmt.Sprintf("%.4f", time), 1)
	return jsAction
}

func imageFileToBase64(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
		return ""
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
		return ""
	}

	//_, format, err := image.DecodeConfig(file)
	//if err != nil {
	//	panic(err)
	//	return ""
	//}
	format := "jpeg"

	str := base64.StdEncoding.EncodeToString(data)

	return "data:image/" + format + ";base64," + str
}
