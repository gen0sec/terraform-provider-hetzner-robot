package hetznerrobot

import (
	"reflect"
	"testing"
)

func TestServerOrderAddons(t *testing.T) {
	cases := []struct {
		name        string
		addons      []interface{}
		primaryIPv4 bool
		want        []string
	}{
		{"none", nil, false, []string{}},
		{"flag only", nil, true, []string{"primary_ipv4"}},
		{"explicit only", []interface{}{"failover_ip"}, false, []string{"failover_ip"}},
		{"flag adds to explicit", []interface{}{"failover_ip"}, true, []string{"failover_ip", "primary_ipv4"}},
		{"flag dedups when already listed", []interface{}{"primary_ipv4"}, true, []string{"primary_ipv4"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{"primary_ipv4": tc.primaryIPv4}
			if tc.addons != nil {
				raw["addons"] = tc.addons
			}
			d := resourceServerOrder().TestResourceData()
			if tc.addons != nil {
				d.Set("addons", tc.addons)
			}
			d.Set("primary_ipv4", tc.primaryIPv4)

			got := serverOrderAddons(d)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("serverOrderAddons() = %v, want %v", got, tc.want)
			}
		})
	}
}
