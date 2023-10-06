package main

import (
	"fmt"

	"github.com/spf13/cobra"
	tinkv1 "github.com/tinkerbell/tink/api/v1alpha2"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LabelCores        = "baremetal.nicklasfrahm.dev/cpus"
	LabelMemory       = "baremetal.nicklasfrahm.dev/memory"
	LabelManufacturer = "baremetal.nicklasfrahm.dev/manufacturer"
	LabelModel        = "baremetal.nicklasfrahm.dev/model"
)

var cpus int
var memoryGB int
var manufacturer string
var model string

var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Enroll a new server",
	Long: `Enroll a new piece of hardware.

This command will create a configuration
file for the new server in the repository
and generate a random passphrase for the
encrypted disk.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		vendorData := fmt.Sprintf(`{
  "cpus": %d,
  "memory": %d,
  "manufacturer": "%s",
  "model": "%s"
}`, cpus, memoryGB, manufacturer, model)

		hw := tinkv1.Hardware{
			ObjectMeta: metav1.ObjectMeta{
				Name: serverName,
				Labels: map[string]string{
					LabelCores:        fmt.Sprint(cpus),
					LabelMemory:       fmt.Sprint(memoryGB),
					LabelManufacturer: manufacturer,
					LabelModel:        model,
				},
			},
			Spec: tinkv1.HardwareSpec{
				NetworkInterfaces: tinkv1.NetworkInterfaces{
					"00:e0:4c:88:00:f1": tinkv1.NetworkInterface{},
				},
				Instance: &tinkv1.Instance{
					Vendordata: &vendorData,
				},
			},
		}

		yamlBytes, err := yaml.Marshal(hw)
		if err != nil {
			return err
		}

		fmt.Printf("Enrolling new server:\n%s\n", string(yamlBytes))

		return nil
	},
}

func init() {
	enrollCmd.Flags().IntVarP(&cpus, "cpus", "c", 4, "number of CPU cores")
	enrollCmd.MarkFlagRequired("cpus")
	enrollCmd.Flags().IntVarP(&memoryGB, "memory", "m", 16, "amount of memory in GB")
	enrollCmd.MarkFlagRequired("memory")
	enrollCmd.Flags().StringVarP(&manufacturer, "manufacturer", "M", "raspberrypi", "manufacturer of the server")
	enrollCmd.MarkFlagRequired("manufacturer")
	enrollCmd.Flags().StringVarP(&model, "model", "B", "4b", "model of the server")
	enrollCmd.MarkFlagRequired("model")

	rootCmd.AddCommand(enrollCmd)
}
