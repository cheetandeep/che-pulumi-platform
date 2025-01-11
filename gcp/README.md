# GCP Platform with Pulumi

This repository contains a Pulumi project to deploy a various Goole Cloud Platform resources. The project is structured to support multiple environments (e.g., dev, prod) and uses modular code to manage different resources.

## Directory Structure

```
pulumi/
  go-examples/
│   ├── dev/
│   │   ├── networking/
│   │   │   └── main.go  (Deploys networking resources)
│   │   ├── Pulumi.networking.yaml  (Configuration for dev network)
│   │   └── gke/
│   │   │   └── main.go (Deploys GKE cluster and node pool)
│   │   ├── Pulumi.gke.yaml  (Configuration for dev gke)
│   ├── prod/
│   │   ├── networking/
│   │   │   └── main.go (Deploys networking resources)
│   │   ├── Pulumi.networking.yaml  (Configuration for prod network)
│   │   ├── gke/
│   │   │   └── main.go (Deploys GKE resources, likely calling modules)
│   │   ├── Pulumi.gke.yaml  (Configuration for prod GKE)
    go.mod
    go.sum
```


## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
- [Go](https://golang.org/doc/install)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)

## Setup

1. **Clone the repository**:
   ```sh
   git clone https://github.com/cheetandeep/che-pulumi-platform

2. **Install dependencies**:
```
go mod tidy
```

3. Configure Pulumi
   * Ensure you are logged in to Pulumi:
```
pulumi login
``` 
   * Configure your Google Cloud credentials:
  
```
gcloud auth login
gcloud auth application-default login
```

## Deploying GKE Resources

### Development Environment

1. **Navigate to the dev GKE directory:**
```
cd dev/gke
```

2. **Create a new Pulumi stack::**
```
pulumi stack init dev
```

3. **Configure the stack:**
```
pulumi config set gcp:project YOUR_GCP_PROJECT_ID
pulumi config set gcp:region YOUR_GCP_REGION
```

4. **Deploy the stack:**
```
pulumi up
```

### Production Environment

1. **Navigate to the prod GKE directory:**
```
cd prod/gke
```

2. **Create a new Pulumi stack::**
```
pulumi stack init prod
```

3. **Configure the stack:**
```
pulumi config set gcp:project YOUR_GCP_PROJECT_ID
pulumi config set gcp:region YOUR_GCP_REGION
```

4. **Deploy the stack:**
```
pulumi up
```

## Project Structure
* `go.mod` and `go.sum`: Manage dependencies for the entire project.
* `dev/` and `prod/` Directories: Contain environment-specific directories (`gke`) and their respective configuration files (`Pulumi.dev.yaml`, `Pulumi.prod.yaml`).
* `pkg/` Directories: Contain reusable modules for creating various resources.
* `internal/` This directory contains packages that are intended to be used only within your project. These packages are not accessible to other projects.
  * `config/`: Contains configuration-related code.
  
    * `config.go`: Code for reading and managing configuration values.
  

### Example Files
`Pulumi.yaml` for `dev/gke`

```
name: gcp-platform-dev-gke
runtime: go
description: A Go program to deploy a GKE cluster on Google Cloud for the dev environment
```

`Pulumi.dev.yaml` for `dev/gke`
```
config:
  gcp-platform:clusterName: "dev-gke-cluster"
  gcp-platform:location: "us-central1"
  gcp-platform:nodeCount: 1
  gcp-platform:machineType: "e2-medium"
  gcp-platform:serviceAccountId: "dev-service-account-id"
  gcp-platform:serviceAccountDisplayName: "Dev Service Account"
  gcp-platform:nodePoolName: "dev-node-pool"
```
## Cleaning Up
To clean up the resources created by Pulumi, run the following command in the respective environment directory:

```
pulumi destroy
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Summary

This `README.md` provides an overview of the project structure, setup instructions, and deployment steps for the GKE resources. It also includes example configuration files and code snippets to help users understand how to manage and deploy the infrastructure using Pulumi.