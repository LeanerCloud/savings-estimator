package screens

import (
	"fmt"

	"github.com/LeanerCloud/savings-estimator/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	preferenceProfile = "profile"
	preferenceRegion  = "region"

	preferenceAutoSpottingVersion = "AutoSpottingVersion"
)

func staticAuth(a fyne.App, c *core.Launcher) *widget.AccordionItem {
	currentPrefRegion := a.Preferences().StringWithFallback(preferenceRegion, "us-east-1")

	accessKey := widget.NewEntry()
	accessKey.SetPlaceHolder("Usually starts with AKIA...")

	secret := widget.NewEntry()
	secret.SetPlaceHolder("Longer string")

	sessionToken := widget.NewEntry()
	sessionToken.SetPlaceHolder("(optional)")

	regionsStaticAuth := widget.NewSelect(c.AWSRegions(), func(s string) {
		a.Preferences().SetString(preferenceRegion, s)
		fmt.Println("selected region", s)
		c.ConnectWithStaticAuth(accessKey.Text, secret.Text, sessionToken.Text)
		c.SetRegion(s)
	})

	regionsStaticAuth.SetSelected(currentPrefRegion)
	regionsStaticAuth.PlaceHolder = "Select region"

	return widget.NewAccordionItem("Static AWS Credentials", container.NewVBox(
		widget.NewLabel("(Not persisted in the configuration)"),
		&widget.Form{
			Items: []*widget.FormItem{
				{Text: "Access Key ID", Widget: accessKey, HintText: ""},
				{Text: "Secret Access Key", Widget: secret, HintText: ""},
				{Text: "Session Token", Widget: sessionToken, HintText: ""},
				{Text: "Region", Widget: regionsStaticAuth, HintText: ""},
			}}))

}

func profileAuth(a fyne.App, c *core.Launcher) *widget.AccordionItem {
	currentPrefRegion := a.Preferences().StringWithFallback(preferenceRegion, "us-east-1")

	currentPrefProfile := a.Preferences().String(preferenceProfile)
	regionsProfileAuth := widget.NewSelect(c.AWSRegions(), func(s string) {})
	profiles := widget.NewSelect(c.ReadAWSProfiles(), func(s string) {
		a.Preferences().SetString(preferenceProfile, s)
		fmt.Println("selected profile", s)
		c.ConnectWithProfileAuth(regionsProfileAuth.Selected)
		c.SetRegion(s)
	})
	profiles.SetSelected(currentPrefProfile)

	regionsProfileAuth.SetSelected(currentPrefRegion)

	regionsProfileAuth.PlaceHolder = "Type or select region"

	regionsProfileAuth.OnChanged = func(s string) {
		a.Preferences().SetString(preferenceRegion, s)
		fmt.Println("selected region", s)
		c.ConnectWithProfileAuth(profiles.Selected)
		c.SetRegion(s)
	}
	return widget.NewAccordionItem("AWS Profile", container.NewVBox(
		widget.NewLabel("(Persisted in the configuration)"),
		&widget.Form{
			Items: []*widget.FormItem{
				{Text: "Region", Widget: regionsProfileAuth, HintText: ""},
				{Text: "Profile Name", Widget: profiles, HintText: ""},
			}}))
}

func authentication(a fyne.App, c *core.Launcher) *container.TabItem {

	acc := widget.NewAccordion(staticAuth(a, c), profileAuth(a, c))
	acc.MultiOpen = true
	return container.NewTabItem("Authentication", acc)
}

