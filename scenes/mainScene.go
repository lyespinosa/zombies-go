package scenes

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/oakmound/oak/v4/render/mod"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/collision/ray"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

var (
	playerX *float64
	playerY *float64

	destroy = event.RegisterEvent[struct{}]()

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
	char := createPlayer(ctx)
	go playerBehavior(ctx, char)
	go camera(ctx, char)
	go enemyGenerator(ctx)

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

func createPlayer(ctx *scene.Context) *entities.Entity {
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

func playerBehavior(ctx *scene.Context, char *entities.Entity) {
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
			event.TriggerForCallerOn(ctx, hit.Zone.CID, destroy, struct{}{})
		}
		ctx.DrawForTime(
			render.NewLine(x, y, mx, my, color.RGBA{0, 128, 0, 128}),
			time.Millisecond*50,
			1, 2)
		return 0
	})
}

func camera(ctx *scene.Context, char *entities.Entity) {
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

func enemyGenerator(ctx *scene.Context) {
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
	enemy := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 30, 45)),
		entities.WithRenderable(enemyR),
		entities.WithDrawLayers([]int{1, 2}),
		entities.WithLabel(Enemy),
	)

	event.Bind(ctx, event.Enter, enemy, func(e *entities.Entity, ev event.EnterPayload) event.Response {
		x, y := enemy.X(), enemy.Y()
		pt := floatgeom.Point2{x, y}
		pt2 := floatgeom.Point2{*playerX, *playerY}
		delta := pt2.Sub(pt).Normalize().MulConst(EnemySpeed * ev.TickPercent)
		enemy.Shift(delta)

		swtch := enemy.Renderable.(*render.Switch)
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

	event.Bind(ctx, destroy, enemy, func(e *entities.Entity, nothing struct{}) event.Response {
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
