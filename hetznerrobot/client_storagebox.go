package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#storage-box

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type StorageBoxWrapper struct {
	StorageBox HetznerRobotStorageBox `json:"storagebox"`
}

type HetznerRobotStorageBox struct {
	ID                   int    `json:"id"`
	Login                string `json:"login"`
	Name                 string `json:"name"`
	Product              string `json:"product"`
	Cancelled            bool   `json:"cancelled"`
	Locked               bool   `json:"locked"`
	Location             string `json:"location"`
	LinkedServer         int    `json:"linked_server"`
	PaidUntil            string `json:"paid_until"`
	DiskQuota            int    `json:"disk_quota"`
	DiskUsage            int    `json:"disk_usage"`
	WebDAV               bool   `json:"webdav"`
	Samba                bool   `json:"samba"`
	SSH                  bool   `json:"ssh"`
	ExternalReachability bool   `json:"external_reachability"`
	ZFS                  bool   `json:"zfs"`
	Server               string `json:"server"`
	HostSystem           string `json:"host_system"`
}

func (c *HetznerRobotClient) getStorageBox(ctx context.Context, id int) (*HetznerRobotStorageBox, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/storagebox/%d", c.url, id), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := StorageBoxWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.StorageBox, nil
}

func (c *HetznerRobotClient) getStorageBoxes(ctx context.Context) ([]HetznerRobotStorageBox, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/storagebox", c.url), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	var wrappers []StorageBoxWrapper
	if err = json.Unmarshal(bytes, &wrappers); err != nil {
		return nil, err
	}
	boxes := make([]HetznerRobotStorageBox, len(wrappers))
	for i := range wrappers {
		boxes[i] = wrappers[i].StorageBox
	}
	return boxes, nil
}

// updateStorageBox updates a storage box's name and service toggles.
func (c *HetznerRobotClient) updateStorageBox(ctx context.Context, id int, name string, ssh, samba, webdav, externalReachability, zfs bool) (*HetznerRobotStorageBox, error) {
	data := url.Values{}
	data.Set("storagebox_name", name)
	data.Set("ssh", strconv.FormatBool(ssh))
	data.Set("samba", strconv.FormatBool(samba))
	data.Set("webdav", strconv.FormatBool(webdav))
	data.Set("external_reachability", strconv.FormatBool(externalReachability))
	data.Set("zfs", strconv.FormatBool(zfs))

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/storagebox/%d", c.url, id), data, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := StorageBoxWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.StorageBox, nil
}

// --- snapshots ---

type SnapshotWrapper struct {
	Snapshot HetznerRobotSnapshot `json:"snapshot"`
}

type HetznerRobotSnapshot struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	Size      int    `json:"size"`
	Automatic bool   `json:"automatic"`
}

func (c *HetznerRobotClient) createStorageBoxSnapshot(ctx context.Context, id int) (*HetznerRobotSnapshot, error) {
	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/storagebox/%d/snapshot", c.url, id), url.Values{}, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	wrapper := SnapshotWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Snapshot, nil
}

func (c *HetznerRobotClient) getStorageBoxSnapshots(ctx context.Context, id int) ([]HetznerRobotSnapshot, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/storagebox/%d/snapshot", c.url, id), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	var wrappers []SnapshotWrapper
	if err = json.Unmarshal(bytes, &wrappers); err != nil {
		return nil, err
	}
	snaps := make([]HetznerRobotSnapshot, len(wrappers))
	for i := range wrappers {
		snaps[i] = wrappers[i].Snapshot
	}
	return snaps, nil
}

func (c *HetznerRobotClient) deleteStorageBoxSnapshot(ctx context.Context, id int, name string) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/storagebox/%d/snapshot/%s", c.url, id, name), nil, []int{http.StatusOK, http.StatusNotFound})
	return err
}

// --- subaccounts ---

type SubaccountWrapper struct {
	Subaccount HetznerRobotSubaccount `json:"subaccount"`
}

type HetznerRobotSubaccount struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	Server               string `json:"server"`
	HomeDirectory        string `json:"homedirectory"`
	Samba                bool   `json:"samba"`
	SSH                  bool   `json:"ssh"`
	ExternalReachability bool   `json:"external_reachability"`
	WebDAV               bool   `json:"webdav"`
	Readonly             bool   `json:"readonly"`
	Comment              string `json:"comment"`
}

func subaccountForm(homeDir string, samba, ssh, externalReachability, webdav, readonly bool, comment string) url.Values {
	data := url.Values{}
	data.Set("homedirectory", homeDir)
	data.Set("samba", strconv.FormatBool(samba))
	data.Set("ssh", strconv.FormatBool(ssh))
	data.Set("external_reachability", strconv.FormatBool(externalReachability))
	data.Set("webdav", strconv.FormatBool(webdav))
	data.Set("readonly", strconv.FormatBool(readonly))
	if comment != "" {
		data.Set("comment", comment)
	}
	return data
}

func (c *HetznerRobotClient) createStorageBoxSubaccount(ctx context.Context, id int, homeDir string, samba, ssh, externalReachability, webdav, readonly bool, comment string) (*HetznerRobotSubaccount, error) {
	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/storagebox/%d/subaccount", c.url, id), subaccountForm(homeDir, samba, ssh, externalReachability, webdav, readonly, comment), []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	wrapper := SubaccountWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Subaccount, nil
}

func (c *HetznerRobotClient) updateStorageBoxSubaccount(ctx context.Context, id int, username, homeDir string, samba, ssh, externalReachability, webdav, readonly bool, comment string) error {
	_, err := c.makeAPICall(ctx, "PUT", fmt.Sprintf("%s/storagebox/%d/subaccount/%s", c.url, id, username), subaccountForm(homeDir, samba, ssh, externalReachability, webdav, readonly, comment), []int{http.StatusOK})
	return err
}

func (c *HetznerRobotClient) getStorageBoxSubaccounts(ctx context.Context, id int) ([]HetznerRobotSubaccount, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/storagebox/%d/subaccount", c.url, id), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	var wrappers []SubaccountWrapper
	if err = json.Unmarshal(bytes, &wrappers); err != nil {
		return nil, err
	}
	subs := make([]HetznerRobotSubaccount, len(wrappers))
	for i := range wrappers {
		subs[i] = wrappers[i].Subaccount
	}
	return subs, nil
}

func (c *HetznerRobotClient) deleteStorageBoxSubaccount(ctx context.Context, id int, username string) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/storagebox/%d/subaccount/%s", c.url, id, username), nil, []int{http.StatusOK, http.StatusNotFound})
	return err
}
