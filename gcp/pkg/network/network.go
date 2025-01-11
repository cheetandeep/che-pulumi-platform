package network

import (
	"gcp-platform/internal/config"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateNetwork(ctx *pulumi.Context, cfg config.Config) (*compute.Network, error) {
	networkName := "shared-" + cfg.Env + "-network"
	net, err := compute.NewNetwork(ctx, networkName, &compute.NetworkArgs{
		Name:                  pulumi.String(networkName),
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("network", net.ID())
	ctx.Export("networkName", pulumi.String(networkName))
	return net, nil
}

func CreateSubnetwork(ctx *pulumi.Context, cfg config.Config, networkID pulumi.IDOutput) (*compute.Subnetwork, error) {
	subnetName := "shared-" + cfg.Env + "-subnet"
	subnets, err := compute.NewSubnetwork(ctx, subnetName, &compute.SubnetworkArgs{
		Name:        pulumi.String(subnetName),
		IpCidrRange: pulumi.String(cfg.PrimaryIpCidrRange),
		Region:      pulumi.String(cfg.Region),
		Network:     networkID,
		SecondaryIpRanges: compute.SubnetworkSecondaryIpRangeArray{
			&compute.SubnetworkSecondaryIpRangeArgs{
				RangeName:   pulumi.String("sec-pod-range"),
				IpCidrRange: pulumi.String(cfg.SecPodIpCidrRange),
			},
			&compute.SubnetworkSecondaryIpRangeArgs{
				RangeName:   pulumi.String("sec-svc-range"),
				IpCidrRange: pulumi.String(cfg.SecSvcIpCidrRange),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("subnet", subnets.ID())
	ctx.Export("subnetName", pulumi.String(subnetName))
	ctx.Export("primaryIpCidrRange", pulumi.String(cfg.PrimaryIpCidrRange))
	ctx.Export("secPodRangeName", pulumi.String("sec-pod-range"))
	ctx.Export("secSvcRangeName", pulumi.String("sec-svc-range"))
	return subnets, nil
}

func SharedVpc(ctx *pulumi.Context, cfg config.Config) (*compute.SharedVPCServiceProject, error) {

	// A host project provides network resources to associated service projects.
	host, err := compute.NewSharedVPCHostProject(ctx, "host", &compute.SharedVPCHostProjectArgs{
		Project: pulumi.String(cfg.HostProject),
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("HostProject", host.ID())
	return nil, err

	// // TODO: Uncomment this block to create a service project.
	// A service project gains access to network resources provided by its
	// associated host project.
	// service, err := compute.NewSharedVPCServiceProject(ctx, "service1", &compute.SharedVPCServiceProjectArgs{
	// 	HostProject:    host.Project,
	// 	ServiceProject: pulumi.String(cfg.ServiceProject),
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// ctx.Export("ServiceProject", service.ID())
	// return service, err
}

func NatRouter(ctx *pulumi.Context, cfg config.Config, networkID pulumi.IDOutput) (*compute.Router, error) {

	cloudRouterName := "shared-" + cfg.Env + "-cloud-nat-router"
	router, err := compute.NewRouter(ctx, "router", &compute.RouterArgs{
		Name:    pulumi.String(cloudRouterName),
		Region:  pulumi.String(cfg.Region),
		Network: networkID,
		Bgp: &compute.RouterBgpArgs{
			Asn: pulumi.Int(64514),
		},
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("router", router.Name)

	natRouterName := "shared-" + cfg.Env + "-nat-router"
	natRouter, err := compute.NewRouterNat(ctx, "nat", &compute.RouterNatArgs{
		Name:                          pulumi.String(natRouterName),
		Router:                        router.Name,
		Region:                        router.Region,
		NatIpAllocateOption:           pulumi.String("AUTO_ONLY"),
		SourceSubnetworkIpRangesToNat: pulumi.String("ALL_SUBNETWORKS_ALL_IP_RANGES"),
		LogConfig: &compute.RouterNatLogConfigArgs{
			Enable: pulumi.Bool(true),
			Filter: pulumi.String("ERRORS_ONLY"),
		},
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("nat", natRouter.ID())
	return router, nil
}
