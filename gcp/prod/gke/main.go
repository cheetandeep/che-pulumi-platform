package main

import (
	"fmt"
	"gcp-platform/pkg/kubernetes"

	"gcp-platform/internal/config"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read configuration values
		cfg := config.ReadConfig(ctx)

		// Create a new service account
		sa, err := kubernetes.CreateServiceAccount(ctx, cfg.ServiceAccountId, cfg.ServiceAccountDisplayName)
		if err != nil {
			return fmt.Errorf("failed to create service account: %w", err)
		}
		ctx.Log.Info("Service account created", &pulumi.LogArgs{})

		// Create a new GKE cluster
		cluster, err := kubernetes.CreateGKECluster(ctx, cfg, cfg.Location, cfg.NodeCount, sa)
		if err != nil {
			return fmt.Errorf("failed to create GKE cluster: %w", err)
		}
		ctx.Log.Info("GKE cluster created", &pulumi.LogArgs{})

		// Create a new node pool
		_, err = kubernetes.CreateNodePool(ctx, cfg, cfg.MachineType, cfg.NodeCount, cluster, sa)
		if err != nil {
			return fmt.Errorf("failed to create node pool: %w", err)
		}
		ctx.Log.Info("Node pool created", &pulumi.LogArgs{})

		return nil
	})
}
