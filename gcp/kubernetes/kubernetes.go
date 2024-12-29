package kubernetes

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func ReadConfig(ctx *pulumi.Context) (string, string, int, string, string, string, string) {
	// Read configuration values
	config := config.New(ctx, "")
	clusterName := config.Require("clusterName")
	location := config.Require("location")
	nodeCount := config.RequireInt("nodeCount")
	machineType := config.Require("machineType")
	serviceAccountId := config.Require("serviceAccountId")
	serviceAccountDisplayName := config.Require("serviceAccountDisplayName")
	nodePoolName := config.Require("nodePoolName")
	return clusterName, location, nodeCount, machineType, serviceAccountId, serviceAccountDisplayName, nodePoolName
}

// Creates a new service account
func CreateServiceAccount(ctx *pulumi.Context, serviceAccountId string, serviceAccountDisplayName string) (*serviceaccount.Account, error) {
	sa, err := serviceaccount.NewAccount(ctx, "default", &serviceaccount.AccountArgs{
		AccountId:   pulumi.String(serviceAccountId),
		DisplayName: pulumi.String(serviceAccountDisplayName),
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("serviceAccountEmail", sa.Email)
	return sa, nil
}

// Creates new GKE cluster
func CreateGKECluster(ctx *pulumi.Context, clusterName string, location string, nodeCount int, sa *serviceaccount.Account) (*container.Cluster, error) {
	primary, err := container.NewCluster(ctx, "primary", &container.ClusterArgs{
		Name:                  pulumi.String(clusterName),
		Location:              pulumi.String(location),
		RemoveDefaultNodePool: pulumi.Bool(true),
		InitialNodeCount:      pulumi.Int(nodeCount),
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("clusterName", primary.Name)
	return primary, nil
}

// Creates a new node pool for the GKE cluster
func CreateNodePool(ctx *pulumi.Context, nodePoolName, machineType string, nodeCount int, primary *container.Cluster, sa *serviceaccount.Account) (*container.NodePool, error) {
	nodePool, err := container.NewNodePool(ctx, "primary_preemptible_nodes", &container.NodePoolArgs{
		Name:      pulumi.String(nodePoolName),
		Cluster:   primary.ID(),
		NodeCount: pulumi.Int(nodeCount),
		NodeConfig: &container.NodePoolNodeConfigArgs{
			Preemptible:    pulumi.Bool(true),
			MachineType:    pulumi.String(machineType),
			ServiceAccount: sa.Email,
			OauthScopes: pulumi.StringArray{
				pulumi.String("https://www.googleapis.com/auth/cloud-platform"),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("nodePoolName", nodePool.Name)
	return nodePool, nil
}
