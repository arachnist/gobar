package main

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
*/
import "C"

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type gobarConfig struct {
	Listen   string `yaml:"listen"`
	CssPath  string `yaml:"css_path"`
	Position struct {
		X int `yaml:"x"`
		Y int `yaml:"y"`
	} `yaml:"position"`
	BarSize struct {
		X int `yaml:"x"`
		Y int `yaml:"y"`
	} `yaml:"bar_size"`
	Actions map[string]struct {
		Command  string  `yaml:"command"`
		Value    string  `yaml:"value"`
		Label    string  `yaml:"label"`
		Duration string  `yaml:"duration"`
		Min      float64 `yaml:"min"`
		Max      float64 `yaml:"max"`
	} `yaml:"actions"`
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
		r.Body.Close()

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

		_, _ = w.Write([]byte("OK\n"))

		go func() {
			bar.SetMinValue(min)
			bar.SetMaxValue(max)
			bar.SetValue(level)
			label.SetLabel(labelText)

			win.ShowAll()

			time.Sleep(duration)

			win.Resize(10, 10)
			win.Hide()
		}()
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

	http.HandleFunc("/api/v1/bar", getLabelHandler(lb, bar, win))

	for key, value := range config.Actions {
		key := key
		value := value
		http.HandleFunc("/api/v1/action/"+key, func(w http.ResponseWriter, r *http.Request) {
			r.Body.Close()
			err := exec.Command("sh", "-c", value.Command).Run()
			if err != nil {
				log.Println("Command:", value.Command, "error:", err)
				fmt.Fprintln(w, "error:", err)
				return
			}
			out, err := exec.Command("sh", "-c", value.Value).Output()
			if err != nil {
				log.Println("Command:", value.Command, "error:", err)
				fmt.Fprintln(w, "error:", err)
				return
			}
			val, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
			if err != nil {
				log.Println("Error parsing as float", string(out), "error:", err)
				fmt.Fprintln(w, "error:", err)
				return
			}
			duration, err := time.ParseDuration(value.Duration)
			if err != nil {
				log.Println("Error paring as duration", duration, "error:", err)
				fmt.Fprintln(w, "error:", err)
			}

			fmt.Fprintln(w, "status: OK\nvalue:", val)

			go func() {
				lb.SetLabel(value.Label)

				bar.SetMinValue(value.Min)
				bar.SetMaxValue(value.Max)
				bar.SetValue(val)

				win.ShowAll()

				time.Sleep(duration)

				win.Resize(10, 10)
				win.Hide()
			}()
		})
	}

	go gtk.Main()

	http.ListenAndServe(config.Listen, nil)
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
