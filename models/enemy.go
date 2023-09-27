package models

import (
	"math/rand"

	"github.com/oakmound/oak/v4/render/mod"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

const (
	Enemy        = 1
	EnemyRefresh = 60
	EnemySpeed   = 2
)

var Destroy = event.RegisterEvent[struct{}]()

func EnemyGenerator(ctx *scene.Context) {
	event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
		if enterPayload.FramesElapsed%EnemyRefresh == 0 {
			go NewEnemy(ctx)
		}
		return 0
	})
}

func NewEnemy(ctx *scene.Context) {
	x, y := enemyPos()

	enemyFrame, err := render.GetSprite("zombie-right.png")
	dlog.ErrorCheck(err)
	enemyR := render.NewSwitch("left", map[string]render.Modifiable{
		"left":  enemyFrame.Copy().Modify(mod.FlipX),
		"right": enemyFrame,
	})
	hitbox := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 30, 45)),
		entities.WithRenderable(enemyR),
		entities.WithDrawLayers([]int{1, 2}),
		entities.WithLabel(Enemy),
	)

	event.Bind(ctx, event.Enter, hitbox, func(e *entities.Entity, ev event.EnterPayload) event.Response {
		x, y := hitbox.X(), hitbox.Y()
		pt := floatgeom.Point2{x, y}
		pt2 := floatgeom.Point2{*playerX, *playerY}
		delta := pt2.Sub(pt).Normalize().MulConst(EnemySpeed * ev.TickPercent)
		hitbox.Shift(delta)

		swtch := hitbox.Renderable.(*render.Switch)
		if delta.X() > 0 {
			if swtch.Get() == "left" {
				swtch.Set("right")
			}
		} else if delta.X() < 0 {
			if swtch.Get() == "right" {
				swtch.Set("left")
			}
		}
		return 0
	})

	event.Bind(ctx, Destroy, hitbox, func(e *entities.Entity, nothing struct{}) event.Response {
		e.Destroy()
		return 0
	})
}

func enemyPos() (float64, float64) {
	perimeter := fieldWidth*2 + fieldHeight*2
	pos := int(rand.Float64() * float64(perimeter))
	if pos < fieldWidth {
		return float64(pos), 0
	}
	pos -= fieldWidth
	if pos < fieldHeight {
		return float64(fieldWidth), float64(pos)
	}
	pos -= fieldHeight
	if pos < fieldWidth {
		return float64(pos), float64(fieldHeight)
	}
	pos -= fieldWidth
	return 0, float64(pos)
}
