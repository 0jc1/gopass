package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gopass/internal/auth"
	"gopass/internal/storage"
)

type MainApp struct {
	window     fyne.Window
	auth       *auth.Auth
	storage    *storage.Storage
	authScreen *AuthScreen
	output     *widget.TextGrid
	passwordTab *PasswordTab
	notesTab    *NotesTab
	dataTabs    *DataTabs
}

func NewMainApp(window fyne.Window) *MainApp {
	app := &MainApp{
		window: window,
		auth:   auth.NewAuth(),
		output: widget.NewTextGrid(),
	}

	app.authScreen = NewAuthScreen(window, app.auth, app.onAuthSuccess)
	app.passwordTab = NewPasswordTab(window, app)
	app.notesTab = NewNotesTab(window, app)
	app.dataTabs = NewDataTabs(window, app)
	return app
}

func (m *MainApp) LoadAuth() {
	err := m.auth.LoadPINHash()
	if err != nil && err.Error() != "PIN not set" {
		m.logOutput("Error loading PIN: " + err.Error())
	}
	m.authScreen.Load()
}

func (m *MainApp) onAuthSuccess() {
	// Initialize storage with PIN
	m.storage = storage.NewStorage(m.auth.GetCurrentPIN())
	if err := m.storage.Load(); err != nil {
		m.logOutput("Error loading data: " + err.Error())
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Passwords", m.createPasswordsTab()),
		container.NewTabItem("Notes", m.createNotesTab()),
		container.NewTabItem("Export Data", m.createExportTab()),
		container.NewTabItem("Import Data", m.createImportTab()),
	)

	content := container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewLabel("System Output:"),
			m.output,
		),
		nil,
		nil,
		tabs,
	)

	m.window.SetContent(content)
	m.logOutput("Successfully authenticated")
}

func (m *MainApp) createPasswordsTab() fyne.CanvasObject {
	return m.passwordTab.createContent()
}

func (m *MainApp) createNotesTab() fyne.CanvasObject {
	return m.notesTab.createContent()
}

func (m *MainApp) createExportTab() fyne.CanvasObject {
	return m.dataTabs.createExportTab()
}

func (m *MainApp) createImportTab() fyne.CanvasObject {
	return m.dataTabs.createImportTab()
}

func (m *MainApp) logOutput(message string) {
	m.output.SetText(m.output.Text() + "\n" + message)
}
