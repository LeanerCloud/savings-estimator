package screens

import (
	"net/url"

	"github.com/LeanerCloud/savings-estimator/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	logo := canvas.NewImageFromFile("logo.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	return container.NewCenter(container.NewVBox(widget.NewLabelWithStyle(
		`Welcome to the LeanerCloud AWS Spot cost savings estimation tool!

This tool allows you to estimate the cost savings you can achieve in your AWS account by converting your AutoScaling Groups to Spot instances.

You can select various scenarios, such as to keep some of your instances as OnDemand in each group (maybe covered by Reserved Instances or Savings Plans),
or only convert some of your AutoScaling Groups to Spot as part of a gradual rollout.

You may use any mechanism to adopt Spot, such as converting the configuration yourself group by group as per what you defined in this tool.


For your convenience, you can also use AutoSpotting, our state of the art cost optimization engine for Spot. AutoSpotting is tightly integrated with
this cost estimator, so you can apply this configuration with a single click, by tagging them as expected by AutoSpotting. You just need to install 
AutoSpotting from the AWS Marketplace link below

AutoSpotting allows you to adopt Spot instances with all the Spot best practices, recommended by AWS such as wide diversification over multiple instance
types and uses capacity optimized allocation strategy optimized for reduced interruptions.


In addition, AutoSpotting also prioritizes for lower cost instances from newer generations and implements a reliable failover to on-demand instances when running out of Spot capacity.

In most situations AutoSpotting doesn't require any configuration changes to your AutoScaling Groups, but uses the existing launch template or launch configuration.`,
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(container.NewHBox(
			widget.NewHyperlink("LeanerCloud GUI on GitHub", parseURL("https://github.com/LeanerCloud/savings-estimator")),
			widget.NewLabel("-"),
			widget.NewHyperlink("LeanerCloud.com", parseURL("https://leanercloud.com/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("AutoSpotting.io", parseURL("https://autospotting.io")),
			widget.NewLabel("-"),
			widget.NewHyperlink("Install AutoSpotting from the AWS Marketplace", parseURL("https://aws.amazon.com/marketplace/pp/prodview-6uj4pruhgmun6")),
		)),
		logo,
	))
}
