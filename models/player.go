package models

import (
	"github.com/oakmound/oak/v4/render/mod"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func CreatePlayer(ctx *scene.Context, playerX *float64, playerY *float64) *entities.Entity {
	zeke, err := render.GetSprite("zeke-left.png")
	dlog.ErrorCheck(err)

	playerR := render.NewSwitch("left", map[string]render.Modifiable{
		"left":  zeke,
		"right": zeke.Copy().Modify(mod.FlipX),
	})
	char := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(100, 100, 12, 35)),
		entities.WithRenderable(playerR),
		entities.WithSpeed(floatgeom.Point2{3, 3}),
		entities.WithDrawLayers([]int{1, 2}),
	)

	playerX = &char.Rect.Min[0]
	playerY = &char.Rect.Min[1]

	return char
}
