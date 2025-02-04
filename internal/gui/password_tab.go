package gui

import (
	"fmt"
	"time"

	"gopass/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type PasswordTab struct {
	window      fyne.Window
	mainApp     *MainApp
	table       *widget.Table
	passwords   []models.Password
	selectedRow int
}

func NewPasswordTab(window fyne.Window, mainApp *MainApp) *PasswordTab {
	return &PasswordTab{
		window:      window,
		mainApp:     mainApp,
		selectedRow: -1,
	}
}

func (p *PasswordTab) createContent() fyne.CanvasObject {
	// Create search entry
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search passwords...")
	searchEntry.OnChanged = func(text string) {
		results := p.mainApp.storage.Search(text)
		p.passwords = results.Passwords
		p.table.Refresh()
	}

	// Create table
	p.passwords = p.mainApp.storage.GetPasswords()
	p.table = widget.NewTable(
		func() (int, int) {
			return len(p.passwords), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			password := p.passwords[i.Row]
			switch i.Col {
			case 0:
				label.SetText(password.Name)
			case 1:
				label.SetText(password.URL)
			case 2:
				label.SetText(password.Username)
			case 3:
				label.SetText("********")
			}
		},
	)

	// Set up table selection
	p.table.OnSelected = func(id widget.TableCellID) {
		p.selectedRow = id.Row
	}

	// Add button
	addBtn := widget.NewButton("Add Password", func() {
		p.showPasswordDialog(nil)
	})

	// Edit button
	editBtn := widget.NewButton("Edit", func() {
		if len(p.passwords) == 0 {
			return
		}
		if p.selectedRow < 0 {
			dialog.ShowInformation("Select Entry", "Please select a password entry to edit", p.window)
			return
		}
		p.showPasswordDialog(&p.passwords[p.selectedRow])
	})

	// Delete button
	deleteBtn := widget.NewButton("Delete", func() {
		if len(p.passwords) == 0 {
			return
		}
		if p.selectedRow < 0 {
			dialog.ShowInformation("Select Entry", "Please select a password entry to delete", p.window)
			return
		}
		dialog.ShowConfirm("Delete Password", "Are you sure you want to delete this password?",
			func(ok bool) {
				if ok {
					err := p.mainApp.storage.DeletePassword(p.passwords[p.selectedRow].ID)
					if err != nil {
						dialog.ShowError(err, p.window)
						return
					}
					// Reset selection and refresh table
					p.selectedRow = -1
					p.passwords = p.mainApp.storage.GetPasswords()
					p.table.UnselectAll()
					p.table.Refresh()
					p.mainApp.logOutput("Password deleted successfully")
				}
			}, p.window)
	})

	// View button
	viewBtn := widget.NewButton("View", func() {
		if len(p.passwords) == 0 {
			return
		}
		if p.selectedRow < 0 {
			dialog.ShowInformation("Select Entry", "Please select a password entry to view", p.window)
			return
		}
		pass := p.passwords[p.selectedRow]
		content := widget.NewTextGrid()
		content.SetText(fmt.Sprintf("Name: %s\nURL: %s\nUsername: %s\nPassword: %s\nNote: %s",
			pass.Name, pass.URL, pass.Username, pass.Password, pass.Note))
		dialog.ShowCustom("Password Details", "Close", content, p.window)
	})

	buttons := container.NewHBox(addBtn, editBtn, deleteBtn, viewBtn)
	count := widget.NewLabel(fmt.Sprintf("Total Passwords: %d", len(p.passwords)))

	return container.NewBorder(
		container.NewVBox(searchEntry, buttons, count),
		nil, nil, nil,
		p.table,
	)
}

func (p *PasswordTab) showPasswordDialog(password *models.Password) {
	isNew := password == nil
	if isNew {
		password = models.NewPassword()
		password.ID = uuid.New().String()
	}

	nameEntry := widget.NewEntry()
	urlEntry := widget.NewEntry()
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	noteEntry := widget.NewMultiLineEntry()

	if !isNew {
		nameEntry.SetText(password.Name)
		urlEntry.SetText(password.URL)
		usernameEntry.SetText(password.Username)
		passwordEntry.SetText(password.Password)
		noteEntry.SetText(password.Note)
	}

	items := []*widget.FormItem{
		{Text: "Name", Widget: nameEntry},
		{Text: "URL", Widget: urlEntry},
		{Text: "Username", Widget: usernameEntry},
		{Text: "Password", Widget: passwordEntry},
		{Text: "Note", Widget: noteEntry},
	}

	var formTxt string
	if isNew {
		formTxt = "Add Password"
	} else {
		formTxt = "Edit Password"
	}

	dialog.ShowForm(formTxt, "Save", "Cancel", items,
		func(ok bool) {
			if !ok {
				return
			}

			password.Name = nameEntry.Text
			password.URL = urlEntry.Text
			password.Username = usernameEntry.Text
			password.Password = passwordEntry.Text
			password.Note = noteEntry.Text
			password.UpdatedAt = time.Now()

			var err error
			if isNew {
				err = p.mainApp.storage.AddPassword(*password)
			} else {
				err = p.mainApp.storage.UpdatePassword(*password)
			}

			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}

			// Use goroutine to refresh UI after dialog closes
			go func() {
				// Reset selection and refresh table
				p.selectedRow = -1
				p.passwords = p.mainApp.storage.GetPasswords()
				p.table.UnselectAll()
				p.table.Refresh()

				// Log the operation
				action := "added"
				if !isNew {
					action = "updated"
				}
				p.mainApp.logOutput(fmt.Sprintf("Password %s successfully", action))
			}()
		}, p.window)
}
