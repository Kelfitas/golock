package main

import (
	"log"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const AppName = "GoLock"

type App struct {
	Window         *gtk.Window
	Label          *gtk.Label
	CapsLabel      *gtk.Label
	AlphaSupported bool
}

var app *App

func setMessage(str string, isTemp bool) {
	app.Label.SetLabel(str)
	if !isTemp {
		return
	}

	go func() {
		select {
		case event := <-state.EventChan:
			state.EventChan <- event
		case <-time.After(5 * time.Second):
			app.Label.SetLabel(" ")
		}
	}()
}

func watchState(done chan bool) {
	for {
		select {
		case event := <-state.EventChan:
			log.Println("New Event!")
			switch event {
			case EmptyPasswordEvent:
				setMessage("Empty", true)
				break
			case CapsChangedEvent:
				if state.IsCapsLockOn {
					app.CapsLabel.SetText("CapsLock on")
				} else {
					app.CapsLabel.SetText(" ")
				}
				break
			case WrongPasswordEvent:
				setMessage("Wrong password", true)
				break
			case KeyPressEvent:
				break
			}
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
	win.SetTitle(AppName)
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
	})

	l, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	capsLabel, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	win.SetDefaultSize(1920, 1080)
	win.SetProperty("opacity", 0.01)
	app = &App{
		Window:    win,
		Label:     l,
		CapsLabel: capsLabel,
	}
	// win.Fullscreen()

	// textView, err := gtk.TextViewNew()
	// if err != nil {
	// 	log.Fatal("Unable to create textView:", err)
	// }

	// textView.Add(l)
	// textView.Add(capsLabel)
	// win.Add(textView)

	win.Add(l)
	// win.Add(capsLabel)

	screenChanged(&win.Widget)

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

	widget.SetVisual(visual)
}

func exposeDraw(w *gtk.Window, ctx *cairo.Context) {
	log.Println("=========")
	log.Printf("AlphaSupported: %t\n", app.AlphaSupported)
	log.Println("=========")
	if app.AlphaSupported {
		// ctx.SetSourceRGBA(0.0, 0.0, 0.0, 0.25)
		ctx.SetSourceRGBA(0.0, 0.0, 1.0, 0.01)
	} else {
		ctx.SetSourceRGB(0.0, 0.0, 0.0)
	}

	ctx.SetOperator(cairo.OPERATOR_SOURCE)
	ctx.Paint()
}
