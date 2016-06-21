package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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

func main() {
	var config gobarConfig

	if len(os.Args) < 4 {
		log.Fatalln("Usage:", os.Args[0], "<configuration file> <level> <label>")
	}

	f, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatal("Can't convert float:", err)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Error reading configuration file:", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalln("Error parsing configuration file:", err)
	}

	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		log.Fatal("Unable to parse duration:", err)
	}

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

	lb := label(os.Args[3])

	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	grid.Add(levelBar(f, config.BarSize.X, config.BarSize.Y))
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
	win.Move(config.Position.X, config.Position.Y)

	go func() {
		time.Sleep(duration)
		gtk.MainQuit()
	}()

	win.ShowAll()
	gtk.Main()
}

func label(text string) *gtk.Widget {
	label, err := gtk.LabelNew(text)

	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	return &label.Widget
}

func levelBar(value float64, x, y int) *gtk.Widget {
	lb, err := gtk.LevelBarNew()

	if err != nil {
		log.Fatal("Unable to create level bar:", err)
	}

	lb.SetValue(value)
	lb.SetSizeRequest(x, y)

	return &lb.Widget
}
