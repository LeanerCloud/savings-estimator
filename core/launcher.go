package core

import (
	"context"
	"fmt"
	"log"
	"math"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	ec2instancesinfo "github.com/LeanerCloud/ec2-instances-info"
	"gopkg.in/ini.v1"
)

type Launcher struct {
	Regions                   map[string]*Region
	CurrentRegion             string
	GlobalServices            *globalServices
	Connected                 bool
	InstanceTypeData          *ec2instancesinfo.InstanceData
	PricingIntervalMultiplier float64

	AutoSpottingCurrentTotalMonthlyCosts     binding.String
	AutoSpottingProjectedMonthlyCosts        binding.String
	AutoSpottingProjectedSpotSavings         binding.String
	AutoSpottingProjectedSpotSavingsPercent  binding.String
	AutoSpottingProjectedAutoSpottingCharges binding.String
	AutoSpottingProjectedNetSavings          binding.String
}

type Region struct {
	services         *services
	AutoSpotting     *AutoSpotting
	Launcher         *Launcher
	name             string
	instanceTypeData *ec2instancesinfo.InstanceData
}

// TODO: remove the hardcoded region list
func (c *Launcher) AWSRegions() []string {
	return []string{
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-north-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}
}

func (c *Launcher) ReadAWSProfiles() []string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	path := filepath.Join(homeDir, ".aws/config")

	cfg, err := ini.Load(path)
	if err != nil {
		log.Printf("Fail to read AWS Credentials file: %v", err)
		return []string{}
	}

	profiles := cfg.SectionStrings()[1:]

	for i, profile := range profiles {
		profiles[i] = strings.Replace(profile, "profile ", "", 1)
	}

	sort.Strings(profiles)
	return profiles
}

func (c *Launcher) Connect(configOption config.LoadOptionsFunc) {
	mainRegion := "us-east-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		configOption,
		config.WithRegion(mainRegion),
	)
	if err != nil {
		log.Printf("unable to load SDK config from profile , %v", err)
	}

	log.Println("Connecting global services in region", mainRegion)

	s := globalServices{
		config:         cfg,
		cloudformation: cloudformation.NewFromConfig(cfg),
		s3:             s3.NewFromConfig(cfg),
	}

	c.GlobalServices = &s

	c.Regions = make(map[string]*Region, 0)

	for _, r := range c.AWSRegions() {

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			configOption,
			config.WithRegion(r),
		)
		if err != nil {
			log.Printf("unable to load SDK config from profile , %v", err)
		}

		log.Println("Connecting services in region", r)

		s := services{
			config:      cfg,
			autoscaling: autoscaling.NewFromConfig(cfg),
			ec2:         ec2.NewFromConfig(cfg),
			//cfn:         cloudformation.NewFromConfig(cfg)
		}

		c.Regions[r] = &Region{
			name:     r,
			services: &s,
			AutoSpotting: &AutoSpotting{
				services: &s,
			},
			instanceTypeData: c.InstanceTypeData,
		}
		c.Regions[r].AutoSpotting.region = c.Regions[r]

	}
	c.Connected = true
}

func (c *Launcher) ConnectWithProfileAuth(profile string) {
	co := config.WithSharedConfigProfile(profile)
	c.Connect(co)
}

func (c *Launcher) ConnectWithStaticAuth(key, secret, token string) {
	co := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, token))
	c.Connect(co)
}

func (c *Launcher) SetRegion(region string) {
	c.CurrentRegion = region
}

func (c *Launcher) SetPricingInterval(s string) {
	if s == "monthly" {
		c.PricingIntervalMultiplier = 730.0

	} else {
		c.PricingIntervalMultiplier = 1.0
	}
}

