package core

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	"fyne.io/fyne/v2/widget"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2instancesinfo "github.com/cristim/ec2-instances-info"
)

type AutoSpotting struct {
	ASGs                       []*ASG
	services                   *services
	region                     *Region
	OverrideOnDemandPercentage int64
	OverrideOnDemandNumber     int64
	OverrideSpotConversion     bool
}

type ASG struct {
	types.AutoScalingGroup
	services                        *services
	HourlyCosts                     float64
	ProjectedCosts                  float64
	ProjectedSavings                float64
	InstanceTypes                   []string
	SpotInstanceNumber              int
	SpotInstancePercent             int
	region                          *Region
	ami                             string
	spotProduct                     *string
	Enabled                         bool
	OnDemandNumber                  int64
	OnDemandPercentage              float64
	OnDemandPercentageEntry         *widget.Entry
	OnDemandNumberEntry             *widget.Entry
	ConvertToSpotCheck              *widget.Check
	EnabledTagExistedInitially      bool
	ODNumberTagExistedInitially     bool
	ODPercentageTagExistedInitially bool
}

func (a *AutoSpotting) LoadASGData() error {

	input := &autoscaling.DescribeAutoScalingGroupsInput{}

	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(a.services.autoscaling, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Println("Error", err)
			return err
		}
		for _, asg := range output.AutoScalingGroups {
			asgData := ASG{
				AutoScalingGroup: asg,
				services:         a.services,
				region:           a.region,
			}
			for _, tag := range asg.Tags {
				switch *tag.Key {
				case "spot-enabled":
					{
						asgData.Enabled, _ = strconv.ParseBool(*tag.Value)
						asgData.EnabledTagExistedInitially = true
					}
				case "autospotting_min_on_demand_number":
					{
						asgData.OnDemandNumber, _ = strconv.ParseInt(*tag.Value, 10, 64)
						asgData.ODNumberTagExistedInitially = true
					}
				case "autospotting_min_on_demand_percentage":
					{
						n, _ := strconv.ParseFloat(*tag.Value, 64)
						if n > 100 {
							n = 100
						}
						if n < 0 {
							n = 0
						}
						asgData.OnDemandPercentage = n
						asgData.ODPercentageTagExistedInitially = true
					}

				}
			}

			fmt.Printf("%#v", asgData)

			asgData.populate()

			a.ASGs = append(a.ASGs, &asgData)
			fmt.Printf("AutoSpotting found ASG: %#v\n", *asg.AutoScalingGroupName)
		}

	}

	return nil
}

func (asg *ASG) populate() {
	asg.readASGConfiguration()

	err := asg.CalculateHourlyPricing()
	if err != nil {
		log.Printf("Couldn't determine hourly pricing for asg: %s, error: %s", *asg.AutoScalingGroupName, err.Error())
		return
	}
}

func (asg *ASG) readASGConfiguration() {
	if asg.LaunchConfigurationName != nil {
		resp, err := asg.services.autoscaling.DescribeLaunchConfigurations(context.TODO(),
			&autoscaling.DescribeLaunchConfigurationsInput{
				LaunchConfigurationNames: []string{*asg.LaunchConfigurationName},
			})
		if err != nil {
			fmt.Printf("Couldn't determint launch configuration for ASG: %#v\n", *asg.AutoScalingGroupName)
			return
		}
		asg.InstanceTypes = []string{*resp.LaunchConfigurations[0].InstanceType}
		asg.ami = *resp.LaunchConfigurations[0].ImageId
	}

	if asg.LaunchTemplate != nil {
		resp, err := asg.services.ec2.DescribeLaunchTemplateVersions(context.TODO(),
			&ec2.DescribeLaunchTemplateVersionsInput{
				LaunchTemplateName: asg.LaunchTemplate.LaunchTemplateName,
				Versions:           []string{*asg.LaunchTemplate.Version},
			})
		if err != nil {
			fmt.Printf("Couldn't determine launch template for ASG: %#v\n", *asg.AutoScalingGroupName)
			return
		}
		for _, lt := range resp.LaunchTemplateVersions {
			asg.InstanceTypes = []string{string(lt.LaunchTemplateData.InstanceType)}
			asg.ami = *lt.LaunchTemplateData.ImageId
		}

	}
	fmt.Printf("ASG Instance types: %v \n", asg.InstanceTypes)
}

