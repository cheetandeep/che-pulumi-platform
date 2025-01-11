package config

import (
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
	HostProject               string
	ServiceProject            string
	PrimaryIpCidrRange        string
	SecPodIpCidrRange         string
	SecSvcIpCidrRange         string
	Region                    string
}

func ReadConfig(ctx *pulumi.Context) Config {
	cfg := config.New(ctx, "")
	return Config{
		Env:                       cfg.Get("env"),
		App:                       cfg.Get("app"),
		Location:                  cfg.Get("location"),
		NodeCount:                 cfg.GetInt("nodeCount"),
		MachineType:               cfg.Get("machineType"),
		ServiceAccountId:          cfg.Get("serviceAccountId"),
		ServiceAccountDisplayName: cfg.Get("serviceAccountDisplayName"),
		NetworkingStackName:       cfg.Get("networkingStackName"),
		HostProject:               cfg.Get("hostProject"),
		ServiceProject:            cfg.Get("serviceProject"),
		PrimaryIpCidrRange:        cfg.Get("primaryIpCidrRange"),
		SecPodIpCidrRange:         cfg.Get("secondarypodIpCidrRange"),
		SecSvcIpCidrRange:         cfg.Get("secondarySvcIpCidrRange"),
		Region:                    cfg.Get("region"),
	}
}
