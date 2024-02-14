// Packascreenshot provides various examples of Fyne API capabilities.
package main

import (
	"log"

	"github.com/LeanerCloud/savings-estimator/core"
	"github.com/LeanerCloud/savings-estimator/screens"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
)

const (
	preferenceCurrentPage = "currentPage"
)

var topWindow fyne.Window

var c *core.Launcher

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	c = &core.Launcher{
		PricingIntervalMultiplier: 1,
	}

	c.AutoSpottingCurrentTotalMonthlyCosts = binding.NewString()
	c.AutoSpottingProjectedMonthlyCosts = binding.NewString()
	c.AutoSpottingProjectedSpotSavings = binding.NewString()
	c.AutoSpottingProjectedSpotSavingsPercent = binding.NewString()
	c.AutoSpottingProjectedAutoSpottingCharges = binding.NewString()
	c.AutoSpottingProjectedNetSavings = binding.NewString()

	c.AutoSpottingCurrentTotalMonthlyCosts.Set("0")
	c.AutoSpottingProjectedMonthlyCosts.Set("0")
	c.AutoSpottingProjectedSpotSavings.Set("0")
	c.AutoSpottingProjectedSpotSavingsPercent.Set("0%")
	c.AutoSpottingProjectedAutoSpottingCharges.Set("0")
	c.AutoSpottingProjectedNetSavings.Set("0")

	data, err := ec2instancesinfo.Data()
	if err != nil {
		log.Fatalln("Couldn't load instance type data")
	}

	c.InstanceTypeData = data

	a := app.NewWithID("com.leanercloud")
	a.SetIcon(theme.FyneLogo())
	//makeTray(a)
	//logLifecycle(a)
	w := a.NewWindow("LeanerCloud")
	topWindow = w

	//w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setScreen := func(t screens.Screen) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow, c))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w, c)}
		content.Refresh()
	}

	screen := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setScreen, false))
	} else {
		split := container.NewHSplit(makeNav(setScreen, true), screen)
		split.Offset = 0.1
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(1500, 768))
	w.ShowAndRun()
}

// func logLifecycle(a fyne.App) {
// 	a.Lifecycle().SetOnStarted(func() {
// 		log.Println("Lifecycle: Started")
// 	})
// 	a.Lifecycle().SetOnStopped(func() {
// 		log.Println("Lifecycle: Stopped")
// 	})
// 	a.Lifecycle().SetOnEnteredForeground(func() {
// 		log.Println("Lifecycle: Entered Foreground")
// 	})
// 	a.Lifecycle().SetOnExitedForeground(func() {
// 		log.Println("Lifecycle: Exited Foreground")
// 	})
// }

// func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
// 	newItem := fyne.NewMenuItem("New", nil)
// 	checkedItem := fyne.NewMenuItem("Checked", nil)
// 	checkedItem.Checked = true
// 	disabledItem := fyne.NewMenuItem("Disabled", nil)
// 	disabledItem.Disabled = true
// 	otherItem := fyne.NewMenuItem("Other", nil)
// 	mailItem := fyne.NewMenuItem("Mail", func() { log.Println("Menu New->Other->Mail") })
// 	mailItem.Icon = theme.MailComposeIcon()
// 	otherItem.ChildMenu = fyne.NewMenu("",
// 		fyne.NewMenuItem("Project", func() { log.Println("Menu New->Other->Project") }),
// 		mailItem,
// 	)
// 	fileItem := fyne.NewMenuItem("File", func() { log.Println("Menu New->File") })
// 	fileItem.Icon = theme.FileIcon()
// 	dirItem := fyne.NewMenuItem("Directory", func() { log.Println("Menu New->Directory") })
// 	dirItem.Icon = theme.FolderIcon()
// 	newItem.ChildMenu = fyne.NewMenu("",
// 		fileItem,
// 		dirItem,
// 		otherItem,
// 	)

// 	openSettings := func() {
// 		w := a.NewWindow("Fyne Settings")
// 		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
// 		w.Resize(fyne.NewSize(480, 480))
// 		w.Show()
// 	}
// 	settingsItem := fyne.NewMenuItem("Settings", openSettings)
// 	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
// 	settingsItem.Shortcut = settingsShortcut
// 	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
// 		openSettings()
// 	})

// 	cutShortcut := &fyne.ShortcutCut{Clipboard: w.Clipboard()}
// 	cutItem := fyne.NewMenuItem("Cut", func() {
// 		shortcutFocused(cutShortcut, w)
// 	})
// 	cutItem.Shortcut = cutShortcut
// 	copyShortcut := &fyne.ShortcutCopy{Clipboard: w.Clipboard()}
// 	copyItem := fyne.NewMenuItem("Copy", func() {
// 		shortcutFocused(copyShortcut, w)
// 	})
// 	copyItem.Shortcut = copyShortcut
// 	pasteShortcut := &fyne.ShortcutPaste{Clipboard: w.Clipboard()}
// 	pasteItem := fyne.NewMenuItem("Paste", func() {
// 		shortcutFocused(pasteShortcut, w)
// 	})
// 	pasteItem.Shortcut = pasteShortcut
// 	performFind := func() { log.Println("Menu Find") }
// 	findItem := fyne.NewMenuItem("Find", performFind)
// 	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
// 	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
// 		performFind()
// 	})

// 	helpMenu := fyne.NewMenu("Help",
// 		fyne.NewMenuItem("Documentation", func() {
// 			u, _ := url.Parse("https://developer.fyne.io")
// 			_ = a.OpenURL(u)
// 		}),
// 		fyne.NewMenuItem("Support", func() {
// 			u, _ := url.Parse("https://fyne.io/support/")
// 			_ = a.OpenURL(u)
// 		}),
// 		fyne.NewMenuItemSeparator(),
// 		fyne.NewMenuItem("Sponsor", func() {
// 			u, _ := url.Parse("https://fyne.io/sponsor/")
// 			_ = a.OpenURL(u)
// 		}))

// 	// a quit item will be appended to our first (File) menu
// 	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem)
// 	device := fyne.CurrentDevice()
// 	if !device.IsMobile() && !device.IsBrowser() {
// 		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
// 	}
// 	main := fyne.NewMainMenu(
// 		file,
// 		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
// 		helpMenu,
// 	)
// 	checkedItem.Action = func() {
// 		checkedItem.Checked = !checkedItem.Checked
// 		main.Refresh()
// 	}
// 	return main
// }

// func makeTray(a fyne.App) {
// 	if desk, ok := a.(desktop.App); ok {
// 		h := fyne.NewMenuItem("Hello", func() {})
// 		menu := fyne.NewMenu("Hello World", h)
// 		h.Action = func() {
// 			log.Println("System tray menu tapped")
// 			h.Label = "Welcome"
// 			menu.Refresh()
// 		}
// 		desk.SetSystemTrayMenu(menu)
// 	}
// }

func unsupportedTutorial(t screens.Screen) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func makeNav(setTutorial func(tutorial screens.Screen), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return screens.ScreenIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := screens.ScreenIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := screens.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := screens.Tutorials[uid]; ok {
				t.Core = c
				if unsupportedTutorial(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentPage, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentPage, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

// func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
// 	switch sh := s.(type) {
// 	case *fyne.ShortcutCopy:
// 		sh.Clipboard = w.Clipboard()
// 	case *fyne.ShortcutCut:
// 		sh.Clipboard = w.Clipboard()
// 	case *fyne.ShortcutPaste:
// 		sh.Clipboard = w.Clipboard()
// 	}
// 	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
// 		focused.TypedShortcut(s)
// 	}
// }
