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

package targets

import (
	"errors"
	"fmt"
)

const (
	globalTid = 1
	globalIQN = "iqn.2017-10.io.opensds:volume:00000001"
	baseNum   = 100
)

var (
	globalLun = 1
)

type Target interface {
	CreateExport(path, initiator string) (Export, error)

	RemoveExport(path, initiator string) error
}

func NewTarget() Target {
	return &iscsiTarget{
		ISCSITarget: NewISCSITarget(globalTid, globalIQN),
	}
}

type iscsiTarget struct {
	ISCSITarget
}

func (t *iscsiTarget) CreateExport(path, initiator string) (Export, error) {
	globalLun = (globalLun + 1) % baseNum

	if t.GetISCSITarget() != globalTid {
		if err := t.CreateISCSITarget(); err != nil {
			return nil, err
		}
	}
	if err := t.AddLun(globalLun, path); err != nil {
		return nil, err
	}
	if err := t.BindInitiator(initiator); err != nil {
		return nil, err
	}

	return Export{
		"targetDiscovered": true,
		"targetIQN":        globalIQN,
		"targetPortal":     "127.0.0.1:2360",
		"discard":          false,
	}, nil
}

func (t *iscsiTarget) RemoveExport(path, initiator string) error {
	if err := t.UnbindInitiator(initiator); err != nil {
		return err
	}

	lun := t.GetLun(path)
	if lun == -1 {
		return errors.New("Can't find lun with path " + path)
	}
	if err := t.RemoveLun(lun); err != nil {
		return err
	}

	return nil
}

type Export map[string]interface{}

func (e Export) ValueToString() map[string]string {
	var out = make(map[string]string)

	for k := range e {
		out[k] = fmt.Sprint(e[k])
	}
	return out
}
