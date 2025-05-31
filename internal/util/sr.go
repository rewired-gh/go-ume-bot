package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SRScale string

const (
	SR2x SRScale = "2"
	SR3x SRScale = "3"
	SR4x SRScale = "4"
)

type ModelName string

const (
	X4Plus      ModelName = "realesrgan-x4plus"       // Slowest for general use
	X4PlusAnime ModelName = "realesrgan-x4plus-anime" // Modest speed, good for anime
	V3Anime     ModelName = "realesr-animevideov3"    // Fatest for anime
)

type SRPreset struct {
	Scale     SRScale
	ModelName ModelName
}

const (
	PresetNameAnimeFast4x   = "af4"
	PresetNameAnimeFast2x   = "af2"
	PresetNameAnimeNormal4x = "a"
	PresetNameGeneral4x     = "g"
)

var SRPresets = map[string]SRPreset{
	PresetNameAnimeFast4x: {
		Scale:     SR4x,
		ModelName: V3Anime,
	},
	PresetNameAnimeFast2x: {
		Scale:     SR2x,
		ModelName: V3Anime,
	},
	PresetNameAnimeNormal4x: {
		Scale:     SR4x,
		ModelName: X4PlusAnime,
	},
	PresetNameGeneral4x: {
		Scale:     SR4x,
		ModelName: X4Plus,
	},
}

var ErrInvalidPreset = errors.New("invalid preset name")

func GetPresetsList() string {
	return `Available presets:
af4 - Anime Fast 4x (fastest for anime)
af2 - Anime Fast 2x (fastest for anime)
a - Anime Normal 4x (better quality for anime)
g - General 4x (slowest, for general use)`
}

func UpscaleImage(imgPath string, presetName string, config Config) (resultPath string, err error) {
	log.Printf("[DEBUG] UpscaleImage called with imgPath: %s, presetName: %s", imgPath, presetName)

	preset, ok := SRPresets[presetName]
	if !ok {
		log.Printf("[DEBUG] Invalid preset: %s", presetName)
		return "", ErrInvalidPreset
	}
	log.Printf("[DEBUG] Using preset - Scale: %s, Model: %s", preset.Scale, preset.ModelName)

	imgFileName := filepath.Base(imgPath)
	ext := strings.LastIndex(imgFileName, ".")
	var relativeResultPath string
	if ext == -1 {
		relativeResultPath = filepath.Join(config.TmpPath, imgFileName+"_sr_"+presetName+".png")
	} else {
		relativeResultPath = filepath.Join(config.TmpPath, imgFileName[:ext]+"_sr_"+presetName+".png")
	}

	resultPath, err = filepath.Abs(relativeResultPath)
	if err != nil {
		log.Printf("[DEBUG] Failed to get absolute path: %v", err)
		return
	}
	log.Printf("[DEBUG] Result path: %s", resultPath)

	cmd := exec.Command("./realesrgan-ncnn-vulkan", "-i", imgPath, "-o", resultPath, "-n", string(preset.ModelName), "-s", string(preset.Scale))
	cmd.Dir = config.RESRGANPath
	log.Printf("[DEBUG] Running command: %s in directory: %s", cmd.String(), cmd.Dir)
	log.Printf("[DEBUG] Command args: %v", cmd.Args)

	// Capture both stdout and stderr for better debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[DEBUG] Command failed: %v", err)
		log.Printf("[DEBUG] Command output: %s", string(output))
		return
	}
	log.Printf("[DEBUG] Command completed successfully")
	log.Printf("[DEBUG] Command output: %s", string(output))

	// Check if the result file was actually created
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		log.Printf("[DEBUG] Result file not created: %s", resultPath)
		return "", fmt.Errorf("upscaling failed: result file not created")
	}
	log.Printf("[DEBUG] Result file confirmed: %s", resultPath)

	return
}