func autoSpottingConfiguration(a fyne.App, c *core.Launcher) *container.TabItem {
	autoSpottingVersions := []string{
		"stable-promo-1.1.1-0",
		"stable-1.1.1-0",
	}

	versions := widget.NewSelect(autoSpottingVersions, func(s string) {
		a.Preferences().SetString(preferenceAutoSpottingVersion, s)
		fmt.Println("selected AutoSpotting Version", s)
		//	c.ProcessAutoSpottingTemplate(a, s)
	})

	// currentPrefRegion := a.Preferences().StringWithFallback(preferenceRegion, "us-east-1")

	// accessKey := widget.NewEntry()
	// accessKey.SetPlaceHolder("Usually starts with AKIA...")

	// secret := widget.NewEntry()
	// secret.SetPlaceHolder("Longer string")

	// sessionToken := widget.NewEntry()
	// sessionToken.SetPlaceHolder("(optional)")

	// regionsStaticAuth := widget.NewSelect(awsRegions(), func(s string) {
	// 	a.Preferences().SetString(preferenceRegion, s)
	// 	fmt.Println("selected region", s)
	// 	c.StaticAuth(accessKey.Text, secret.Text, sessionToken.Text, s)
	// })

	// regionsStaticAuth.SetSelected(currentPrefRegion)
	// regionsStaticAuth.PlaceHolder = "Select region"

	// return widget.NewAccordionItem("Static AWS Credentials", container.NewVBox(
	// 	widget.NewLabel("(Not persisted in the configuration)"),
	// 	&widget.Form{
	// 		Items: []*widget.FormItem{
	// 			{Text: "Access Key ID", Widget: accessKey, HintText: ""},
	// 			{Text: "Secret Access Key", Widget: secret, HintText: ""},
	// 			{Text: "Session Token", Widget: sessionToken, HintText: ""},
	// 			{Text: "Region", Widget: regionsStaticAuth, HintText: ""},
	// 		}}))

	return container.NewTabItem("AutoSpotting", &widget.Form{
		Items: []*widget.FormItem{
			{Text: "AutoSpotting Versions", Widget: versions, HintText: ""},
			// {Text: "Secret Access Key", Widget: secret, HintText: ""},
			// {Text: "Session Token", Widget: sessionToken, HintText: ""},
			// {Text: "Region", Widget: regionsStaticAuth, HintText: ""},
		}})
}

func ebsOptimizerConfiguration(a fyne.App, c *core.Launcher) *container.TabItem {
	// currentPrefRegion := a.Preferences().StringWithFallback(preferenceRegion, "us-east-1")

	// accessKey := widget.NewEntry()
	// accessKey.SetPlaceHolder("Usually starts with AKIA...")

	// secret := widget.NewEntry()
	// secret.SetPlaceHolder("Longer string")

	// sessionToken := widget.NewEntry()
	// sessionToken.SetPlaceHolder("(optional)")

	// regionsStaticAuth := widget.NewSelect(awsRegions(), func(s string) {
	// 	a.Preferences().SetString(preferenceRegion, s)
	// 	fmt.Println("selected region", s)
	// 	c.StaticAuth(accessKey.Text, secret.Text, sessionToken.Text, s)
	// })

	// regionsStaticAuth.SetSelected(currentPrefRegion)
	// regionsStaticAuth.PlaceHolder = "Select region"

	// return widget.NewAccordionItem("Static AWS Credentials", container.NewVBox(
	// 	widget.NewLabel("(Not persisted in the configuration)"),
	// 	&widget.Form{
	// 		Items: []*widget.FormItem{
	// 			{Text: "Access Key ID", Widget: accessKey, HintText: ""},
	// 			{Text: "Secret Access Key", Widget: secret, HintText: ""},
	// 			{Text: "Session Token", Widget: sessionToken, HintText: ""},
	// 			{Text: "Region", Widget: regionsStaticAuth, HintText: ""},
	// 		}}))
	return container.NewTabItem("EBS Optimizer", widget.NewLabel("ToDo"))
}

func configuration(_ fyne.Window, c *core.Launcher) fyne.CanvasObject {

	a := fyne.CurrentApp()

	return container.NewAppTabs(
		authentication(a, c),
		// autoSpottingConfiguration(a, c),
		// ebsOptimizerConfiguration(a, c),
	)

}

// selectEntry := widget.NewSelectEntry([]string{"Option A", "Option B", "Option C"})
// selectEntry.PlaceHolder = "Type or select"
// disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
// disabledCheck.Disable()
// checkGroup := widget.NewCheckGroup([]string{"CheckGroup Item 1", "CheckGroup Item 2"}, func(s []string) { fmt.Println("selected", s) })
// checkGroup.Horizontal = true
// radio := widget.NewRadioGroup([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { fmt.Println("selected", s) })
// radio.Horizontal = true
// disabledRadio := widget.NewRadioGroup([]string{"Disabled radio"}, func(string) {})
// disabledRadio.Disable()

// return container.NewVBox(
// 	widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) { fmt.Println("selected", s) }),
// 	selectEntry,
// 	widget.NewCheck("Check", func(on bool) { fmt.Println("checked", on) }),
// 	disabledCheck,
// 	checkGroup,
// 	radio,
// 	disabledRadio,
// 	widget.NewSlider(0, 100),
// )
