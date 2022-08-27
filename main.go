package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"github.com/mearaj/bhagad-house-booking/ui"
	"log"
	"os"
)

func main() {
	go func() {
		title := app.Title("Bhagad House Booking")
		w := app.NewWindow(title)
		w.Perform(system.ActionCenter)
		if err := ui.Loop(w); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
