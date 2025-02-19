package gui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type DataTabs struct {
	window  fyne.Window
	mainApp *MainApp
}

func NewDataTabs(window fyne.Window, mainApp *MainApp) *DataTabs {
	return &DataTabs{
		window:  window,
		mainApp: mainApp,
	}
}

func (d *DataTabs) createExportTab() fyne.CanvasObject {
	description := widget.NewTextGrid()
	description.SetText("Export your encrypted data to a file. You can use this file to backup your data or transfer it to another device.")

	exportBtn := widget.NewButton("Export Data", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, d.window)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()

			// Handle export in goroutine
			go func() {
				data, err := d.mainApp.storage.Export()
				if err != nil {
					d.window.Canvas().Refresh(d.window.Content())
					dialog.ShowError(err, d.window)
					return
				}

				_, err = writer.Write(data)
				if err != nil {
					d.window.Canvas().Refresh(d.window.Content())
					dialog.ShowError(err, d.window)
					return
				}

				d.window.Canvas().Refresh(d.window.Content())
				d.mainApp.logOutput(fmt.Sprintf("Data exported successfully to %s", writer.URI().Path()))
			}()
		}, d.window)
	})

	return container.NewVBox(
		description,
		exportBtn,
	)
}

func (d *DataTabs) createImportTab() fyne.CanvasObject {
	description := widget.NewTextGrid()
	description.SetText("Import encrypted data from a file. This will add the imported passwords and notes to your existing data.")

	importBtn := widget.NewButton("Import Data", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, d.window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			// Handle import in goroutine
			go func() {
				data, err := os.ReadFile(reader.URI().Path())
				if err != nil {
					d.window.Canvas().Refresh(d.window.Content())
					dialog.ShowError(err, d.window)
					return
				}

				err = d.mainApp.storage.Import(data)
				if err != nil {
					d.window.Canvas().Refresh(d.window.Content())
					dialog.ShowError(err, d.window)
					return
				}

				d.window.Canvas().Refresh(d.window.Content())
				d.mainApp.logOutput(fmt.Sprintf("Data imported successfully from %s", filepath.Base(reader.URI().Path())))
			}()
		}, d.window)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fd.Show()
	})

	return container.NewVBox(
		description,
		importBtn,
	)
}
