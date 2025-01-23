package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gopass/internal/auth"
)

type AuthScreen struct {
	window fyne.Window
	auth   *auth.Auth
	onAuth func()
}

func NewAuthScreen(window fyne.Window, auth *auth.Auth, onAuth func()) *AuthScreen {
	return &AuthScreen{
		window: window,
		auth:   auth,
		onAuth: onAuth,
	}
}

func (a *AuthScreen) Load() {
	var content fyne.CanvasObject

	if !a.auth.IsPINSet() {
		content = a.createSetPINScreen()
	} else {
		content = a.createLoginScreen()
	}

	a.window.SetContent(content)
}

func (a *AuthScreen) createSetPINScreen() fyne.CanvasObject {
	pinEntry := widget.NewPasswordEntry()
	confirmEntry := widget.NewPasswordEntry()
	message := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Create PIN", Widget: pinEntry},
			{Text: "Confirm PIN", Widget: confirmEntry},
		},
		OnSubmit: func() {
			if pinEntry.Text != confirmEntry.Text {
				message.SetText("PINs do not match")
				return
			}
			if len(pinEntry.Text) < 4 {
				message.SetText("PIN must be at least 4 characters")
				return
			}

			err := a.auth.SetPIN(pinEntry.Text)
			if err != nil {
				message.SetText("Error setting PIN: " + err.Error())
				return
			}

			a.onAuth()
		},
	}

	return container.NewVBox(
		widget.NewLabel("Welcome to GoPass"),
		widget.NewLabel("Please set a PIN to secure your data"),
		form,
		message,
	)
}

func (a *AuthScreen) createLoginScreen() fyne.CanvasObject {
	pinEntry := widget.NewPasswordEntry()
	message := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Enter PIN", Widget: pinEntry},
		},
		OnSubmit: func() {
			if a.auth.ValidatePIN(pinEntry.Text) {
				a.onAuth()
			} else {
				message.SetText("Invalid PIN")
				pinEntry.SetText("")
			}
		},
	}

	return container.NewVBox(
		widget.NewLabel("Welcome back to GoPass"),
		form,
		message,
	)
}
