package models

import (
	"image/color"
	"time"

	"github.com/oakmound/oak/v4/render/mod"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/collision/ray"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

const (
	fieldWidth  = 1000
	fieldHeight = 1000
)

var (
	playerX *float64
	playerY *float64
)

func CreatePlayer(ctx *scene.Context) *entities.Entity {
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

func PlayerBehavior(ctx *scene.Context, char *entities.Entity) {
	event.Bind(ctx, event.Enter, char, func(char *entities.Entity, ev event.EnterPayload) event.Response {
		if oak.IsDown(key.W) {
			char.Delta[1] += (-char.Speed.Y() * ev.TickPercent)
		}
		if oak.IsDown(key.A) {
			char.Delta[0] += (-char.Speed.X() * ev.TickPercent)
		}
		if oak.IsDown(key.S) {
			char.Delta[1] += (char.Speed.Y() * ev.TickPercent)
		}
		if oak.IsDown(key.D) {
			char.Delta[0] += (char.Speed.X() * ev.TickPercent)
		}

		if char.X() < 0 {
			char.SetX(0)
		} else if char.X() > fieldWidth-char.W() {
			char.SetX(fieldWidth - char.W())
		}
		if char.Y() < 0 {
			char.SetY(0)
		} else if char.Y() > fieldHeight-char.H() {
			char.SetY(fieldHeight - char.H())
		}

		hit := char.HitLabel(Enemy)
		if hit != nil {
			ctx.Window.NextScene()
		}

		swtch := char.Renderable.(*render.Switch)
		if char.Delta.X() > 0 {
			if swtch.Get() == "left" {
				swtch.Set("right")
			}
		} else if char.Delta.X() < 0 {
			if swtch.Get() == "right" {
				swtch.Set("left")
			}
		}

		return 0
	})

	event.Bind(ctx, mouse.Press, char, func(char *entities.Entity, mevent *mouse.Event) event.Response {
		x := char.X() + char.W()/2
		y := char.Y() + char.H()/2
		vp := ctx.Window.Viewport()
		mx := mevent.X() + float64(vp.X())
		my := mevent.Y() + float64(vp.Y())
		ray.DefaultCaster.CastDistance = floatgeom.Point2{x, y}.Sub(floatgeom.Point2{mx, my}).Magnitude()
		hits := ray.CastTo(floatgeom.Point2{x, y}, floatgeom.Point2{mx, my})
		for _, hit := range hits {
			event.TriggerForCallerOn(ctx, hit.Zone.CID, Destroy, struct{}{})
		}
		ctx.DrawForTime(
			render.NewLine(x, y, mx, my, color.RGBA{0, 128, 0, 128}),
			time.Millisecond*50,
			1, 2)
		return 0
	})
}

func Camera(ctx *scene.Context, char *entities.Entity) {
	screenCenter := ctx.Window.Bounds().DivConst(2)
	event.Bind(ctx, event.Enter, char, func(char *entities.Entity, ev event.EnterPayload) event.Response {
		ctx.Window.(*oak.Window).DoBetweenDraws(func() {
			char.ShiftDelta()
			oak.SetViewport(
				intgeom.Point2{int(char.X()), int(char.Y())}.Sub(screenCenter),
			)
			char.Delta = floatgeom.Point2{}
		})
		return 0
	})
}
