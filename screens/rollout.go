package screens

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/LeanerCloud/savings-estimator/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type widgetType int64

const (
	preferenceAutoSpottingRolloutRegion   = "AutoSpottingRolloutRegion"
	preferenceAutoSpottingPricingInterval = "AutoSpottingPricingInterval"

	Label widgetType = iota
	Check
	Entry
)

type ColumnInfo struct {
	Header         string
	Type           widgetType
	DataKey        string // New field to identify the data
	EntryValidator fyne.StringValidator
	PlaceHolder    string
}

// ActiveHeader represents a table header that can handle taps.
type ActiveHeader struct {
	widget.Label
	OnTapped func()
}

// newActiveHeader creates a new header cell with the specified label.
func newActiveHeader(label string) *ActiveHeader {
	h := &ActiveHeader{}
	h.ExtendBaseWidget(h)
	h.SetText(label)
	return h
}

// Tapped is called when the header is tapped.
func (h *ActiveHeader) Tapped(_ *fyne.PointEvent) {
	if h.OnTapped != nil {
		h.OnTapped()
	}
}

// TappedSecondary is called on a secondary tap (right-click, long-press, etc.).
func (h *ActiveHeader) TappedSecondary(_ *fyne.PointEvent) {
}

func autoSpottingRollout(a fyne.App, w fyne.Window, c *core.Launcher) *container.TabItem {

	return container.NewTabItem("Convert ASGs to Spot", makeASGTable(w, c))
}

func formatFloat(f float64) string {
	var format string
	if f < 1 {
		format = "%.4f"
	} else {
		format = "%.2f"
	}
	return fmt.Sprintf(format, f)
}

func getColumnInfoData() []ColumnInfo {
	return []ColumnInfo{
		{Header: "AutoScaling Group Name", Type: Label, DataKey: "AutoScalingGroupName"},
		{Header: "Instance Type", Type: Label, DataKey: "InstanceTypes"},
		{Header: "Instances", Type: Label, DataKey: "DesiredCapacity"},
		{Header: "Cost $", Type: Label, DataKey: "HourlyCosts"},
		{Header: "Projected Cost $", Type: Label, DataKey: "ProjectedCosts"},
		{Header: "Projected Savings $", Type: Label, DataKey: "ProjectedSavings"},
		{Header: "Projected Savings %", Type: Label, DataKey: "ProjectedSavingsPercent"},
		{Header: "OnDemand %", Type: Entry, DataKey: "OnDemandPercentage", EntryValidator: validation.NewRegexp(`^([0-9]|[1-9][0-9]|100)$`, "Must contain an integer number between 0 and 100"), PlaceHolder: "0-100"},
		{Header: "OnDemand #", Type: Entry, DataKey: "OnDemandNumber", EntryValidator: validation.NewRegexp(`^([0-9]|[1-9][0-9]+)$`, "Must contain a natural number"), PlaceHolder: "Number"},
		{Header: "Enabled", Type: Check, DataKey: "Enabled"},
	}
}

func createTableWithHeaders(c *core.Launcher, data []ColumnInfo) *widget.Table {
	t := widget.NewTableWithHeaders(
		func() (int, int) {
			if c.Regions == nil || c.Regions[c.CurrentRegion] == nil ||
				c.Regions[c.CurrentRegion].AutoSpotting == nil {
				return 0, 0
			} else {
				return len(c.Regions[c.CurrentRegion].AutoSpotting.ASGs), len(data)
			}
		},
		func() fyne.CanvasObject {
			return container.NewStack(
				widget.NewLabel(""),
				widget.NewCheck("", func(bool) {}),
				widget.NewEntry(),
			)
		}, func(id widget.TableCellID, o fyne.CanvasObject) {})

	t.CreateHeader = func() fyne.CanvasObject {
		return newActiveHeader("")
	}

	t.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		header := o.(*ActiveHeader)
		header.TextStyle.Bold = true
		if id.Col >= 0 && id.Col < len(data) {
			header.SetText(data[id.Col].Header)
		}

		header.OnTapped = func() {
			log.Printf("Header %d tapped\n", id.Col)
		}
	}
	return t
}
func makeASGTable(w fyne.Window, c *core.Launcher) fyne.CanvasObject {
	data := getColumnInfoData()
	t := createTableWithHeaders(c, data)

	t.UpdateCell = func(id widget.TableCellID, o fyne.CanvasObject) {
		generateAutoSpottingTableData(&o, &id, &data, c)
		c.UpdateAutoSpottingTotals(c.CurrentRegion)
	}

	for i, col := range data {
		t.SetColumnWidth(i, float32(30+7*len(col.Header)))
	}

	return t
}

