package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"gopass/internal/models"
)

type NotesTab struct {
	window  fyne.Window
	mainApp *MainApp
	table   *widget.Table
	notes   []models.Note
}

func NewNotesTab(window fyne.Window, mainApp *MainApp) *NotesTab {
	return &NotesTab{
		window:  window,
		mainApp: mainApp,
	}
}

func (n *NotesTab) createContent() fyne.CanvasObject {
	// Create search entry
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search notes...")
	searchEntry.OnChanged = func(text string) {
		results := n.mainApp.storage.Search(text)
		n.notes = results.Notes
		n.table.Refresh()
	}

	// Create table
	n.notes = n.mainApp.storage.GetNotes()
	n.table = widget.NewTable(
		func() (int, int) {
			return len(n.notes), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			note := n.notes[i.Row]
			switch i.Col {
			case 0:
				label.SetText(note.Title)
			case 1:
				// Show preview of content
				if len(note.Content) > 50 {
					label.SetText(note.Content[:47] + "...")
				} else {
					label.SetText(note.Content)
				}
			}
		},
	)

	// Add button
	addBtn := widget.NewButton("Add Note", func() {
		n.showNoteDialog(nil)
	})

	// Edit button
	editBtn := widget.NewButton("Edit", func() {
		if len(n.notes) == 0 {
			return
		}
		selected := n.table.SelectedCell()
		if selected.Row < 0 {
			dialog.ShowInformation("Select Entry", "Please select a note to edit", n.window)
			return
		}
		n.showNoteDialog(&n.notes[selected.Row])
	})

	// Delete button
	deleteBtn := widget.NewButton("Delete", func() {
		if len(n.notes) == 0 {
			return
		}
		selected := n.table.SelectedCell()
		if selected.Row < 0 {
			dialog.ShowInformation("Select Entry", "Please select a note to delete", n.window)
			return
		}
		dialog.ShowConfirm("Delete Note", "Are you sure you want to delete this note?",
			func(ok bool) {
				if ok {
					err := n.mainApp.storage.DeleteNote(n.notes[selected.Row].ID)
					if err != nil {
						dialog.ShowError(err, n.window)
						return
					}
					n.notes = n.mainApp.storage.GetNotes()
					n.table.Refresh()
					n.mainApp.logOutput("Note deleted successfully")
				}
			}, n.window)
	})

	// View button
	viewBtn := widget.NewButton("View", func() {
		if len(n.notes) == 0 {
			return
		}
		selected := n.table.SelectedCell()
		if selected.Row < 0 {
			dialog.ShowInformation("Select Entry", "Please select a note to view", n.window)
			return
		}
		note := n.notes[selected.Row]
		content := widget.NewTextGrid()
		content.SetText(fmt.Sprintf("Title: %s\n\n%s", note.Title, note.Content))
		dialog.ShowCustom("Note Details", "Close", content, n.window)
	})

	buttons := container.NewHBox(addBtn, editBtn, deleteBtn, viewBtn)
	count := widget.NewLabel(fmt.Sprintf("Total Notes: %d", len(n.notes)))

	return container.NewBorder(
		container.NewVBox(searchEntry, buttons, count),
		nil, nil, nil,
		n.table,
	)
}

func (n *NotesTab) showNoteDialog(note *models.Note) {
	isNew := note == nil
	if isNew {
		note = models.NewNote()
		note.ID = uuid.New().String()
	}

	titleEntry := widget.NewEntry()
	contentEntry := widget.NewMultiLineEntry()

	if !isNew {
		titleEntry.SetText(note.Title)
		contentEntry.SetText(note.Content)
	}

	items := []*widget.FormItem{
		{Text: "Title", Widget: titleEntry},
		{Text: "Content", Widget: contentEntry},
	}

	dialog.ShowForm(isNew?"Add Note":"Edit Note", "Save", "Cancel", items,
		func(ok bool) {
			if !ok {
				return
			}

			note.Title = titleEntry.Text
			note.Content = contentEntry.Text
			note.UpdatedAt = time.Now()

			var err error
			if isNew {
				err = n.mainApp.storage.AddNote(*note)
			} else {
				err = n.mainApp.storage.UpdateNote(*note)
			}

			if err != nil {
				dialog.ShowError(err, n.window)
				return
			}

			n.notes = n.mainApp.storage.GetNotes()
			n.table.Refresh()
			n.mainApp.logOutput(fmt.Sprintf("Note %s successfully", isNew?"added":"updated"))
		}, n.window)
}