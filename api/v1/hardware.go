package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HardwareChassisDimensions describes the dimensions of a chassis.
type HardwareChassisDimensions struct {
	DepthMm  int `json:"depthMm,omitempty" yaml:"depthMm"`
	HeightMm int `json:"heightMm,omitempty" yaml:"heightMm"`
	WidthMm  int `json:"widthMm,omitempty" yaml:"widthMm"`
}

// HardwareChassis describes a chassis.
type HardwareChassis struct {
	Dimensions   HardwareChassisDimensions `json:"dimensions,omitempty" yaml:"dimensions"`
	Manufacturer string                    `json:"manufacturer,omitempty" yaml:"manufacturer"`
	Model        string                    `json:"model,omitempty" yaml:"model"`
	Type         string                    `json:"type,omitempty" yaml:"type"`
}

// HardwareCPU describes a CPU.
type HardwareCPU struct {
	Architecture string `json:"architecture,omitempty" yaml:"architecture"`
	Cores        int    `json:"cores,omitempty" yaml:"cores"`
	ClockMHz     int    `json:"clockMHz,omitempty" yaml:"clockMHz"`
	Manufacturer string `json:"manufacturer,omitempty" yaml:"manufacturer"`
	Model        string `json:"model,omitempty" yaml:"model"`
	Threads      int    `json:"threads,omitempty" yaml:"threads"`
}

// HardwareMemory describes a memory module.
type HardwareMemory struct {
	CapacityGiB  int    `json:"capacityGiB,omitempty" yaml:"capacityGiB"`
	ClockMHz     int    `json:"clockMHz,omitempty" yaml:"clockMHz"`
	Manufacturer string `json:"manufacturer,omitempty" yaml:"manufacturer"`
	Model        string `json:"model,omitempty" yaml:"model"`
}

// HardwareDisk describes a disk.
type HardwareDisk struct {
	CapacityGiB  int     `json:"capacityGiB,omitempty" yaml:"capacityGiB"`
	FormFactorIn float64 `json:"formFactorIn,omitempty" yaml:"formFactorIn"`
	Interface    string  `json:"interface,omitempty" yaml:"interface"`
	Manufacturer string  `json:"manufacturer,omitempty" yaml:"manufacturer"`
	Model        string  `json:"model,omitempty" yaml:"model,omitempty"`
	Name         string  `json:"name,omitempty" yaml:"name"`
	Type         string  `json:"type,omitempty" yaml:"type"`
}

// HardwareInterface describes a network interface.
type HardwareInterface struct {
	MAC       string `json:"mac,omitempty" yaml:"mac"`
	Name      string `json:"name,omitempty" yaml:"name"`
	SpeedMbps int    `json:"speedMbps,omitempty" yaml:"speedMbps"`
	Type      string `json:"type,omitempty" yaml:"type"`
}

// HardwareSpec defines the configuration of a physical server.
type HardwareSpec struct {
	Hostname   string              `json:"hostname,omitempty" yaml:"hostname"`
	Chassis    HardwareChassis     `json:"chassis,omitempty" yaml:"chassis"`
	CPUs       []HardwareCPU       `json:"cpus,omitempty" yaml:"cpus"`
	Memory     []HardwareMemory    `json:"memory,omitempty" yaml:"memory"`
	Disks      []HardwareDisk      `json:"disks,omitempty" yaml:"disks"`
	Interfaces []HardwareInterface `json:"interfaces,omitempty" yaml:"interfaces"`
}

// Hardware describes a physical server.
type Hardware struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec HardwareSpec `json:"spec,omitempty" yaml:"spec"`
}
