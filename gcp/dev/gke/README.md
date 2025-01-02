# Pulumi GCP Kubernetes Setup

This repository contains a Pulumi setup for creating and managing Google Cloud Platform (GCP) resources, including a GKE cluster and its associated node pool. The project is structured to support multiple environments (e.g., dev, prod) and uses modular code to manage different resources.

The setup demonstrates the use of `stack references` to share resources between different Pulumi stacks


## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
- [Go](https://golang.org/doc/install)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- A GCP project with the necessary permissions

## Configuration

Before running the Pulumi program, you need to configure the necessary values. These values are read from the Pulumi configuration file.

### Required Configuration Values

- `env`: The environment (e.g., `dev`, `prod`).
- `app`: The application name.
- `location`: The GCP region where the resources will be created.
- `nodeCount`: The number of nodes in the GKE cluster.
- `machineType`: The machine type for the GKE nodes.
- `serviceAccountId`: The ID of the service account to be used.
- `serviceAccountDisplayName`: The display name of the service account.
- `networkingStackName`: The name of the stack that creates the networking resources.

### Setting Configuration Values

You can set the configuration values using the Pulumi CLI:

```sh
pulumi config set env <your-env>
pulumi config set app <your-app>
pulumi config set location <your-location>
pulumi config set nodeCount <your-node-count>
pulumi config set machineType <your-machine-type>
pulumi config set serviceAccountId <your-service-account-id>
pulumi config set serviceAccountDisplayName <your-service-account-display-name>
pulumi config set networkingStackName <your-networking-stack-name>
```


## Running the Pulumi Program
1. **Clone the repository**:
```
   git clone https://github.com/cheetandeep/che-pulumi-platform
   cd gcp/dev/gke

```

2. **Install dependencies**:
```
go mod tidy
```

3. **Configure Pulumi**
   * Ensure you are logged in to Pulumi:
```
pulumi login
``` 
   * Configure your Google Cloud credentials:
  
```
gcloud auth login
gcloud auth application-default login
```


## Project Structure
* `main.go`: The main entry point for the Pulumi program
* `kubernetes/`: Contains functions for creating and managing GCP Kubernetes resources
  * `kubernetes.go`: Contains functions for reading configuration, creating service accounts, GKE clusters, and node pools

### Stack Reference Concept

This project demonstrates the use of Pulumi's stack references to share resources between different stacks. The `networkingStackName` configuration value specifies the name of the stack that creates the networking resources. The `CreateNodePool` function uses a stack reference to retrieve the necessary outputs from the networking stack.

#### Example of Using Stack References
In the `CreateNodePool` function, a stack reference is created to access the outputs from the networking stack:


```
// Read stack name from configuration
stackName := cfg.NetworkingStackName

// Get outputs from the stack that created the network
stackReference, err := pulumi.NewStackReference(ctx, stackName, nil)
if err != nil {
    return nil, err
}

// Retrieve the necessary outputs from the referenced stack
networkID := stackReference.GetOutput(pulumi.String("network")).ApplyT(func(id interface{}) string {
    if id == nil {
        ctx.Log.Error("networkID is nil", &pulumi.LogArgs{})
        return ""
    }
    return id.(string)
}).(pulumi.StringOutput)

subnetID := stackReference.GetOutput(pulumi.String("subnet")).ApplyT(func(id interface{}) string {
    if id == nil {
        ctx.Log.Error("subnetID is nil", &pulumi.LogArgs{})
        return ""
    }
    return id.(string)
}).(pulumi.StringOutput)
```

This allows the `CreateNodePool` function to use the `networkID` and `subnet` outputs from the networking stack to create the node pool.

## License
This project is licensed under the MIT License. See the LICENSE file for details.