package main

import (
	"embed"
	"io"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"

	_ "image/png"
)

//go:embed all:_assets
var gameAssets embed.FS

// todo:
// * all tiles are hidden by default, unless you stand next to them
// * different trap kinds
// * change hp<->trap damage mechanic

const (
	ActionMoveRight input.Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveUp
)

const (
	ImageNone resource.ImageID = iota
	ImageSugar
	ImageGopher
	ImageRoomTile
	ImageSpikeTrap
	ImageFireTrap
	ImageBearTrap
)

const (
	AudioNone resource.AudioID = iota
	AudioSpikesTrap
	AudioFireTrap
	AudioBearTrap
	AudioVictorySound
	AudioTheme
)

func main() {
	ctx := ge.NewContext(
		ge.ContextConfig{
			Mute:       false,
			FixedDelta: false,
		},
	)
	ctx.WindowTitle = "Sweets hunter"
	ctx.WindowWidth = 640
	ctx.WindowHeight = 480
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open("_assets/" + path)
		if err != nil {
			panic(err)
		}
		return f
	}

	keymap := input.Keymap{
		ActionMoveRight: {input.KeyD, input.KeyRight, input.KeyGamepadRight},
		ActionMoveDown:  {input.KeyS, input.KeyDown, input.KeyGamepadDown},
		ActionMoveLeft:  {input.KeyA, input.KeyLeft, input.KeyGamepadLeft},
		ActionMoveUp:    {input.KeyW, input.KeyUp, input.KeyGamepadUp},
	}
	h := ctx.Input.NewHandler(0, keymap)

	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageSugar:     {Path: "image/sugar.png"},
		ImageGopher:    {Path: "image/gopher.png", FrameWidth: 32},
		ImageRoomTile:  {Path: "image/floor_tile.png"},
		ImageSpikeTrap: {Path: "image/spikes_trap_tile.png", FrameWidth: 32},
		ImageFireTrap:  {Path: "image/fire_trap_tile.png", FrameWidth: 32},
		ImageBearTrap:  {Path: "image/bear_trap_tile.png", FrameWidth: 32},
	}
	for id, info := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, info)
		ctx.Loader.LoadImage(id)
	}

	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioSpikesTrap:   {Path: "audio/spikes.wav"},
		AudioFireTrap:     {Path: "audio/fire.wav"},
		AudioBearTrap:     {Path: "audio/beartrap.wav"},
		AudioVictorySound: {Path: "audio/pickup.wav"},
		AudioTheme:        {Path: "audio/btheme.ogg"},
	}
	for id, info := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, info)
		ctx.Loader.LoadAudio(id)
	}

	gameController := newGameController(h)

	if err := ge.RunGame(ctx, gameController); err != nil {
		panic(err)
	}
}
