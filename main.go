package main

import (
	"embed"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"

	"zombies/scenes"
)

//go:embed assets
var assets embed.FS

func main() {
	oak.AddScene("zombies", scene.Scene{Start: scenes.StartGame})
	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)

	oak.SetFS(assets)
	oak.Init("zombies", configure)
}

func configure(c oak.Config) (oak.Config, error) {
	c.BatchLoad = true
	c.Assets.ImagePath = "assets/images"
	return c, nil
}
