package screens

import (
	"errors"
	"fmt"
	"leanercloud/core"
	"log"
	"strconv"
	"strings"

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
	Header             string
	Type               widgetType
	LabelContentFunc   func(id *widget.TableCellID, label *widget.Label)
	EntryContentFunc   func(id *widget.TableCellID, entry *widget.Entry)
	CheckContentFunc   func(id *widget.TableCellID, check *widget.Check)
	CheckOnChangedFunc func(id *widget.TableCellID) (f func(set bool))
	EntryValidator     fyne.StringValidator
	PlaceHolder        string
	EntryOnChangedFunc func(id *widget.TableCellID, entry *widget.Entry) (f func(s string))
	EntryLinkToASGFunc func(id *widget.TableCellID, entry *widget.Entry)
	LinkCheckToASGFunc func(id *widget.TableCellID, check *widget.Check)
}

func autoSpottingRollout(a fyne.App, w fyne.Window, c *core.Launcher) *container.TabItem {

	return container.NewTabItem("AutoSpotting", makeASGTable(w, c))
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

func makeASGTable(_ fyne.Window, c *core.Launcher) fyne.CanvasObject {

	data := []ColumnInfo{
		{
			Header: "",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(fmt.Sprintf("%d", id.Row))
			},
		},
		{
			Header: "AutoScaling Group Name",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(*c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].AutoScalingGroupName)
			},
		},
		{
			Header: "Instance Type",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(strings.Join(c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].InstanceTypes, ","))
			},
		},
		{
			Header: "Instances",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(fmt.Sprintf("%d",
					*c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].DesiredCapacity))
			},
		},
		{
			Header: "Cost $",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {

				label.SetText(formatFloat(
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].HourlyCosts * c.PricingIntervalMultiplier))
			},
		},

		{

			Header: "Projected Cost $",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(formatFloat(
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].ProjectedCosts * c.PricingIntervalMultiplier))
			},
		},

		{
			Header: "Projected Savings $",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(formatFloat(
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].ProjectedSavings * c.PricingIntervalMultiplier))
			},
		},
		{
			Header: "Projected Savings %",
			Type:   Label,
			LabelContentFunc: func(id *widget.TableCellID, label *widget.Label) {
				label.SetText(fmt.Sprintf("%d%%",
					int(c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].ProjectedSavings/
						c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].HourlyCosts*100)))
			},
		},

		{
			Header:         "OnDemand %",
			Type:           Entry,
			EntryValidator: validation.NewRegexp(`^([0-9]|[1-9][0-9]|100)$`, "Must contain an integer number between 0 and 100"),
			PlaceHolder:    "0-100",
			EntryContentFunc: func(id *widget.TableCellID, entry *widget.Entry) {
				entry.SetText(fmt.Sprintf("%.0f",
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].OnDemandPercentage))
			},
			EntryOnChangedFunc: func(id *widget.TableCellID, entry *widget.Entry) (f func(s string)) {
				//log.Printf("Entry for cell %d, %d changed to %v", id.Row, id.Col, entry.Text)
				return func(s string) {
					var n float64
					if entry.Validate() == nil {
						log.Println("Data is valid", s)
						n, _ = strconv.ParseFloat(s, 64)
					}

					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].OnDemandPercentage = n
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].CalculateHourlyPricing()
					c.UpdateAutoSpottingTotals(c.CurrentRegion)
				}
			},
			EntryLinkToASGFunc: func(id *widget.TableCellID, entry *widget.Entry) {
				if asg := c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1]; asg.OnDemandPercentageEntry == nil {
					log.Printf("Setting Percentage entry to ASG %v", *asg.AutoScalingGroupName)
					asg.OnDemandPercentageEntry = entry
				}
			},
		},
		{
			Header:         "OnDemand #",
			Type:           Entry,
			EntryValidator: validation.NewRegexp(`^([0-9]|[1-9][0-9]+)$`, "Must contain a natural number"),
			PlaceHolder:    "Number",
			EntryContentFunc: func(id *widget.TableCellID, entry *widget.Entry) {
				entry.SetText(fmt.Sprintf("%d",
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].OnDemandNumber))
			},
			EntryOnChangedFunc: func(id *widget.TableCellID, entry *widget.Entry) (f func(s string)) {
				//log.Printf("Entry for cell %d, %d changed to %v", id.Row, id.Col, entry.Text)
				return func(s string) {
					var n int64
					if entry.Validate() == nil {
						log.Println("Data is valid", s)
						n, _ = strconv.ParseInt(s, 10, 64)
					}

					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].OnDemandNumber = n
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].CalculateHourlyPricing()
					c.UpdateAutoSpottingTotals(c.CurrentRegion)
				}
			},
			EntryLinkToASGFunc: func(id *widget.TableCellID, entry *widget.Entry) {
				if asg := c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1]; asg.OnDemandNumberEntry == nil {
					log.Printf("Setting Number entry to ASG %v", *asg.AutoScalingGroupName)
					asg.OnDemandNumberEntry = entry
				}
			},
		},

		{
			Header: "Enabled",
			Type:   Check,
			CheckContentFunc: func(id *widget.TableCellID, check *widget.Check) {
				check.SetChecked(c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].Enabled)
			},
			CheckOnChangedFunc: func(id *widget.TableCellID) (f func(set bool)) {
				return func(set bool) {
					log.Printf("Checkbox for cell %d, %d changed to %v", id.Row, id.Col, set)
					c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1].Enabled = set
					c.UpdateAutoSpottingTotals(c.CurrentRegion)
				}
			},
			LinkCheckToASGFunc: func(id *widget.TableCellID, check *widget.Check) {
				if asg := c.Regions[c.CurrentRegion].AutoSpotting.ASGs[id.Row-1]; asg.ConvertToSpotCheck == nil {
					log.Printf("Setting Check to ASG %v", *asg.AutoScalingGroupName)
					asg.ConvertToSpotCheck = check
				}
			},
		},
	}

	t := widget.NewTable(
		func() (int, int) {
			if c.Regions == nil || c.Regions[c.CurrentRegion] == nil ||
				c.Regions[c.CurrentRegion].AutoSpotting == nil {
				return 1, 1
			} else {
				return len(c.Regions[c.CurrentRegion].AutoSpotting.ASGs) + 1, len(data)
			}
		},
		func() fyne.CanvasObject {
			return container.NewMax(
				widget.NewLabel(""),
				widget.NewCheck("", func(bool) {}),
				widget.NewEntry(),
			)
		}, func(id widget.TableCellID, o fyne.CanvasObject) {})

	t.UpdateCell = func(id widget.TableCellID, o fyne.CanvasObject) {
		generateAutoSpottingTableData(&o, &id, &data)
		c.UpdateAutoSpottingTotals(c.CurrentRegion)
	}

	for i, col := range data {
		t.SetColumnWidth(i, float32(30+7*len(col.Header)))
	}

	return t
}

