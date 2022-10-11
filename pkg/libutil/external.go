package libutil

import (
	"encoding/json"
	"os/exec"
)

// run external image maker
func RunExternalImageMaker(path string, username string, publicKey string, mac string, vlan int32,
	ip string, gateway string, dns string) (string, error) {

	type JsonData struct {
		Path      string   `json:"path"`
		Username  string   `json:"username"`
		PublicKey string   `json:"public_key"`
		Mac       string   `json:"mac"`
		Vlan      int32    `json:"vlan"`
		Ip        string   `json:"ip"`
		Gateway   string   `json:"gateway"`
		Dns       string   `json:"dns"`
		CmdList   []string `json:"cmdList"`
	}

	data := JsonData{
		Path:      path,
		Username:  username,
		PublicKey: publicKey,
		Mac:       mac,
		Vlan:      vlan,
		Ip:        ip,
		Gateway:   gateway,
		Dns:       dns,
		CmdList:   []string{},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("external/image_maker", string(jsonData))
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
