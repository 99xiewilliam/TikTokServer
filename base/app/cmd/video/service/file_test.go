package service

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestFile(t *testing.T) {
	filename := "temp/11a12dsf88.mp4"
	cmd := exec.Command("ffmpeg", "-i", filename, "-y", "-f", "image2", "-ss", "1", "-t", "0.001", "-s", "1080x1920", "temp/111.jpg")
	if err := cmd.Run(); err != nil {
		fmt.Printf("run err:%v\n", err)
	}
}
