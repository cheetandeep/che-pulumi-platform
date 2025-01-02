package main

import (
	"fmt"
	"gcp-platform/network"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := network.ReadConfig(ctx)

		// Create Virtual Private Cloud (VPC) network
		net, err := network.CreateNetwork(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create network: %w", err)
		}
		ctx.Log.Info("Network created", &pulumi.LogArgs{})

		// Create Subnetwork and Secondary IP ranges
		_, err = network.CreateSubnetwork(ctx, cfg, net.ID())
		if err != nil {
			return fmt.Errorf("failed to create subnetwork: %w", err)
		}
		ctx.Log.Info("Subnetwork created", &pulumi.LogArgs{})

		// Create Shared VPC
		_, err = network.SharedVpc(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create shared VPC: %w", err)
		}
		ctx.Log.Info("Shared VPC created", &pulumi.LogArgs{})

		// Create NAT router
		_, err = network.NatRouter(ctx, cfg, net.ID())
		if err != nil {
			return fmt.Errorf("failed to create NAT router: %w", err)
		}
		ctx.Log.Info("NAT router created", &pulumi.LogArgs{})

		return nil
	})
}
