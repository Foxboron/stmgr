package provision

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"

	"github.com/u-root/u-root/pkg/efivarfs"
)

const defaultFilePerm fs.FileMode = 0o600

// HostCfgSimplified mirrors what stboots host configuration
// looks like, but does not implement their custom types and
// the extensive parsing options. It's enough for the cfgtool
// to create a configuration, but might get streamlined with stboot
// after the final layout is set.
type HostCfgSimplified struct {
	Version           int               `json:"version"`
	IPAddrMode        *string           `json:"network_mode"`
	HostIP            *string           `json:"host_ip"`
	DefaultGateway    *string           `json:"gateway"`
	DNSServer         *string           `json:"dns"`
	NetworkInterface  *string           `json:"network_interface"`
	ProvisioningURLs  []string          `json:"provisioning_urls"`
	ID                *string           `json:"identity"`
	Auth              *string           `json:"authentication"`
	Timestamp         *int64            `json:"timestamp"`
	NetworkInterfaces []string          `json:"network_interfaces"`
	BondingMode       *string           `json:"bonding_mode"`
	BondName          *string           `json:"bond_name"`
	Custom            map[string]string `json:"custom,omitempty"`
}

// MarshalCfg takes a HostCfgSimplified struct and depending
// on the efi bool either writes it to disk as "host_configuration.json"
// in the current directory or into the efivarfs.
func MarshalCfg(cfg *HostCfgSimplified, efi bool) error {
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if efi {
		name := "STHostConfig-f401f2c1-b005-4be0-8cee-f2e5945bcbe7"
		attrs := efivarfs.AttributeBootserviceAccess | efivarfs.AttributeRuntimeAccess | efivarfs.AttributeNonVolatile

		e, err := efivarfs.New()
		if err != nil {
			return err
		}

		return efivarfs.SimpleWriteVariable(e, name, attrs, bytes.NewBuffer(jsonBytes))
	}

	return os.WriteFile("host_configuration.json", jsonBytes, defaultFilePerm)
}
