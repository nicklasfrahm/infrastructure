package kubernetes

import (
	"fmt"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// ComponentTypeK3se is the ID of the component type.
	ComponentTypeK3se = "nicklasfrahm:kubernetes:K3se"
	// K3seSpecDir is the directory where k3se configurations are stored.
	K3seSpecDir = "deploy/k3se"

	// kubeConfigDir is the path where the kubeconfig files will be stored.
	kubeConfigDir = "./"
)

// K3se manages a k3s cluster using k3se.
type K3se struct {
	pulumi.ResourceState

	// Server is the address of the Kubernetes API server.
	Server pulumi.StringOutput `pulumi:"server"`
}

// K3seArgs are the arguments for creating a k3s cluster using k3se.
type K3seArgs struct {
	// Name is the name of the cluster.
	Name string `validate:"required,alphanum"`
}

// NewK3se configures CNAME DNS records for the given hostname.
func NewK3se(ctx *pulumi.Context, name string, args *K3seArgs, opts ...pulumi.ResourceOption) (*K3se, error) {
	component := &K3se{}
	if err := ctx.RegisterComponentResource(ComponentTypeK3se, name, component, opts...); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(args); err != nil {
		return nil, err
	}
	outputs := pulumi.Map{}

	// TODO: Remove this once we verified that this won't break anything.
	if args.Name == "charlie" {
		specPath := path.Join(K3seSpecDir, fmt.Sprintf("%s.yaml", args.Name))
		kubeConfigPath := path.Join(kubeConfigDir, fmt.Sprintf("%s.kubeconfig.yaml", args.Name))
		cmd, err := local.NewCommand(ctx, fmt.Sprintf("%s-r.command-k3se", name), &local.CommandArgs{
			Create: pulumi.Sprintf("k3se up %s -k %s", specPath, kubeConfigPath),
			Delete: pulumi.Sprintf("k3se down %s", specPath, kubeConfigPath),
			AssetPaths: pulumi.StringArray{
				pulumi.String(kubeConfigPath),
			},
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}

		component.Server = cmd.Assets.ApplyT(func(item map[string]pulumi.AssetOrArchive) string {
			// Ignore empty assets during the preview.
			if item == nil {
				return ""
			}

			asset, ok := item[kubeConfigPath].(pulumi.Asset)
			if !ok {
				return fmt.Errorf("expected pulumi.Asset: received %T", item).Error()
			}

			config, err := clientcmd.LoadFromFile(asset.Path())
			if err != nil {
				return err.Error()
			}

			return config.Clusters[config.CurrentContext].Server
		}).(pulumi.StringOutput)

		outputs["server"] = component.Server
	}

	if err := ctx.RegisterResourceOutputs(component, outputs); err != nil {
		return nil, err
	}

	return component, nil
}
