---
layout: ""
page_title: "Provider: DataStax Astra - Serverless Cassandra DBaaS"
description: |-
  The Astra provider provides Terraform resources to interact with DataStax AstraDB and Astra Streaming, DataStax's cloud offerings based on Apache Cassandra, Apache Pulsar, and Kubernetes.
---

# Manage tenants

## Prerequisites

* [Astra](https://astra.datastax.com/) account
* Astra token for authentication, which looks like "AstraCS:xxyyzz..."
* [Terraform](https://www.terraform.io/) version 1.0 or higher

## Create tenant

1. Create a new directory, and within it, create a `main.tf` file.
```
mkdir terraform
touch main.tf
```
2. Modify main.tf with your favorite text editor to add your tenant
```terraform
terraform {
  required_providers {
    astra = {
      source = "datastax/astra"
      version = "2.1.15"
    }
  }
}

variable "token" {}

provider "astra" {
  // This can also be set via ASTRA_API_TOKEN environment variable.
  token = var.token
}

resource "astra_database" "example" {
  name           = "terraform-db"
  keyspace       = "ks1"
  cloud_provider = "aws"
  regions        = ["us-east-1"]
}
```
3. Initialize terraform
```
terraform init
...
Terraform has been successfully initialized!
```
4. Apply your `main.tf` file.
```
terraform apply
```
3. You will be prompted to input your Astra token (AstraCS:xxyyzz...). Input your token.
4. You will be prompted to create the listed tenant resources. Type yes to continue. Terraform will begin creating the tenant.
```
astra_database.example: Creation complete after 2m0s [id=5a31c262-14ab-4569-bba7-d1fb38349413]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```
5. If you modify values in your `main.tf` file, use `terraform apply` to update your deployment.

## Additional Info

To report bugs or feature requests for the provider [file an issue on github](https://github.com/datastax/terraform-provider-astra/issues).

For help, contact [DataStax Support](https://support.datastax.com/).