func (asg *ASG) CalculateHourlyPricing() error {

	log.Printf("Calculating Hourly Pricing for ASG %s", *asg.AutoScalingGroupName)

	if asg.spotProduct == nil {
		spotProduct, err := asg.determineSpotProduct()
		if err != nil {
			log.Printf("Couldn't determine Operating System for asg: %s, error: %s", *asg.AutoScalingGroupName, err.Error())
			return err
		}
		asg.spotProduct = spotProduct
	}

	var odCosts, projectedCosts, projectedSavings float64
	for i, instance := range asg.Instances {
		pricing := asg.getHourlyPricing("cost", *instance.InstanceType, asg.region.name, *asg.spotProduct)
		if pricing == nil {
			log.Printf("Couldn't calculate hourly pricing for instance type %s in region %s", *instance.InstanceType, asg.region.name)
			continue
		}
		odCosts += pricing.OnDemand

		log.Printf("ASG %s on demand number %d and percentage %.2f, processing instance %d",
			*asg.AutoScalingGroupName, asg.OnDemandNumber, asg.OnDemandPercentage, i)

		keepOnDemand := int(math.Max(float64(asg.OnDemandNumber), float64(*asg.DesiredCapacity)*asg.OnDemandPercentage/100.0))

		if i >= keepOnDemand {
			log.Printf("ASG %s on demand number %d and percentage %.2f, adding instance number %d",
				*asg.AutoScalingGroupName, asg.OnDemandNumber, asg.OnDemandPercentage, i)
			projectedCosts += pricing.SpotMin
			projectedSavings += (pricing.OnDemand - pricing.SpotMin)
		} else {
			projectedCosts += pricing.OnDemand
		}
	}
	asg.HourlyCosts = odCosts
	asg.ProjectedCosts = projectedCosts
	asg.ProjectedSavings = projectedSavings

	fmt.Printf("ASG current OnDemand costs %v, projected costs: %v, savings: %v \n", odCosts, projectedCosts, projectedSavings)

	return nil
}

func (asg *ASG) getHourlyPricing(figureType, instanceType, region, spotProduct string) *ec2instancesinfo.Pricing {
	var ret ec2instancesinfo.Pricing

	//log.Printf("ASG: %#v", asg)
	//log.Printf("ASG region: %#v", asg.region)
	//log.Printf("ASG region instance type data: %#v", asg.region.instanceTypeData)

	for _, i := range *asg.region.instanceTypeData {
		if i.InstanceType != instanceType {
			continue
		}

		//log.Printf("Found instance type, %#v with instance type information in %v %#v", i.InstanceType, region, i)

		pricing := i.Pricing[region]

		switch spotProduct {
		case "Linux/UNIX":
			ret = pricing.Linux
		case "Windows":
			ret = pricing.MSWin
		case "Red Hat Enterprise Linux":
			ret = pricing.RHEL
		case "SUSE Linux":
			ret = pricing.SLES

		}

		log.Printf("Hourly pricing information: %#v", ret)
		break
	}
	return &ret
}

func (asg *ASG) determineSpotProduct() (*string, error) {

	resp, err := asg.services.ec2.DescribeImages(context.TODO(),
		&ec2.DescribeImagesInput{
			ImageIds: []string{asg.ami},
		})
	if err != nil {
		log.Printf("Couldn't describe image %s: %s", asg.ami, err.Error())
		return nil, err
	}
	product := resp.Images[0].PlatformDetails

	log.Printf("Spot Product: %s", *product)
	return product, err
}

func checkRI(instanceID string, ec2Svc *ec2.Client) (bool, error) {
	dir, err := ec2Svc.DescribeInstances(context.TODO(),
		&ec2.DescribeInstancesInput{
			InstanceIds: []string{instanceID},
		})
	if err != nil {
		return false, fmt.Errorf("error describing instances: %v", err)
	}

	reservationID := dir.Reservations[0].ReservationId

	if reservationID == nil {
		return false, fmt.Errorf("error finding reserved instance ID: %v", err)
	}

	drir, err := ec2Svc.DescribeReservedInstances(context.TODO(),
		&ec2.DescribeReservedInstancesInput{
			ReservedInstancesIds: []string{*reservationID},
		})
	if err != nil {
		return false, fmt.Errorf("error describing reserved instances: %v", err)
	}

	if len(drir.ReservedInstances) > 0 {
		return true, nil
	}

	return false, nil
}