func generateAutoSpottingTableData(o *fyne.CanvasObject, id *widget.TableCellID, data *[]ColumnInfo, c *core.Launcher) {
	container := (*o).(*fyne.Container)
	label := container.Objects[0].(*widget.Label)
	check := container.Objects[1].(*widget.Check)
	entry := container.Objects[2].(*widget.Entry)

	// Clear previous state
	label.Hide()
	check.Hide()
	entry.Hide()

	if id.Row < 0 || id.Row >= len(c.Regions[c.CurrentRegion].AutoSpotting.ASGs) {
		return
	}

	colInfo := (*data)[id.Col]
	asg := c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row]

	// Determine cell width for this column (you might have this set elsewhere in your app)
	// cellWidth := int(t.ColumnWidth(id.Col)) - cellPadding

	// maxChars := estimateMaxChars(cellWidth, avgCharWidth)

	switch colInfo.Type {
	case Label:

		var text string
		switch colInfo.DataKey {
		case "AutoScalingGroupName":
			text = *asg.AutoScalingGroupName
		case "InstanceTypes":
			text = strings.Join(asg.InstanceTypes, ",")
		case "DesiredCapacity":
			text = fmt.Sprintf("%d", *asg.DesiredCapacity)
		case "HourlyCosts":
			text = formatFloat(asg.HourlyCosts * c.PricingIntervalMultiplier)
		case "ProjectedCosts":
			text = formatFloat(asg.ProjectedCosts * c.PricingIntervalMultiplier)
		case "ProjectedSavings":
			text = formatFloat(asg.ProjectedSavings * c.PricingIntervalMultiplier)
		case "ProjectedSavingsPercent":
			percentage := 0
			if asg.HourlyCosts > 0 {
				percentage = int(asg.ProjectedSavings / asg.HourlyCosts * 100)
			}
			text = fmt.Sprintf("%d%%", percentage)
		}
		// truncatedText := truncateTextToFitCell(text, maxChars)

		// label.SetText(truncatedText)
		label.SetText(text)
		label.Show()
	case Check:
		if colInfo.DataKey == "Enabled" {
			check.Show()
			check.SetChecked(asg.Enabled)
			check.OnChanged = func(checked bool) {
				// Update the ASG Enabled status based on checkbox
				asg.Enabled = checked
				// Additional logic to handle the change can be added here
			}
		}
	case Entry:
		entry.Show()
		entry.Validator = colInfo.EntryValidator
		entry.SetPlaceHolder(colInfo.PlaceHolder)
		var entryText string
		switch colInfo.DataKey {
		case "OnDemandPercentage":
			entryText = fmt.Sprintf("%.0f", asg.OnDemandPercentage)
		case "OnDemandNumber":
			entryText = fmt.Sprintf("%d", asg.OnDemandNumber)
		}
		entry.SetText(entryText)
		entry.OnChanged = func(text string) {
			switch colInfo.DataKey {
			case "OnDemandPercentage":
				if newValue, err := strconv.ParseFloat(text, 64); err == nil {
					// Assuming asg has a method to update on-demand percentage safely
					asg.OnDemandPercentage = newValue
					// Optionally, recalculate or update related fields
					// asg.CalculateHourlyPricing()
					// Remember to handle errors or invalid input as needed
					log.Printf("Updated OnDemandPercentage for %s to %.2f", *asg.AutoScalingGroupName, newValue)
				} else {
					log.Printf("Invalid OnDemandPercentage input: %s", text)
				}
			case "OnDemandNumber":
				if newValue, err := strconv.ParseInt(text, 10, 64); err == nil {
					// Assuming asg has a method to update on-demand number safely
					asg.OnDemandNumber = newValue
					// Optionally, recalculate or update related fields
					// asg.CalculateHourlyPricing()
					// Remember to handle errors or invalid input as needed
					log.Printf("Updated OnDemandNumber for %s to %d", *asg.AutoScalingGroupName, newValue)
				} else {
					log.Printf("Invalid OnDemandNumber input: %s", text)
				}
			}

			// After updating the ASG, you might want to trigger a refresh of your UI
			// or push changes to a server or configuration system as needed.
			// This is highly dependent on your application's architecture.
		}
	}
}

