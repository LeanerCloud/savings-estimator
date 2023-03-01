# LeanerCloud GUI - Spot Savings Estimator GUI application.


Welcome to the LeanerCloud AWS Spot cost savings estimation tool!

This tool allows you to estimate the cost savings you can achieve in your AWS account by converting your AutoScaling Groups to Spot instances.

You can select various scenarios, such as to keep some of your instances as OnDemand in each group (maybe covered by Reserved Instances or Savings Plans),
or only convert some of your AutoScaling Groups to Spot as part of a gradual rollout.

You may use any mechanism to adopt Spot, such as converting the configuration yourself group by group as per what you defined in this tool.

For your convenience, you can also use [AutoSpotting](AutoSpotting.io), our state of the art cost optimization engine for Spot. 

AutoSpotting is tightly integrated with
this cost savings estimator, so you can apply this configuration with a single click, by tagging your ASGs as expected by AutoSpotting.

AutoSpotting allows you to adopt Spot instances with all the Spot best practices recommended by AWS:
-  wide diversification over multiple instance types
-  uses a capacity optimized allocation strategy to reduce the frequency of Spot interruptions.


In addition, AutoSpotting also prioritizes for lower cost instances from newer generations and implements a reliable failover to on-demand instances when running out of Spot capacity.

In most situations AutoSpotting doesn't require any configuration changes to your AutoScaling Groups, but uses the existing launch template or launch configuration.


For more details about AutoSpotting, see [LeanerCloud.com](LeanerCloud.com).

## Demo

[![LeanerCloud GUI Demo](https://img.youtube.com/vi/2D6IMm6dFDo/0.jpg)](https://www.youtube.com/watch?v=2D6IMm6dFDo)

(click above to play the demo video)

## Precompiled Binaries

Binaries for Windows and Linux are available at [Releases](https://github.com/LeanerCloud/leaner-cloud-gui/releases).

On Linux you may need to make them executable after you download them, and it's recommended to put them somewhere in the PATH 

## Install from source code

On any other OS you should be able to install it from source if you have Go installed:

`go install github.com/LeanerCloud/leaner-cloud-gui@latest`

## Credential management

- It assumes you have some AWS credentials configured in the Configuration view, either as profiles or pasted access key/secret.
- The profile configuration is persisted across runs in the Fyne config path (macOS: ~/Library/Preferences/fyne, Linux or BSD: ~/.config/fyne, Windows: ~\AppData\Roaming\fyne), together with some additional settings. 
- The pasted access key and secret is ephemeral and only used for the curent run.

## Required IAM permissions

It's recommended to use an IAM role with limited permissions. 

The only permissions required at the moment are listed below, but these are subject to change over time as new features are implemented which may require more permissions:
```
autoscaling:CreateOrUpdateTags
autoscaling:DescribeAutoScalingGroups
ec2:DescribeImages
ec2:DescribeInstances
```

## Dependency on AutoSpotting

`leaner-cloud-gui` can be executed independent of AutoSpotting for cost savings simulation/estimation and configuration purposes.
These configurations will be persisted as tags on your ASGs, but nothing else will happen unless AutoSpotting is installed in the AWS account.

The latest version of AutoSpotting is available on the [AWS Marketplace](https://aws.amazon.com/marketplace/pp/prodview-6uj4pruhgmun6), and you will need to follow the installation instructions:
- Continue to Subscribe/Configuration/Launch
- Install AutoSpotting using either CloudFormation or Terraform from the "Launch this software" view

Once AutoSpotting is installed, any settings created as ASG tags through `leaner-cloud-gui` will be gradually applied on your AutoScaling groups.

## Local development

`go run .`


## Future plans

Please refer to our public [roadmap](https://github.com/orgs/LeanerCloud/projects/1)

## License

This tool is available under the AGPL-3 Open Source license.


-- 

This tool is proudly written in Go using the [fyne](fyne.io) GUI toolkit.
