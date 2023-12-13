package kubernetes

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
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

	// Server is the address of the k3s server.
	Server pulumi.StringOutput
	// Provider is the provider for the k3s cluster.
	Provider *kubernetes.Provider
}

// K3seArgs are the arguments for creating a k3s cluster using k3se.
type K3seArgs struct {
	// Name is the name of the cluster.
	Name string `validate:"required,alphanum"`
}

// NewK3se manages a k3s cluster using k3se.
func NewK3se(ctx *pulumi.Context, name string, args *K3seArgs, opts ...pulumi.ResourceOption) (*K3se, error) {
	component := &K3se{
		Server:   pulumi.String("").ToStringOutput(),
		Provider: nil,
	}
	if err := ctx.RegisterComponentResource(ComponentTypeK3se, name, component, opts...); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(args); err != nil {
		return nil, err
	}

	outputs := pulumi.Map{
		"server":   component.Server,
		"provider": component.Provider,
	}

	// TODO: Remove this once we verified that this won't break anything.
	if args.Name == "charlie" {
		buffer := make([]byte, 16)
		if _, err := rand.Read(buffer); err != nil {
			return nil, err
		}
		updateTrigger := hex.EncodeToString(buffer)

		specPath := path.Join(K3seSpecDir, fmt.Sprintf("%s.yaml", args.Name))
		kubeConfigPath := path.Join(kubeConfigDir, fmt.Sprintf("%s.kubeconfig.yaml", args.Name))
		cmd, err := local.NewCommand(ctx, fmt.Sprintf("%s-r.command-k3se", name), &local.CommandArgs{
			Create: pulumi.Sprintf("k3se up %s -k %s", specPath, kubeConfigPath),
			Update: pulumi.Sprintf("k3se up %s -k %s", specPath, kubeConfigPath),
			Delete: pulumi.Sprintf("k3se down %s", specPath, kubeConfigPath),
			AssetPaths: pulumi.StringArray{
				pulumi.String(kubeConfigPath),
			},
			Environment: pulumi.StringMap{
				// HACK: This ensures that the command is run on every update.
				"UPDATE_TRIGGER": pulumi.String(updateTrigger),
			},
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}

		kubeConfigAsset := cmd.Assets.MapIndex(pulumi.String(kubeConfigPath))

		component.Server = kubeConfigAsset.ApplyT(func(item pulumi.AssetOrArchive) string {
			configString, err := readKubeConfigFromAsset(item)
			if err != nil {
				return err.Error()
			}
			// During preview the kubeconfig is not downloaded and hence empty.
			if configString == "" {
				return ""
			}

			config, err := clientcmd.Load([]byte(configString))
			if err != nil {
				return err.Error()
			}

			context := config.Contexts[config.CurrentContext]
			server := config.Clusters[context.Cluster].Server

			return server
		}).(pulumi.StringOutput)
		outputs["server"] = component.Server

		// TODO: Figure out how we can prevent the preview from always showing a replacement
		//       of this provider. One possible idea is to create the provider using the
		//       `k3se up -s <config>` command. Another idea is to create a native k3se
		//       provider that avoids the quirks of using a command.
		component.Provider, err = kubernetes.NewProvider(ctx, fmt.Sprintf("%s-p.kubernetes", name), &kubernetes.ProviderArgs{
			Kubeconfig: pulumi.ToSecret(kubeConfigAsset.ApplyT(func(item pulumi.AssetOrArchive) string {
				configString, err := readKubeConfigFromAsset(item)
				if err != nil {
					return err.Error()
				}

				return configString
			})).(pulumi.StringOutput),
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}
		outputs["provider"] = component.Provider
	}

	if err := ctx.RegisterResourceOutputs(component, outputs); err != nil {
		return nil, err
	}

	return component, nil
}

// readKubeConfigFromAsset parses a kubeconfig from an asset. During preview
// the kubeconfig is empty. This is expected behavior and needs to be handled
// accordingly by the consumer of this function.
func readKubeConfigFromAsset(item pulumi.AssetOrArchive) (string, error) {
	asset, ok := item.(pulumi.Asset)
	if !ok {
		return "", fmt.Errorf("expected pulumi.Asset: received %T", item)
	}

	// This avoids issues during the preview, because
	// the command is not executed during the preview
	// and thus the kubeconfig file is not created.
	if asset == nil || asset.Path() == "" {
		return "", nil
	}

	bytes, err := os.ReadFile(asset.Path())
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
