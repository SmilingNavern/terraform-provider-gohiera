
# Terraform Go Hiera Provider
This provider is forked from puppet hiera provider and implements data sources that can be used to perform hierachical data lookups with go implementation of [Hiera](https://github.com/lyraproj/hiera).

This is useful for providing configuration values in an environment with a high level of dimensionality.

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.10.x
* [Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)
* [Hiera written in go](https://github.com/lyraproj/hiera)

## Usage
To configure the provider:
```hcl
provider "hiera" {}
```
Hiera isn't very useful without scope variables to inform its lookups:
```hcl
provider "hiera" {
  scope {
    environment = "live"
    service     = "api"
  }
}
```
By default, the provider expects `lookup` to be available in your `$PATH`. It will also look for the configuration file at current working directory with name `hiera.yaml`. You can override both those values if you so wish:
```hcl
provider "hiera" {
  config = "~/hiera.yaml"
  bin    = "/usr/local/bin/lookup"
  scope {
    environment = "live"
    service     = "api"
  }
}
```

You can configure merge strategy for variables as well:
```hcl
provider "hiera" {
  merge = "deep"
}
```

## Data Sources
This provider only implements data sources.
#### Hash
To retrieve a hash:
```hcl
data "hiera_hash" "aws_tags" {
  key = "aws_tags"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the hash, represented as a map

Terraform doesn't support nested maps or other more complex data structures. Any keys containing nested elements won't be returned.

#### Array
To retrieve an array:
```hcl
data "hiera_array" "java_opts" {
  key = "java_opts"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the array (list)

#### Value
To retrieve any other flat value:
```hcl
data "hiera" "aws_cloudwatch_enable" {
  key = "aws_cloudwatch_enable"
}
```
The following output parameters are returned:
* `id` - matches the key
* `key` - the queried key
* `value` - the value

All values are returned as strings because Terraform doesn't implement other types like int, float or bool. The values will be implicitly converted into the appropriate type depending on usage.

## Example
Here's an example of setting different values and data types at multiple levels in Hiera and then retrieving those values as data sources for use in outputs.

**hiera.yaml**
```yaml
---
:backends:
  - yaml

:hierarchy:
  - service/%{service}
  - environment/%{environment}
  - common

:yaml:
  :datadir: /data/
```
**/data/common.yaml**
```yaml
---
aws_cloudwatch_enable: false
aws_instance_size: t2.micro
aws_tags: {}
java_opts: []
```
**/data/environment/live.yaml**
```yaml
---
aws_cloudwatch_enable: true
aws_tags:
  tier: 1
java_opts:
  - '-Dspring.profiles.active=live'
```
**/data/service/api.yaml**
```yaml
---
aws_instance_size: t2.large
aws_tags:
  team: A
java_opts:
  - '-Xms512m'
  - '-Xmx2g'
```
**main.tf**
```hcl
provider "hiera" {
  scope {
    environment = "live"
    service = "api"
  }
}

data "hiera" "aws_instance_size" {
  key = "aws_instance_size"
}

data "hiera" "aws_cloudwatch_enable" {
  key = "aws_cloudwatch_enable"
}

data "hiera_hash" "aws_tags" {
  key = "aws_tags"
}

data "hiera_array" "java_opts" {
  key = "java_opts"
}

output "aws_instance_size" {
  value = "${data.hiera.aws_instance_size.value}"
}

output "aws_cloudwatch_enable" {
  value = "${data.hiera.aws_cloudwatch_enable.value}"
}

output "aws_tags" {
  value = "${data.hiera_hash.aws_tags.value}"
}

output "java_opts" {
  value = "${data.hiera_array.java_opts.value}"
}

```
Then, plan and apply:
```sh
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.hiera_hash.aws_tags: Refreshing state...
data.hiera_array.java_opts: Refreshing state...
data.hiera.aws_cloudwatch_enable: Refreshing state...
data.hiera.aws_instance_size: Refreshing state...

------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.


$ terraform apply
data.hiera.aws_cloudwatch_enable: Refreshing state...
data.hiera_array.java_opts: Refreshing state...
data.hiera_hash.aws_tags: Refreshing state...
data.hiera.aws_instance_size: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

aws_cloudwatch_enable = true
aws_instance_size = t2.large
aws_tags = {
  team = A
  tier = 1
}
java_opts = [
    -Xms512m, 
    -Xmx2g,
    -Dspring.profiles.active=live
]
```