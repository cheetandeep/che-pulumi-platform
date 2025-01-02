package kubernetes

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Config struct {
	Env                       string
	App                       string
	Location                  string
	NodeCount                 int
	MachineType               string
	ServiceAccountId          string
	ServiceAccountDisplayName string
	NetworkingStackName       string
}

func ReadConfig(ctx *pulumi.Context) Config {
	// Read configuration values
	cfg := config.New(ctx, "")

	return Config{
		Env:                       cfg.Require("env"),
		App:                       cfg.Require("app"),
		Location:                  cfg.Require("location"),
		NodeCount:                 cfg.RequireInt("nodeCount"),
		MachineType:               cfg.Require("machineType"),
		ServiceAccountId:          cfg.Require("serviceAccountId"),
		ServiceAccountDisplayName: cfg.Require("serviceAccountDisplayName"),
		NetworkingStackName:       cfg.Require("networkingStackName"),
	}
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
func CreateGKECluster(ctx *pulumi.Context, cfg Config, location string, nodeCount int, sa *serviceaccount.Account) (*container.Cluster, error) {

	// Read stack name from configuration
	stackName := cfg.NetworkingStackName

	// Get outputs from the stack that created the network
	stackReference, err := pulumi.NewStackReference(ctx, stackName, nil)
	//stackReference, err := pulumi.NewStackReference(ctx, "siche2824/gcp-platform/networking", nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the necessary outputs from the referenced stack
	networkName := stackReference.GetOutput(pulumi.String("networkName")).ApplyT(func(id interface{}) string {
		if id == nil {
			ctx.Log.Error("networkName is nil", &pulumi.LogArgs{})
			return ""
		}
		return id.(string)
	}).(pulumi.StringOutput)

	subnetName := stackReference.GetOutput(pulumi.String("subnetName")).ApplyT(func(id interface{}) string {
		if id == nil {
			ctx.Log.Error("subnetName is nil", &pulumi.LogArgs{})
			return ""
		}
		return id.(string)
	}).(pulumi.StringOutput)

	secPodIpCidrRange := stackReference.GetOutput(pulumi.String("secPodRangeName")).ApplyT(func(id interface{}) string {
		if id == nil {
			ctx.Log.Error("secPodIpCidrRange is nil", &pulumi.LogArgs{})
			return ""
		}
		return id.(string)
	}).(pulumi.StringOutput)

	secSvcRangeName := stackReference.GetOutput(pulumi.String("secSvcRangeName")).ApplyT(func(id interface{}) string {
		if id == nil {
			ctx.Log.Error("secSvcRangeName is nil", &pulumi.LogArgs{})
			return ""
		}
		return id.(string)
	}).(pulumi.StringOutput)

	// Log the retrieved values for debugging
	networkName.ApplyT(func(id string) error {
		ctx.Log.Info(fmt.Sprintf("networkName: %s", id), &pulumi.LogArgs{})
		return nil
	})

	subnetName.ApplyT(func(id string) error {
		ctx.Log.Info(fmt.Sprintf("subnetName: %s", id), &pulumi.LogArgs{})
		return nil
	})

	secPodIpCidrRange.ApplyT(func(id string) error {
		ctx.Log.Info(fmt.Sprintf("secPodIpCidrRange: %s", id), &pulumi.LogArgs{})
		return nil
	})

	clusterName := cfg.App + "-" + cfg.Env + "-gke"
	primary, err := container.NewCluster(ctx, "primary", &container.ClusterArgs{
		Name:                  pulumi.String(clusterName),
		Location:              pulumi.String(location),
		RemoveDefaultNodePool: pulumi.Bool(true),
		InitialNodeCount:      pulumi.Int(nodeCount),
		Subnetwork:            subnetName.ApplyT(func(v interface{}) string { return v.(string) }).(pulumi.StringOutput),
		Network:               networkName.ApplyT(func(v interface{}) string { return v.(string) }).(pulumi.StringOutput),
		IpAllocationPolicy: &container.ClusterIpAllocationPolicyArgs{
			AdditionalPodRangesConfig: &container.ClusterIpAllocationPolicyAdditionalPodRangesConfigArgs{
				PodRangeNames: pulumi.StringArray{
					pulumi.String("string"),
				},
			},
			ClusterSecondaryRangeName: secPodIpCidrRange.ApplyT(func(v interface{}) string { return v.(string) }).(pulumi.StringOutput),

			ServicesSecondaryRangeName: secSvcRangeName.ApplyT(func(v interface{}) string { return v.(string) }).(pulumi.StringOutput),
			PodCidrOverprovisionConfig: &container.ClusterIpAllocationPolicyPodCidrOverprovisionConfigArgs{
				Disabled: pulumi.Bool(false),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("clusterName", primary.Name)
	return primary, nil
}

// Creates a new node pool for the GKE cluster
func CreateNodePool(ctx *pulumi.Context, cfg Config, machineType string, nodeCount int, primary *container.Cluster, sa *serviceaccount.Account) (*container.NodePool, error) {

	// Create a new node pool
	nodePoolName := cfg.App + "-" + cfg.Env + "-nodepool"
	nodePool, err := container.NewNodePool(ctx, "primary", &container.NodePoolArgs{
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