func ebsOptimizerRollout(a fyne.App, w fyne.Window, c *core.Launcher) *container.TabItem {
	// currentPrefRegion := a.Preferences().StringWithFallback(preferenceRegion, "us-east-1")

	// accessKey := widget.NewEntry()
	// accessKey.SetPlaceHolder("Usually starts with AKIA...")

	// secret := widget.NewEntry()
	// secret.SetPlaceHolder("Longer string")

	// sessionToken := widget.NewEntry()
	// sessionToken.SetPlaceHolder("(optional)")

	// regionsStaticAuth := widget.NewSelect(awsRegions(), func(s string) {
	// 	a.Preferences().SetString(preferenceRegion, s)
	// 	log.Println("selected region", s)
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

func rollout(w fyne.Window, c *core.Launcher) fyne.CanvasObject {

	a := fyne.CurrentApp()

	regions := widget.NewSelect(c.AWSRegions(), func(s string) {
		a.Preferences().SetString(preferenceAutoSpottingRolloutRegion, s)
		log.Println("selected AWS region", s)

		if p := a.Preferences().String(preferenceProfile); p != "" {
			log.Println("selected profile", p)
			c.ConnectWithProfileAuth(p)
			c.SetRegion(p)
		}

		if !c.Connected {
			err := errors.New("missing credentials")
			dialog.ShowError(err, w)
			return
		}

		if c.Regions == nil {
			c.Regions = make(map[string]*core.Region, 1)

			// c.Regions[s] = &core.Region{
			// 	AutoSpotting: &core.AutoSpotting{
			// 		services: c.Regions[s].
			// },

			//}
		}
		c.Regions[s].AutoSpotting.LoadASGData()
		c.SetRegion(s)
	})

	priceMode := widget.NewSelect([]string{"hourly", "monthly"}, func(s string) {
		a.Preferences().SetString(preferenceAutoSpottingPricingInterval, s)
		log.Println("selected pricing interval", s)

		c.SetPricingInterval(s)
	})

	pref := a.Preferences().String(preferenceAutoSpottingPricingInterval)
	if pref != "" {
		priceMode.SetSelected(pref)
	}

	odPercentage := widget.NewEntry()
	odPercentage.Validator = validation.NewRegexp(`^([0-9]|[1-9][0-9]|100)$`, "0 - 100")
	odPercentage.OnChanged = func(s string) {
		if odPercentage.Validate() != nil {
			return
		}
		log.Println("OnChanged OD percentage is valid", s)
		p, _ := strconv.ParseFloat(s, 64)

		if c.Regions != nil && c.Regions[c.CurrentRegion] != nil && c.Regions[c.CurrentRegion].AutoSpotting != nil {
			for _, asg := range c.Regions[c.CurrentRegion].AutoSpotting.ASGs {
				asg.OnDemandPercentage = p
				if asg.OnDemandPercentageEntry != nil {
					asg.OnDemandPercentageEntry.SetText(s)
				}
			}

		}
	}

	odNumber := widget.NewEntry()
	odNumber.Validator = validation.NewRegexp(`^([0-9]|[1-9][0-9]+)$`, "Invalid")

	odNumber.OnChanged = func(s string) {
		if odNumber.Validate() != nil {
			return
		}
		log.Println("OnChanged OD number is valid", s)
		n, _ := strconv.ParseInt(s, 10, 64)

		if c.Regions != nil && c.Regions[c.CurrentRegion] != nil && c.Regions[c.CurrentRegion].AutoSpotting != nil {

			for _, asg := range c.Regions[c.CurrentRegion].AutoSpotting.ASGs {
				asg.OnDemandNumber = n
				log.Printf("Setting Number Entry for ASG before nil check %s", *asg.AutoScalingGroupName)
				if asg.OnDemandNumberEntry != nil {
					log.Printf("Setting Number Entry for ASG passed nil check %s", *asg.AutoScalingGroupName)
					asg.OnDemandNumberEntry.SetText(s)
				}
			}
		}
	}

	convertCheck := widget.NewCheck("", func(set bool) {
		log.Printf("%v", set)
		if c.Regions != nil && c.Regions[c.CurrentRegion] != nil && c.Regions[c.CurrentRegion].AutoSpotting != nil {
			for _, asg := range c.Regions[c.CurrentRegion].AutoSpotting.ASGs {
				asg.Enabled = set
				if asg.ConvertToSpotCheck != nil {
					log.Printf("Setting Checkbox for ASG passed nil check %v", *asg.AutoScalingGroupName)
					asg.ConvertToSpotCheck.SetChecked(set)
				}
			}
		}
	})

	return container.NewBorder(
		container.NewHBox(
			container.NewGridWithRows(2,
				&widget.Form{
					Items: []*widget.FormItem{
						{Text: "AWS Region", Widget: regions, HintText: ""},
						//{Text: "Pricing interval", Widget: priceMode, HintText: ""},

					}},
				&widget.Form{
					Items: []*widget.FormItem{
						//	{Text: "AWS Region", Widget: regions, HintText: ""},
						{Text: "Pricing interval", Widget: priceMode, HintText: ""},
					}},

				&widget.Form{
					Items: []*widget.FormItem{
						//	{Text: "Override OnDemand Percentage", Widget: odPercentage, HintText: ""},
						{Text: "Override OnDemand Number", Widget: odNumber, HintText: ""},
					}},
				&widget.Form{
					Items: []*widget.FormItem{
						{Text: "Override OnDemand Percent", Widget: odPercentage, HintText: ""},
						//	{Text: "Override OnDemand Number", Widget: odNumber, HintText: ""},

					}},

				&widget.Form{
					Items: []*widget.FormItem{
						{Text: "Override Convert to Spot", Widget: convertCheck, HintText: ""},
						//	{Text: "Override OnDemand Number", Widget: odNumber, HintText: ""},

					}},

				// &widget.Form{
				// 	Items: []*widget.FormItem{
				// 		{Text: "Override Convert to Spot", Widget: convertCheck, HintText: ""},
				// 	}},
			)),

		container.NewBorder(
			nil, nil, nil,
			container.NewHBox(
				&widget.Form{
					Items: []*widget.FormItem{
						{Text: "Total current monthly costs", Widget: widget.NewLabelWithData(
							c.AutoSpottingCurrentTotalMonthlyCosts), HintText: ""},
						{Text: "Total projected monthly costs", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedMonthlyCosts), HintText: ""},
					},
				},
				&widget.Form{
					Items: []*widget.FormItem{

						{Text: "Total projected Spot monthly savings", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedSpotSavings), HintText: ""},
						{Text: "Total projected Spot savings percentage", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedSpotSavingsPercent), HintText: ""},
					},
				},
				&widget.Form{
					Items: []*widget.FormItem{
						{Text: "AutoSpotting charges (~10% of savings)", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedAutoSpottingCharges), HintText: ""},
						{Text: "Total Monthly net savings", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedNetSavings), HintText: ""},
					},
				},
				&widget.Form{
					Items: []*widget.FormItem{

						{Text: "", Widget: widget.NewButton("Generate AutoSpotting\n configuration", func() {
							c.ApplyAutoSpottingTags()
							dialog.ShowInformation("Information",

								"The configuration was persisted to your AutoScaling group tags. "+
									"In order for it to be applied, \nyou need to install AutoSpotting "+
									"from the AWS Marketplace using the link from the Welcome tab...",
								w)
						}), HintText: ""},
					},
				},
				// widget.NewButton("Apply\nconfiguration", func() {
				// 	log.Println("Apply mock...")
				// }),

			)),
		nil, nil,
		container.NewStack(container.NewAppTabs(
			autoSpottingRollout(a, w, c),
			//ebsOptimizerRollout(a, w, c),
		)),
	)
}

// selectEntry := widget.NewSelectEntry([]string{"Option A", "Option B", "Option C"})
// selectEntry.PlaceHolder = "Type or select"
// disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
// disabledCheck.Disable()
// checkGroup := widget.NewCheckGroup([]string{"CheckGroup Item 1", "CheckGroup Item 2"}, func(s []string) { log.Println("selected", s) })
// checkGroup.Horizontal = true
// radio := widget.NewRadioGroup([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { log.Println("selected", s) })
// radio.Horizontal = true
// disabledRadio := widget.NewRadioGroup([]string{"Disabled radio"}, func(string) {})
// disabledRadio.Disable()

// return container.NewVBox(
// 	widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) { log.Println("selected", s) }),
// 	selectEntry,
// 	widget.NewCheck("Check", func(on bool) { log.Println("checked", on) }),
// 	disabledCheck,
// 	checkGroup,
// 	radio,
// 	disabledRadio,
// 	widget.NewSlider(0, 100),
// )
