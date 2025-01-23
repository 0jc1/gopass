package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gopass/internal/gui"
)

func main() {
	a := app.New()
	w := a.NewWindow("GoPass - Password Manager")
	w.Resize(fyne.NewSize(800, 600))

	mainApp := gui.NewMainApp(w)
	mainApp.LoadAuth()

	w.ShowAndRun()
}
