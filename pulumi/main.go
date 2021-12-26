package main

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cluster, err := digitalocean.NewKubernetesCluster(ctx, "roo-k8s", &digitalocean.KubernetesClusterArgs{
			NodePool: digitalocean.KubernetesClusterNodePoolArgs{
				Name:      pulumi.String("roo-k8s-pool"),
				Size:      pulumi.String("s-2vcpu-4gb"),
				NodeCount: pulumi.Int(3),
			},
			Region:  pulumi.String("lon1"),
			Version: pulumi.String("1.21.5-do.0"),
		})
		if err != nil {
			return err
		}

		ctx.Export("kubeconfig", cluster.KubeConfigs.Index(pulumi.Int(0)).RawConfig())

		return nil
	})
}
