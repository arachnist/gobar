package main

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
*/
import "C"

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type gobarConfig struct {
	Duration string `yaml:"duration"`
	CssPath  string `yaml:"css_path"`
	Position struct {
		X int `yaml:"x"`
		Y int `yaml:"y"`
	} `yaml:"position"`
	BarSize struct {
		X int `yaml:"x"`
		Y int `yaml:"y"`
	} `yaml:"bar_size"`
}

func getLabelHandler(label *gtk.Label, bar *gtk.LevelBar, win *gtk.Window) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		min := 0.0
		max := 100.0
		level := 50.0
		duration, _ := time.ParseDuration("700ms")
		labelText := "label string xD"

		vars := r.URL.Query()

		if val, ok := vars["label"]; ok {
			labelText = val[0]
		}
		if val, ok := vars["min"]; ok {
			min, err = strconv.ParseFloat(val[0], 64)
			if err != nil {
				log.Println("Got a wrong min value", val, err)
			}
		}
		if val, ok := vars["max"]; ok {
			max, err = strconv.ParseFloat(val[0], 64)
			if err != nil {
				log.Println("Got a wrong max value", val, err)
			}
		}
		if val, ok := vars["level"]; ok {
			level, err = strconv.ParseFloat(val[0], 64)
			if err != nil {
				log.Println("Got a wrong level value", val, err)
			}
		}
		if val, ok := vars["duration"]; ok {
			duration, err = time.ParseDuration(val[0])
			if err != nil {
				log.Println("Got a wrong duration value", val, err)
			}
		}

		bar.SetMinValue(min)
		bar.SetMaxValue(max)
		bar.SetValue(level)
		label.SetLabel(labelText)

		win.Resize(10, 10)
		win.ShowAll()

		time.Sleep(duration)

		win.Hide()
	}
}

func main() {
	var config gobarConfig

	if len(os.Args) != 2 {
		log.Fatalln("Usage:", os.Args[0], "<configuration file>")
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Error reading configuration file:", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalln("Error parsing configuration file:", err)
	}

	C.XInitThreads()
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}

	css, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Unable to create css provider:", err)
	}

	lb := label()
	bar := levelBar(config.BarSize.X, config.BarSize.Y)

	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	grid.Add(bar)
	grid.Add(lb)

	css.LoadFromPath(config.CssPath)

	style_context, err := lb.GetStyleContext()
	if err != nil {
		log.Fatal("Unable to get label style context:", err)
	}
	style_context.AddProvider(css, gtk.STYLE_PROVIDER_PRIORITY_USER)

	win.SetWMClass("gobar", "gobar")
	win.SetTitle("gobar")
	win.Add(grid)
	win.SetAcceptFocus(false)
	win.Move(config.Position.X, config.Position.Y)

	http.HandleFunc("/bar", getLabelHandler(lb, bar, win))

	go gtk.Main()

	http.ListenAndServe("localhost:8080", nil)
}

func label() *gtk.Label {
	label, err := gtk.LabelNew("xD")

	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	return label
}

func levelBar(x, y int) *gtk.LevelBar {
	lb, err := gtk.LevelBarNew()

	if err != nil {
		log.Fatal("Unable to create level bar:", err)
	}

	lb.SetSizeRequest(x, y)

	return lb
}
