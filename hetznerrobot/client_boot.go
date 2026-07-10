package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#boot-configuration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

type BootProfile struct {
	ActiveProfile   string // linux/rescue/vnc/windows
	AuthorizedKeys  []string
	HostKeys        []string
	Language        string
	OperatingSystem string
	Password        string
	ServerNumber    int
	ServerIPv4      string
	ServerIPv6      string
}

// parseBootProfile extracts the currently-active boot profile from a /boot
// response body. Each profile carries its OS under a different key: linux and
// vnc use "dist", rescue and windows use "os".
func parseBootProfile(jsonStr string) *BootProfile {
	bootProfile := BootProfile{}
	activeBoot := ""

	if gjson.Get(jsonStr, "boot.linux.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.linux").String()
		bootProfile.ActiveProfile = "linux"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "dist").String()
	}
	if gjson.Get(jsonStr, "boot.rescue.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.rescue").String()
		bootProfile.ActiveProfile = "rescue"
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "os").String()
	}
	if gjson.Get(jsonStr, "boot.vnc.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.vnc").String()
		bootProfile.ActiveProfile = "vnc"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "dist").String()
	}
	if gjson.Get(jsonStr, "boot.windows.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.windows").String()
		bootProfile.ActiveProfile = "windows"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "os").String()
	}

	bootProfile.Password = gjson.Get(activeBoot, "password").String()
	bootProfile.ServerNumber = int(gjson.Get(activeBoot, "server_number").Int())
	bootProfile.ServerIPv4 = gjson.Get(activeBoot, "server_ip").String()
	bootProfile.ServerIPv6 = gjson.Get(activeBoot, "server_ipv6_net").String()

	return &bootProfile
}

func (c *HetznerRobotClient) getBoot(ctx context.Context, serverNumber int) (*BootProfile, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/boot/%d", c.url, serverNumber), nil, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		return nil, err
	}
	return parseBootProfile(string(bytes)), nil
}

// setBootProfile activates a boot profile. Parameter mapping follows Hetzner's
// API: linux/vnc take dist+lang, rescue takes os, windows takes os+lang. The
// operating_system resource field maps to "dist" or "os" depending on profile.
func (c *HetznerRobotClient) setBootProfile(ctx context.Context, serverNumber int, activeBootProfile string, os string, lang string, authorizedKeys []string) (*BootProfile, error) {
	data := url.Values{}
	for _, key := range authorizedKeys {
		data.Add("authorized_key", key)
	}
	switch activeBootProfile {
	case "linux", "vnc":
		data.Set("dist", os)
		data.Set("lang", lang)
	case "rescue":
		data.Set("os", os)
	case "windows":
		data.Set("os", os)
		data.Set("lang", lang)
	}

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/boot/%d/%s", c.url, serverNumber, activeBootProfile), data, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		if strings.Contains(err.Error(), "BOOT_ALREADY_ENABLED") {
			return c.getBoot(ctx, serverNumber)
		}
		return nil, err
	}
	return parseBootProfile(string(bytes)), nil
}

// deleteBootProfile deactivates the given boot profile (rescue/linux/vnc/windows)
// for a server via DELETE /boot/{server-number}/{profile}. A 404 is tolerated so
// destroying a resource whose profile was already consumed/inactive is not fatal.
func (c *HetznerRobotClient) deleteBootProfile(ctx context.Context, serverNumber int, profile string) error {
	_, err := c.makeAPICall(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/boot/%d/%s", c.url, serverNumber, profile),
		nil,
		[]int{http.StatusOK, http.StatusNotFound},
	)
	return err
}
