package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	// App
	appName = "GoLock"

	// Sizes
	buttonRadius   = 90
	buttonSpace    = buttonRadius + 5
	buttonCenter   = buttonRadius + 5
	buttonDiameter = 2 * buttonSpace

	// Math
	pi = math.Pi

	// Draw
	redrawTo = 2 * time.Second

	// Font
	fontFamily = "Hack"
	fontSize   = 16.0
)

var (
	height float64
	width  float64
)

func init() {
	flag.Float64Var(&height, "h", 1080, "window height")
	flag.Float64Var(&width, "w", 1920, "window width")
}

type App struct {
	Window         *gtk.Window
	Ctx            *cairo.Context
	AlphaSupported bool
}

func (a *App) Repaint() {
	screenChanged(&a.Window.Widget)
}

func (a *App) SetMainSourceRGBA() {
	switch state.LastEvent {
	case AuthSuccessEvent:
		a.Ctx.SetSourceRGB(51.0/255, 125.0/255, 0)
	case AuthCheckEvent:
		a.Ctx.SetSourceRGBA(0, 144.0/255, 255.0/255, 0.75)
	case CapsChangedEvent:
		if state.IsCapsLockOn {
			a.Ctx.SetSourceRGBA(0, 144.0/255, 255.0/255, 0.75)
		} else {
			a.Ctx.SetSourceRGBA(0, 144.0/255, 255.0/255, 0.75)
		}
		break
	case EmptyPasswordEvent, WrongPasswordEvent, BackSpaceEvent:
		a.Ctx.SetSourceRGBA(250.0/255, 0, 0, 0.75)
		break
	case KeyPressEvent:
		a.Ctx.SetSourceRGBA(0, 0, 0, 0.75)
		break
	}
}

func (a *App) SetMainSourceRGB() {
	switch state.LastEvent {
	case AuthSuccessEvent:
		a.Ctx.SetSourceRGB(51.0/255, 125.0/255, 0)
	case AuthCheckEvent:
		a.Ctx.SetSourceRGB(51.0/255, 0, 250.0/255)
	case CapsChangedEvent:
		if state.IsCapsLockOn {
			a.Ctx.SetSourceRGB(51.0/255, 0, 250.0/255)
		} else {
			a.Ctx.SetSourceRGB(51.0/255, 0, 250.0/255)
		}
		break
	case EmptyPasswordEvent, WrongPasswordEvent, BackSpaceEvent:
		a.Ctx.SetSourceRGB(125.0/255, 51.0/255, 0)
		break
	case KeyPressEvent:
		a.Ctx.SetSourceRGB(51.0/255, 125.0/255, 0)
		break
	}
}

func (a *App) DrawStateText() {
	// Setup font settings
	a.Ctx.SetSourceRGB(255, 255, 255)
	a.Ctx.SetLineWidth(10.0)
	a.Ctx.SelectFontFace(fontFamily, cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	a.Ctx.SetFontSize(fontSize)

	// Set text
	var text string
	switch state.LastEvent {
	case AuthSuccessEvent:
		text = "Success!"
	case AuthCheckEvent:
		text = "Checking..."
	case EmptyPasswordEvent:
		text = "No input"
	case WrongPasswordEvent:
		text = "Wrong"
	default:
		if state.IsCapsLockOn {
			text = "CapsLock ON"
		}
	}

	if text == "" {
		return
	}

	extents := a.Ctx.TextExtents(text)
	x := width/2 - ((extents.Width / 2) + extents.XBearing)
	y := height/2 - ((extents.Height / 2) + extents.YBearing)

	a.Ctx.MoveTo(x, y)
	a.Ctx.ShowText(text)
	a.Ctx.ClosePath()
}

func (a *App) CreateKeyPressArc() {
	switch state.LastEvent {
	case BackSpaceEvent:
		if state.PasswordLength == 0 {
			return
		}

		fallthrough
	case KeyPressEvent:
		a.Ctx.SetSourceRGB(0, 0, 0)
		a.Ctx.SetLineWidth(10.0)
		startRadians := 2 * pi * 100
		var start int32
		for start == state.LastStart {
			start = (rand.Int31() % int32(startRadians)) / 100
		}

		state.LastStart = start
		endRadians := (pi / 3)
		a.Ctx.Arc(width/2, height/2, buttonRadius, float64(start), float64(start+int32(endRadians)))
		a.Ctx.Stroke()
		break
	}
}

func (a *App) QueueRedraw(willClear bool) {
	state.ShouldDraw = true
	a.Window.QueueDraw()
	if !willClear {
		return
	}

	go func() {
		now := time.Now()
		state.LastQueuedRedraw = now
		time.Sleep(redrawTo)

		if !state.LastQueuedRedraw.Equal(now) {
			return
		}

		state.ShouldDraw = false
		a.Window.QueueDraw()
	}()
}

func (a *App) Draw() {
	if state.LastEvent == NoEvent || !state.ShouldDraw {
		return
	}

	app.Ctx.SetLineWidth(10.0)
	app.Ctx.Arc(width/2, height/2, buttonRadius, 0, 2*pi)
	app.SetMainSourceRGBA()
	app.Ctx.FillPreserve()
	app.SetMainSourceRGB()
	app.Ctx.Stroke()

	app.Ctx.SetSourceRGB(0, 0, 0)
	app.Ctx.SetLineWidth(2.0)
	app.Ctx.Arc(width/2, height/2, buttonRadius-5, 0, 2*pi)
	app.Ctx.Stroke()

	app.CreateKeyPressArc()

	app.DrawStateText()
}

var app *App

func watchState(done chan bool) {
	for {
		select {
		case event := <-state.EventChan:
			log.Println("New Event!")
			state.LastEvent = event
			app.QueueRedraw(true)
		case <-done:
			done <- true
			log.Println("watchState done!")
			return
		}
	}
}

func startApp(done chan bool) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_POPUP)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(appName)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Needed for transparency
	win.SetAppPaintable(true)

	win.Connect("screen-changed", func(widget *gtk.Widget, oldScreen *gdk.Screen, userData ...interface{}) {
		screenChanged(widget)
	})

	win.Connect("draw", func(window *gtk.Window, context *cairo.Context) {
		exposeDraw(window, context)
		app.Draw()
	})

	win.SetDefaultSize(int(width), int(height))
	app = &App{
		Window: win,
	}

	// win.Add(vbox)

	app.Repaint()
	go watchState(done)

	win.ShowAll()
	gtk.Main()
}

func screenChanged(widget *gtk.Widget) {
	screen, _ := widget.GetScreen()
	visual, _ := screen.GetRGBAVisual()

	if visual != nil {
		app.AlphaSupported = true
	} else {
		log.Println("Alpha not supported")
		app.AlphaSupported = false
	}

	log.Printf("%#v\n", visual)

	widget.SetVisual(visual)
}

func exposeDraw(w *gtk.Window, ctx *cairo.Context) {
	if app.AlphaSupported {
		ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.25)
	} else {
		ctx.SetSourceRGB(0.0, 0.0, 0.0)
	}

	ctx.SetOperator(cairo.OPERATOR_SOURCE)
	ctx.Paint()

	app.Ctx = ctx
}