func (c *Launcher) UpdateAutoSpottingTotals(region string) {
	var currentTotalMonthlyCosts, projectedMonthlyCosts, projectedSpotSavings, projectedSpotSavingsPercent, projectedNetSavings float64

	if c.Regions == nil || c.Regions[region] == nil || c.Regions[region].AutoSpotting == nil || len(c.Regions[region].AutoSpotting.ASGs) == 0 {
		return
	}

	for _, asg := range c.Regions[region].AutoSpotting.ASGs {

		currentTotalMonthlyCosts += asg.HourlyCosts * 730
		if !asg.Enabled {
			projectedMonthlyCosts += asg.HourlyCosts * 730
			continue
		}
		projectedMonthlyCosts += asg.ProjectedCosts * 730
		projectedSpotSavings += asg.ProjectedSavings * 730

	}

	projectedAutoSpottingCharges := math.Floor(projectedSpotSavings/14.6) * 0.73
	projectedSpotSavingsPercent = projectedSpotSavings / currentTotalMonthlyCosts

	projectedNetSavings = projectedSpotSavings - projectedAutoSpottingCharges

	c.AutoSpottingCurrentTotalMonthlyCosts.Set(fmt.Sprintf("%.2f", currentTotalMonthlyCosts))
	c.AutoSpottingProjectedMonthlyCosts.Set(fmt.Sprintf("%.2f", projectedMonthlyCosts))
	c.AutoSpottingProjectedSpotSavings.Set(fmt.Sprintf("%.2f", projectedSpotSavings))
	c.AutoSpottingProjectedSpotSavingsPercent.Set(fmt.Sprintf("%d%%", int(projectedSpotSavingsPercent*100)))
	c.AutoSpottingProjectedAutoSpottingCharges.Set(fmt.Sprintf("%.2f", projectedAutoSpottingCharges))
	c.AutoSpottingProjectedNetSavings.Set(fmt.Sprintf("%.2f", projectedNetSavings))

}

func (c *Launcher) ApplyAutoSpottingTags() {
	log.Printf("Appling tags for all ASGs")
	if c.Regions == nil || c.CurrentRegion == "" || c.Regions[c.CurrentRegion] == nil || c.Regions[c.CurrentRegion].AutoSpotting == nil || len(c.Regions[c.CurrentRegion].AutoSpotting.ASGs) == 0 {
		log.Printf("Appling tags for all ASGs failes nil checks")
		return
	}

	as := c.Regions[c.CurrentRegion].AutoSpotting
	for _, asg := range as.ASGs {
		log.Printf("Appling tags for ASG %s", *asg.AutoScalingGroupName)

		// ResourceId:         []string{*asg.AutoScalingGroupName},
		// Tags:              tags,
		// ResourceType:      autoscaling.ResourceTypeAutoScalingGroup,
		// PropagateAtLaunch: aws.Bool(true),

		resourceType := "auto-scaling-group"
		propagate := false
		var tags []types.Tag

		if asg.Enabled || asg.EnabledTagExistedInitially {
			key := "spot-enabled"
			value := fmt.Sprintf("%t", asg.Enabled)
			tags = append(tags, types.Tag{
				Key:   &key,
				Value: &value,
			})
		}

		if asg.OnDemandNumber > 0 || asg.ODNumberTagExistedInitially {
			key := "autospotting_min_on_demand_number"
			value := fmt.Sprintf("%d", asg.OnDemandNumber)
			tags = append(tags, types.Tag{
				Key:   &key,
				Value: &value,
			})
		}

		if asg.OnDemandPercentage > 0 || asg.ODPercentageTagExistedInitially {
			key := "autospotting_min_on_demand_percentage"
			value := fmt.Sprintf("%.2f", float64(asg.OnDemandPercentage))
			tags = append(tags, types.Tag{
				Key:   &key,
				Value: &value,
			})
		}

		for i, _ := range tags {
			tags[i].ResourceId = asg.AutoScalingGroupName
			tags[i].ResourceType = &resourceType
			tags[i].PropagateAtLaunch = &propagate
		}

		//	spew.Dump("Tags: %#v", tags)

		_, err := as.services.autoscaling.CreateOrUpdateTags(context.TODO(), &autoscaling.CreateOrUpdateTagsInput{
			Tags: tags,
		})

		if err != nil {
			log.Printf("Could not create tags for AutoScalingGroup %s, error: %s", *asg.AutoScalingGroupName, err.Error())
		}
	}
}
