package main

import (
	"fmt"
	"gcp-platform/kubernetes"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read configuration values
		clusterName, location, nodeCount, machineType, serviceAccountId, serviceAccountDisplayName, nodePoolName := kubernetes.ReadConfig(ctx)

		// Create a new service account
		sa, err := kubernetes.CreateServiceAccount(ctx, serviceAccountId, serviceAccountDisplayName)
		if err != nil {
			return fmt.Errorf("failed to create service account: %w", err)
		}
		ctx.Log.Info("Service account created", &pulumi.LogArgs{})

		// Create a new GKE cluster
		primary, err := kubernetes.CreateGKECluster(ctx, clusterName, location, nodeCount, sa)
		if err != nil {
			return fmt.Errorf("failed to create GKE cluster: %w", err)
		}
		ctx.Log.Info("GKE cluster created", &pulumi.LogArgs{})

		// Create a new node pool
		_, err = kubernetes.CreateNodePool(ctx, nodePoolName, machineType, nodeCount, primary, sa)
		if err != nil {
			return fmt.Errorf("failed to create node pool: %w", err)
		}
		ctx.Log.Info("Node pool created", &pulumi.LogArgs{})

		return nil
	})
}