func generateAutoSpottingTableData(o *fyne.CanvasObject, id *widget.TableCellID, data *[]ColumnInfo) {

	//log.Printf("Generating table cell row: %d, col: %d", id.Row, id.Col)

	label := (*o).(*fyne.Container).Objects[0].(*widget.Label)
	check := (*o).(*fyne.Container).Objects[1].(*widget.Check)
	entry := (*o).(*fyne.Container).Objects[2].(*widget.Entry)

	if id.Row == 0 {
		check.Hide()
		entry.Hide()
		label.TextStyle = fyne.TextStyle{Bold: true}
		label.Alignment = fyne.TextAlignCenter
		label.SetText((*data)[id.Col].Header)
		label.Show()
		return
	}

	switch (*data)[id.Col].Type {
	case Label:
		{
			label.Show()
			check.Hide()
			entry.Hide()
			(*data)[id.Col].LabelContentFunc(id, label)
		}
	case Check:
		{
			(*data)[id.Col].CheckContentFunc(id, check)
			check.OnChanged = (*data)[id.Col].CheckOnChangedFunc(id)
			(*data)[id.Col].LinkCheckToASGFunc(id, check)

			check.Show()
			entry.Hide()
			label.Hide()
		}
	case Entry:
		{
			entry.Validator = (*data)[id.Col].EntryValidator
			entry.SetPlaceHolder((*data)[id.Col].PlaceHolder)
			(*data)[id.Col].EntryContentFunc(id, entry)
			entry.OnChanged = (*data)[id.Col].EntryOnChangedFunc(id, entry)
			(*data)[id.Col].EntryLinkToASGFunc(id, entry)

			label.Hide()
			check.Hide()
			entry.Show()

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

func rollout(w fyne.Window, c *core.Launcher) fyne.CanvasObject {

	a := fyne.CurrentApp()

	regions := widget.NewSelect(c.AWSRegions(), func(s string) {
		a.Preferences().SetString(preferenceAutoSpottingRolloutRegion, s)
		fmt.Println("selected AWS region", s)

		if p := a.Preferences().String(preferenceProfile); p != "" {
			fmt.Println("selected profile", p)
			c.ConnectWithProfileAuth(p)
			c.SetRegion(p)
		}

		if !c.Connected {
			err := errors.New("missing credentials")
			dialog.ShowError(err, w)
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
		fmt.Println("selected pricing interval", s)

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
						{Text: "AutoSpotting charges", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedAutoSpottingCharges), HintText: ""},
						{Text: "Total Monthly net savings", Widget: widget.NewLabelWithData(
							c.AutoSpottingProjectedNetSavings), HintText: ""},
					},
				},
				&widget.Form{
					Items: []*widget.FormItem{

						{Text: "", Widget: widget.NewButton("Apply\nconfiguration", func() {
							c.ApplyAutoSpottingTags()
						}), HintText: ""},
					},
				},
				// widget.NewButton("Apply\nconfiguration", func() {
				// 	log.Println("Apply mock...")
				// }),

			)),
		nil, nil,
		container.NewMax(container.NewAppTabs(
			autoSpottingRollout(a, w, c),
			//ebsOptimizerRollout(a, w, c),
		)),
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
