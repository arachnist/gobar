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

	ic := icon()
	pb := progressBar()
	pb.SetSizeRequest(700, 3)

	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.Add(ic)
	grid.Add(pb)
	win.Add(grid)

	go func() {
		time.Sleep(time.Second * 3)
		gtk.MainQuit()
	}()

	win.ShowAll()
	gtk.Main()
}

func icon() *gtk.Widget {
	image, err := gtk.ImageNewFromFile(os.Args[1])

	if err != nil {
		log.Fatal("Unable to create image:", err)
	}

	return &image.Widget
}

func progressBar() *gtk.Widget {
	pb, err := gtk.ProgressBarNew()

	if err != nil {
		log.Fatal("Unable to create progress bar:", err)
	}

	f, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatal("Can't convert float:", err)
	}

	pb.SetFraction(f)

	return &pb.Widget
}
