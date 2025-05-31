package util

import (
	"errors"
	"fmt"
	"log/slog"
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
	slog.Debug("UpscaleImage called", "imgPath", imgPath, "presetName", presetName)

	preset, ok := SRPresets[presetName]
	if !ok {
		slog.Debug("Invalid preset", "preset", presetName)
		return "", ErrInvalidPreset
	}
	slog.Debug("Using preset", "scale", preset.Scale, "model", preset.ModelName)

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
		slog.Debug("Failed to get absolute path", "error", err)
		return
	}
	slog.Debug("Result path", "path", resultPath)

	cmd := exec.Command("./realesrgan-ncnn-vulkan", "-i", imgPath, "-o", resultPath, "-n", string(preset.ModelName), "-s", string(preset.Scale))
	cmd.Dir = config.RESRGANPath
	slog.Debug("Running command", "command", cmd.String(), "directory", cmd.Dir)
	slog.Debug("Command args", "args", cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("Command failed", "error", err)
		slog.Debug("Command output", "output", string(output))
		return
	}
	slog.Debug("Command completed successfully")
	slog.Debug("Command output", "output", string(output))
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		slog.Debug("Result file not created", "path", resultPath)
		return "", fmt.Errorf("upscaling failed: result file not created")
	}
	slog.Debug("Result file confirmed", "path", resultPath)
	return
}
