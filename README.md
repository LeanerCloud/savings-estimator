# leaner-cloud-gui

Desktop application for Windows, Mac and Linux that can be used to configure [LeanerCloud](LeanerCloud.com) tools such as [AutoSpotting](AutoSpotting.io), and later [EBS Optimizer](https://leanercloud.com/ebs-optimizer).

## Demo

[![LeanerCloud GUI Demo](https://img.youtube.com/vi/2D6IMm6dFDo/0.jpg)](https://www.youtube.com/watch?v=2D6IMm6dFDo)

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
