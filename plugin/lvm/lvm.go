// Copyright (c) 2017 OpenSDS Authors.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package lvm

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/leonwanghui/lvm-csi/plugin/lvm/targets"
	"github.com/satori/go.uuid"
)

const (
	vgName = "vg001"
)

var lvmDriver = &Driver{}

type Driver struct{}

func (d *Driver) CreateVolume(name string, size uint64) (*csi.VolumeInfo, error) {
	var capacity = fmt.Sprint(size) + "G"

	cmd := strings.Join([]string{"lvcreate", "-n", name, "-L", capacity, vgName}, " ")
	if _, err := d.execCmd(cmd); err != nil {
		log.Println("Failed to create logic volume:", err)
		return nil, err
	}

	var lvPath, lvStatus string
	// Display and parse some metadata in logic volume returned.
	lvPath = strings.Join([]string{"/dev", vgName, name}, "/")
	lv, err := d.execCmd("lvdisplay " + lvPath)
	if err != nil {
		log.Println("Failed to display logic volume:", err)
		return nil, err
	}

	for _, line := range strings.Split(lv, "\n") {
		if strings.Contains(line, "LV Path") {
			lvPath = strings.Fields(line)[2]
		}
		if strings.Contains(line, "LV Status") {
			lvStatus = strings.Fields(line)[2]
		}
	}

	// return volume info
	return &csi.VolumeInfo{
		Id: uuid.NewV4().String(),
		Attributes: map[string]string{
			"name":   name,
			"status": lvStatus,
			"lvPath": lvPath,
		},
		CapacityBytes: size,
	}, nil
}

func (d *Driver) ListVolumes() (*csi.VolumeInfo, error) {
	// Display and parse some metadata in logic volume returned.
	_, err := d.execCmd("lvdisplay " + vgName)
	if err != nil {
		log.Println("Failed to display logic volume:", err)
		return nil, err
	}

	// return volume info
	return nil, nil
}

func (d *Driver) DeleteVolume(lvPath string) error {
	cmd := strings.Join([]string{"lvremove", "-f", lvPath}, " ")
	if _, err := d.execCmd(cmd); err != nil {
		log.Println("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

func (*Driver) InitializeConnection(lvPath, initiator string) (map[string]string, error) {
	if initiator == "" {
		initiator = "ALL"
	}

	t := targets.NewTarget()
	expt, err := t.CreateExport(lvPath, initiator)
	if err != nil {
		log.Println("Failed to initialize connection of logic volume:", err)
		return nil, err
	}

	expt["driverVolumeType"] = "iscsi"
	return expt.ValueToString(), nil
}

func (*Driver) TerminateConnection(lvPath, initiator string) error {
	if initiator == "" {
		initiator = "ALL"
	}

	return targets.NewTarget().RemoveExport(lvPath, initiator)
}

func (*Driver) execCmd(cmd string) (string, error) {
	ret, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return string(ret), nil
}
