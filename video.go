package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type orientation string

const (
	LANDSCAPE orientation = "landscape"
	POTRAIT   orientation = "portrait"
	OTHER     orientation = "other"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %w, stderr: %s", err, stderr.String())
	}

	type Stream struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	type ProbeResult struct {
		Streams []Stream `json:"streams"`
	}

	var streams ProbeResult

	if err := json.Unmarshal(stdout.Bytes(), &streams); err != nil {
		return "", fmt.Errorf("error unmarshaling to buffer: %w", err)
	}

	floatW := float32(streams.Streams[0].Width)
	floatH := float32(streams.Streams[0].Height)

	var orient orientation

	num := floatW / floatH
	denom := floatH / floatW

	if 1.5 < num && num < 1.9 && 0.4 < denom && denom < 0.6 {
		orient = LANDSCAPE
	} else if 1.5 < denom && denom < 1.9 && 0.4 < num && num < 0.6 {
		orient = POTRAIT
	} else {
		orient = OTHER
	}

	return string(orient), nil
}

func processVideoForFastStart(filePath string) (string, error) {
	outPath := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outPath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg  error: %w, stderr: %s", err, stderr.String())
	}

	return outPath, nil
}
