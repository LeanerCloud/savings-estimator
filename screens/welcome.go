package screens

import (
	"leanercloud/core"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func welcomeScreen(_ fyne.Window, _ *core.Launcher) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to the LeanerCloud configuration tool. Currently supported tools: AutoSpotting and EBS Optimizer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		//logo,
		container.NewHBox(
			widget.NewHyperlink("leanercloud.com", parseURL("https://leanercloud.com/")),
			// widget.NewHyperlink("leanercloud.com", parseURL("https://leanercloud.com/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("autospotting.io", parseURL("https://autospotting.io")),
			// widget.NewLabel("-"),

		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}
