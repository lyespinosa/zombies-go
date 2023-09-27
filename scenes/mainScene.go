package scenes

import (
	"math/rand"
	"zombies/models"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

var (
	sheet [][]*render.Sprite
)

const (
	Enemy collision.Label = 1

	fieldWidth  = 1000
	fieldHeight = 1000

	EnemyRefresh = 60
	EnemySpeed   = 2
)

func StartGame(ctx *scene.Context) {

	generateGround()

	char := models.CreatePlayer(ctx)
	go models.PlayerBehavior(ctx, char)
	go models.Camera(ctx, char)
	go models.EnemyGenerator(ctx)

}

func generateGround() {
	sprites, err := render.GetSheet("ground.png")
	dlog.ErrorCheck(err)
	sheet = sprites.ToSprites()

	oak.SetViewportBounds(intgeom.NewRect2(0, 0, fieldWidth, fieldHeight))

	for x := 0; x < fieldWidth; x += 16 {
		for y := 0; y < fieldHeight; y += 16 {
			ix := rand.Intn(3)
			iy := rand.Intn(3)
			sp := sheet[iy][ix].Copy()
			sp.SetPos(float64(x), float64(y))
			render.Draw(sp, 1, 1)
		}
	}
}
