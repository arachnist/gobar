package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}

	css, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Unable to create css provider:", err)
	}

	lb := label()

	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	grid.Add(lb)
	grid.Add(levelBar())

	css.LoadFromPath("gobar.css")

	style_context, err := lb.GetStyleContext()
	if err != nil {
		log.Fatal("Unable to get label style context:", err)
	}
	style_context.AddProvider(css, 0)

	win.Add(grid)
	win.Move(333, 700)

	go func() {
		time.Sleep(time.Second * 1)
		gtk.MainQuit()
		//		win.Hide()
	}()

	win.ShowAll()
	gtk.Main()
}

func label() *gtk.Widget {
	label, err := gtk.LabelNew(os.Args[1])

	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	return &label.Widget
}

func levelBar() *gtk.Widget {
	lb, err := gtk.LevelBarNew()

	if err != nil {
		log.Fatal("Unable to create level bar:", err)
	}

	f, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatal("Can't convert float:", err)
	}

	lb.SetValue(f)
	lb.SetSizeRequest(700, 30)

	return &lb.Widget
}
