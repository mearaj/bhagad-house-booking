package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"github.com/mearaj/bhagad-house-booking/frontend/ui"
	"log"
	"os"
)

func main() {
	// fmt.Print("API_URL IS:  ")
	// fmt.Println(frontend.LoadConfig().ApiURL)

	go func() {
		title := app.Title("Bhagad House Booking")
		w := app.NewWindow(title, app.Size(1024, 768))
		w.Perform(system.ActionCenter)
		if err := ui.Loop(w); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
