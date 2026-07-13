package hetznerrobot

import "testing"

// TestProvider validates the full provider schema (provider config, every
// resource and data source) for internal consistency.
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider internal validation failed: %s", err)
	}
}

func TestProvider_registrations(t *testing.T) {
	p := Provider()

	wantResources := []string{
		"hetzner-robot_boot",
		"hetzner-robot_firewall",
		"hetzner-robot_ssh_key",
		"hetzner-robot_vswitch",
		"hetzner-robot_reset",
		"hetzner-robot_ip",
		"hetzner-robot_ip_mac",
		"hetzner-robot_subnet",
		"hetzner-robot_subnet_mac",
		"hetzner-robot_wol",
		"hetzner-robot_storagebox",
		"hetzner-robot_storagebox_snapshot",
		"hetzner-robot_storagebox_subaccount",
	}
	for _, name := range wantResources {
		if _, ok := p.ResourcesMap[name]; !ok {
			t.Errorf("expected resource %q to be registered", name)
		}
	}

	wantData := []string{
		"hetzner-robot_boot",
		"hetzner-robot_server",
		"hetzner-robot_servers",
		"hetzner-robot_ssh_key",
		"hetzner-robot_vswitch",
		"hetzner-robot_ip",
		"hetzner-robot_subnet",
		"hetzner-robot_traffic",
		"hetzner-robot_storagebox",
		"hetzner-robot_storageboxes",
	}
	for _, name := range wantData {
		if _, ok := p.DataSourcesMap[name]; !ok {
			t.Errorf("expected data source %q to be registered", name)
		}
	}
}
