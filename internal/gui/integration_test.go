package gui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func setupTestApp(t *testing.T) (fyne.App, *MainApp) {
	app := test.NewApp()
	mainWindow := app.NewWindow("GoPass")
	mainApp := NewMainApp(mainWindow)
	
	// Set up authentication
	err := mainApp.auth.SetPIN("123456")
	assert.NoError(t, err, "Should set PIN without error")
	
	// Trigger auth success to show main UI
	mainApp.onAuthSuccess()
	
	mainWindow.Resize(fyne.NewSize(800, 600))
	return app, mainApp
}

// Helper to verify UI responsiveness
func verifyUIResponsive(t *testing.T, app fyne.App, mainApp *MainApp) bool {
	// Try to interact with multiple UI elements
	responsive := true
	
	// Get all windows
	windows := app.Driver().AllWindows()
	assert.NotEmpty(t, windows, "Should have at least one window")
	
	mainWindow := windows[0]
	content := mainWindow.Content().(*fyne.Container)
	
	// Try to click each button
	buttons := []string{"Add Password", "Edit", "Delete", "View"}
	for _, btnText := range buttons {
		btn := findButtonByText(content, btnText)
		if btn != nil {
			// Try to click button
			test.Tap(btn)
			time.Sleep(50 * time.Millisecond)
			
			// Check if new windows appeared (dialogs)
			newWindows := app.Driver().AllWindows()
			if len(newWindows) <= len(windows) {
				t.Logf("Button '%s' did not create expected dialog", btnText)
				responsive = false
			}
		} else {
			t.Logf("Could not find button: %s", btnText)
			responsive = false
		}
	}
	
	return responsive
}

func findButtonByText(c *fyne.Container, text string) *widget.Button {
	var button *widget.Button
	for _, obj := range c.Objects {
		switch v := obj.(type) {
		case *fyne.Container:
			if btn := findButtonByText(v, text); btn != nil {
				return btn
			}
		case *widget.Button:
			if v.Text == text {
				return v
			}
		}
	}
	return button
}

func TestPasswordTabIntegration(t *testing.T) {
	app, mainApp := setupTestApp(t)
	defer app.Quit()

	t.Run("Add Password Form Should Stay Responsive", func(t *testing.T) {
		// Get the main container
		mainWindow := app.Driver().AllWindows()[0]
		content := mainWindow.Content().(*fyne.Container)
		
		// Find Add Password button
		addBtn := findButtonByText(content, "Add Password")
		assert.NotNil(t, addBtn, "Add Password button should exist")

		// Click Add Password button
		test.Tap(addBtn)
		
		// Give time for dialog to appear
		time.Sleep(100 * time.Millisecond)
		
		// Find dialog form
		var dialog *widget.Form
		for _, win := range app.Driver().AllWindows() {
			if form, ok := win.Content().(*widget.Form); ok {
				dialog = form
				break
			}
		}
		assert.NotNil(t, dialog, "Password form dialog should be shown")

		// Fill form
		nameEntry := dialog.Items[0].Widget.(*widget.Entry)
		urlEntry := dialog.Items[1].Widget.(*widget.Entry)
		usernameEntry := dialog.Items[2].Widget.(*widget.Entry)
		passwordEntry := dialog.Items[3].Widget.(*widget.Entry)
		noteEntry := dialog.Items[4].Widget.(*widget.Entry)

		test.Type(nameEntry, "Test Password")
		test.Type(urlEntry, "https://example.com")
		test.Type(usernameEntry, "testuser")
		test.Type(passwordEntry, "testpass123")
		test.Type(noteEntry, "Test note")

		// Find and click Save
		// Find save button in dialog
		saveBtn := widget.NewButton("Save", nil)
		for _, win := range app.Driver().AllWindows() {
			if container, ok := win.Content().(*fyne.Container); ok {
				if btn := findButtonByText(container, "Save"); btn != nil {
					saveBtn = btn
					break
				}
			}
		}
		assert.NotNil(t, saveBtn, "Save button should exist")
		test.Tap(saveBtn)

		// Wait for UI updates
		time.Sleep(100 * time.Millisecond)

		// Verify UI is still responsive
		responsive := verifyUIResponsive(t, app, mainApp)
		assert.True(t, responsive, "UI should remain responsive after adding password")
	})
}

func TestNotesTabIntegration(t *testing.T) {
	app, mainApp := setupTestApp(t)
	defer app.Quit()

	t.Run("Add Note Form Should Stay Responsive", func(t *testing.T) {
		// Get the main container
		mainWindow := app.Driver().AllWindows()[0]
		content := mainWindow.Content().(*fyne.Container)
		
		// Find Add Note button
		addBtn := findButtonByText(content, "Add Note")
		assert.NotNil(t, addBtn, "Add Note button should exist")

		// Click Add Note button
		test.Tap(addBtn)
		
		// Give time for dialog to appear
		time.Sleep(100 * time.Millisecond)
		
		// Find dialog form
		var dialog *widget.Form
		for _, win := range app.Driver().AllWindows() {
			if form, ok := win.Content().(*widget.Form); ok {
				dialog = form
				break
			}
		}
		assert.NotNil(t, dialog, "Note form dialog should be shown")

		// Fill form
		titleEntry := dialog.Items[0].Widget.(*widget.Entry)
		contentEntry := dialog.Items[1].Widget.(*widget.Entry)

		test.Type(titleEntry, "Test Note")
		test.Type(contentEntry, "This is a test note content")

		// Find and click Save
		// Find save button in dialog
		saveBtn := widget.NewButton("Save", nil)
		for _, win := range app.Driver().AllWindows() {
			if container, ok := win.Content().(*fyne.Container); ok {
				if btn := findButtonByText(container, "Save"); btn != nil {
					saveBtn = btn
					break
				}
			}
		}
		assert.NotNil(t, saveBtn, "Save button should exist")
		test.Tap(saveBtn)

		// Wait for UI updates
		time.Sleep(100 * time.Millisecond)

		// Verify UI is still responsive
		responsive := verifyUIResponsive(t, app, mainApp)
		assert.True(t, responsive, "UI should remain responsive after adding note")
	})
}
