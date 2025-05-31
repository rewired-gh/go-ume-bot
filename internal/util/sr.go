package util

import (
	"errors"
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
	preset, ok := SRPresets[presetName]
	if !ok {
		return "", ErrInvalidPreset
	}

	imgFileName := filepath.Base(imgPath)
	ext := strings.LastIndex(imgFileName, ".")
	if ext == -1 {
		resultPath = filepath.Join(config.TmpPath, imgFileName+"_sr_"+presetName+".png")
	} else {
		resultPath = filepath.Join(config.TmpPath, imgFileName[:ext]+"_sr_"+presetName+".png")
	}

	cmd := exec.Command("./realesrgan-ncnn-vulkan", "-i", imgPath, "-o", resultPath, "-n", string(preset.ModelName), "-s", string(preset.Scale))
	cmd.Dir = config.RESRGANPath

	if err = cmd.Run(); err != nil {
		return
	}

	return
}
