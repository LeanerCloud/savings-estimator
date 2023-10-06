# Spot Savings Estimator

Welcome to the LeanerCloud AWS Spot cost savings estimation tool!

This is a native desktop application that allows you to estimate the cost
savings you can achieve in your AWS account by converting your AutoScaling
Groups to Spot instances.

You can simulate various scenarios, such as to keep some of your instances as
OnDemand in each group (maybe covered by Reserved Instances or Savings Plans),
or only convert some of your AutoScaling Groups to Spot as part of a gradual
rollout.

## Demo

[Demo](https://www.youtube.com/watch?v=VXfCOXXtLwA)

And here's a screenshot from that demo:

![Screenshot](https://github.com/LeanerCloud/savings-estimator/assets/95209/897a2fb4-2838-401a-bac0-53b67b1eb77d)



## Precompiled Binaries

Binaries for Windows and Linux are available at
[Releases](https://github.com/LeanerCloud/savings-estimator/releases).

On Linux you may need to make them executable after you download them, and it's
recommended to put them in one of the directories listed in the `PATH` variable

## Install from source code

On any OS you should be able to build and install it from source if you have Go
installed:

`go install github.com/LeanerCloud/savings-estimator@latest`

Then of Unix-like systems you can run it from `$HOME/go/bin/savings-estimator`
or just `savings-estimator` if `$HOME/go/bin` is already in your PATH.

On Windows you should similarly find it in `%USERPROFILE%\go\bin\savings-estimator.exe`

## Credential management

- It assumes you have some AWS credentials configured in the Configuration
view, either as profiles sourced from the AWS CLI/SDK config file, or use a
access key/secret for one-off execution. - The selected profile configuration
is persisted across runs in the Fyne config path, but pasted access key and
secrets are ephemeral and only used for the curent run.

## Required IAM permissions

It's recommended to use an IAM role with limited permissions.

The only permissions required at the moment are listed below, but these are
subject to change over time as new features are implemented which may require
more permissions:

```text
autoscaling:CreateOrUpdateTags
autoscaling:DescribeAutoScalingGroups
ec2:DescribeImages
ec2:DescribeInstances
```

You can also use our CloudFomation [template](/cloudformation/template.yaml) to
create a role with these permissions, and then assume it using the following
snippet added to your `.aws/config`:

```ini
[profile SavingsEstimator]
role_arn = arn:aws:iam::ACCOUNT_ID:role/NAME_OF_THE_ROLE # You can see it in the CloudFormation stack outputs.
source_profile = default # or any other profile from the same AWS account
```

## Integration with AutoSpotting

Spot Savings Estimator can be executed independent of AutoSpotting for cost
savings simulation/estimation purposes, maybe to evaluate the opportunity for
your particular setup.

You may then use any mechanism to adopt Spot, such as applying the
configuration yourself group by group as per your simulation.

But for more convenience you may use [AutoSpotting](AutoSpotting.io), our state
of the art cost optimization engine for Spot.

AutoSpotting allows you to apply the scenario you simulated with a single click
for minimal time and effort spent, and also in the end getting a more reliable
setup.

### About AutoSpotting

AutoSpotting allows you to adopt Spot instances with all the Spot best
practices recommended by AWS:

- wide diversification over multiple instance types.
- uses a capacity optimized allocation strategy to reduce the frequency of
  Spot interruptions.

In addition, AutoSpotting also prioritizes for lower cost instances from newer
generations and implements a reliable failover to on-demand instances when
running out of Spot capacity, which the native ASGs won't do.

### How the integration with AutoSpotting works

AutoSpotting uses tags as configuration mechanism, and most of the times it
works without requiring configuration changes on your OnDemand AutoScaling
groups, as long as they fit its main requirements/recommendations:

- Use Launch Template/Configuration without instance type overrides
- Span across all AZs from the region

The Savings Estimator can conveniently create the AutoSpotting configuration
tags with a single click, so that AutoSpotting will implement the simulated
scenario without the need for additional configuration changes.

### Applying the simulation with AutoSpotting

You can create the AutoSpotting configuration tags by clicking the "Generate
AutoSpotting configuration" button on the bottom right corner in the Savings
view.

These configurations will be persisted as tags on your ASGs, but nothing else
will happen until AutoSpotting is installed in the AWS account.

The latest version of AutoSpotting is available on the
[AWS Marketplace](https://aws.amazon.com/marketplace/pp/prodview-6uj4pruhgmun6),
and you will need to follow these installation instructions to install it:

- Continue to Subscribe/Configuration/Launch
- Install AutoSpotting using either CloudFormation or Terraform from the
  "Launch this software" view.

You may just use the default parameters and adjust them later if needed.

Once AutoSpotting is installed, any settings created as ASG tags through
Savings Estimator will be gradually applied on your AutoScaling groups.

For more details about AutoSpotting, see [AutoSpotting.io](AutoSpotting.io).

## Local development

You need to have Go installed, then it's as easy as running

`go run .`

## Known issues

- The Fyne GUI library we use requires a graphics card with support for
  OpenGL 2.x,. Unfortunately VMs and cloud computing instances typically
  lack the required GPU drivers and will fail to run the Savings Estimator
  As a workaround you may get it work using
  [mesa-dist-win](https://github.com/pal1000/mesa-dist-win#installation-and-usage),
  which offers OpenGL drivers with software rendering support.

## Contributions

Any contributions are welcome through the usual GitHub mechanisms (Issues, Pull
Requests, Discussions, etc.)

## Future plans

Please refer to our public [roadmap](https://github.com/orgs/LeanerCloud/projects/1).

## License

This software is available under the AGPL-3 Open Source license.

## Credits

Savings Estimator is proudly written in Go using the [fyne](fyne.io) GUI
toolkit and leverages a lot of OSS code under the hood. Thanks to everyone tho
worked on those projects for their hard work!
