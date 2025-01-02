# Pulumi GCP Networking Setup

This repository contains a Pulumi setup for creating and managing Google Cloud Platform (GCP) networking resources, including Virtual Private Cloud (VPC) networks, subnetworks, shared VPCs, and NAT routers.

## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
- [Go](https://golang.org/doc/install)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- A GCP project with the necessary permissions

## Configuration

Before running the Pulumi program, you need to configure the necessary values. These values are read from the Pulumi configuration file.

### Required Configuration Values

- `env`: The environment (e.g., `dev`, `prod`).
- `hostProject`: The GCP project ID for the host project.
- `serviceProject`: The GCP project ID for the service project.
- `primaryIpCidrRange`: The primary IP CIDR range for the subnetwork.
- `secondarypodIpCidrRange`: The secondary IP CIDR range for pods.
- `secondarySvcIpCidrRange`: The secondary IP CIDR range for services.
- `region`: The GCP region where the resources will be created.

### Setting Configuration Values

You can set the configuration values using the Pulumi CLI:

```sh
pulumi config set env <your-env>
pulumi config set hostProject <your-host-project-id>
pulumi config set serviceProject <your-service-project-id>
pulumi config set primaryIpCidrRange <your-primary-ip-cidr-range>
pulumi config set secondarypodIpCidrRange <your-secondary-pod-ip-cidr-range>
pulumi config set secondarySvcIpCidrRange <your-secondary-svc-ip-cidr-range>
pulumi config set region <your-region>