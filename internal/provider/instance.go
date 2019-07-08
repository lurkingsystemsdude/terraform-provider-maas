package provider

import (
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
)

// MACAddress is used by the maas_instance resource
type MACAddress maas.MACAddress

// Instance models the maas_instance Terraform resource
type Instance struct {
	Architecture           string        `optional:"true" forcenew:"true"`
	BootType               string        `optional:"true" forcenew:"true"`
	CPUCount               int           `optional:"true" forcenew:"true"`
	DisableIPv4            bool          `optional:"true" name:"disable_ipv4"`
	DistroSeries           string        `optional:"true" forcenew:"true"`
	Hostname               string        `optional:"true" forcenew:"true"`
	DeployHostname         string        `optional:"true" forcenew:"true"`
	DeployTags             []string      `optional:"true" forcenew:"true"`
	Tags                   []string      `optional:"true" forcenew:"true"`
	ReleaseErase           bool          `optional:"true" forcenew:"true" default:"false"`
	ReleaseEraseSecure     bool          `optional:"true" forcenew:"true" default:"false"`
	ReleaseEraseQuick      bool          `optional:"true" forcenew:"true" default:"false"`
	IPAddresses            []string      `optional:"true" forcenew:"true"`
	MACAddressSet          []MACAddress  `optional:"true" forcenew:"true"`
	Memory                 int           `optional:"true" forcenew:"true"`
	Netboot                bool          `optional:"true" forcenew:"true"`
	OSystem                string        `optional:"true" forcenew:"true"`
	Owner                  string        `optional:"true" forcenew:"true"`
	PhysicalBlockDeviceSet []BlockDevice `optional:"true" forcenew:"true" name:"physicalblockdevice_set"`
	PowerState             string        `optional:"true"`
	PowerType              string        `optional:"true"`
	PXEMac                 []MACAddress  `optional:"true" type:"Set"`
	ResourceURI            string        `optional:"true" forcenew:"true"`
	Routers                []string      `optional:"true"`
	Status                 int           `optional:"true"`
	Storage                int           `optional:"true"`
	SwapSize               int           `optional:"true"`
	SystemID               string        `optional:"true" forcenew:"true"`
	TagNames               []string      `optional:"true"`
	Zone                   []Zone        `optional:"true" type:"Set"`
	UserData               string        `optional:"true" forcenew:"true" statefunc:"true"`
	HWEKernel              string        `optional:"true" forcenew:"true"`
	Comment                string        `optional:"true"`
	Lock                   bool          `optional:"true" default:"false"`
}

// NewInstance creates a new instance from the value of the Terraform resource
func NewInstance(resource *schema.ResourceData) Instance {
	var instance Instance
	st := reflect.TypeOf(instance)
	sv := reflect.ValueOf(instance)

	for i := 0; i < st.NumField(); i++ {
		// Get the name of the schema field
		key := st.Field(i).Name
		if tag, ok := st.Field(i).Tag.Lookup("name"); ok {
			if tag == "-" {
				continue
			}
			key = tag
		}
		key = strings.ToLower(key)

		// Set the value if one exists
		if schemaVal, ok := resource.GetOk(key); ok {
			field := sv.FieldByName(key)
			reflectVal := reflect.ValueOf(schemaVal)
			field.Set(reflectVal)
		}
	}
	return instance
}

// FromMachine updates the instance to reflect the state of a Machine
func (i Instance) FromMachine(m maas.Machine) Instance {
	return i
}

// UpdateState updates the Terraform state to match the Instance state
func (i Instance) UpdateState(resource *schema.ResourceData) {
	st := reflect.TypeOf(i)
	sv := reflect.ValueOf(i)

	for i := 0; i < st.NumField(); i++ {
		// Get the name of the schema field
		key := st.Field(i).Name
		if tag, ok := st.Field(i).Tag.Lookup("name"); ok {
			if tag == "-" {
				continue
			}
			key = tag
		}
		key = strings.ToLower(key)

		// Set the value on the resource if it is not a zero value
		field := sv.FieldByName(key)
		if field.IsValid() {
			resource.Set(key, field.Interface())
		}
	}
}

// AllocateParams creates parameters based on the current value of the Instance
func (i Instance) AllocateParams() maas.MachinesAllocateParams {
	var params maas.MachinesAllocateParams
	if i.SystemID != "" {
		params.SystemID = i.SystemID
		return params
	}
	return params
}

// DeployParams creates parameters based on the current value of the Instance
func (i Instance) DeployParams() maas.MachineDeployParams {
	var params maas.MachineDeployParams
	return params
}

func (i Instance) UserDataStateFunc(v interface{}) string {
	switch v.(type) {
	case string:
		hash := sha1.Sum([]byte(v.(string)))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}
