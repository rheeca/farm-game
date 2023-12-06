package loader

import (
	"embed"
	"fmt"
	"guion-2d-project3/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type SoundCollection struct {
	BGMFirstTown   *audio.Player
	SFXChangeMap   *audio.Player
	SFXChicken     *audio.Player
	SFXChopTree    *audio.Player
	SFXCloseDoor   *audio.Player
	SFXCow         *audio.Player
	SFXCraft       *audio.Player
	SFXOpenDoor    *audio.Player
	SFXTillSoil    *audio.Player
	SFXWateringCan *audio.Player
}

func NewSoundCollection(EmbeddedAssets embed.FS, assetPath string) SoundCollection {
	audioContext := audio.NewContext(utils.SoundSampleRate)
	return SoundCollection{
		BGMFirstTown:   loadWavFromEmbedded(EmbeddedAssets, assetPath, "first_town.wav", audioContext),
		SFXChangeMap:   loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_change_map.wav", audioContext),
		SFXChicken:     loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_chicken.wav", audioContext),
		SFXChopTree:    loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_chop_tree.wav", audioContext),
		SFXCloseDoor:   loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_close_door.wav", audioContext),
		SFXCow:         loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_cow.wav", audioContext),
		SFXCraft:       loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_craft.wav", audioContext),
		SFXOpenDoor:    loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_open_door.wav", audioContext),
		SFXTillSoil:    loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_till_soil.wav", audioContext),
		SFXWateringCan: loadWavFromEmbedded(EmbeddedAssets, assetPath, "sfx_watering_can.wav", audioContext),
	}
}

func (s *SoundCollection) PlaySound(sound *audio.Player) {
	err := sound.Rewind()
	if err != nil {
		fmt.Println("failed to rewind sound")
	}
	sound.Play()
}

func loadWavFromEmbedded(EmbeddedAssets embed.FS, assetPath, name string, context *audio.Context) (soundPlayer *audio.Player) {
	soundFile, err := EmbeddedAssets.Open(path.Join(assetPath, "sounds", name))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	sound, err := wav.DecodeWithoutResampling(soundFile)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	soundPlayer, err = context.NewPlayer(sound)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return soundPlayer
}
