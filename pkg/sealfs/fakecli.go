/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sealfs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeCli struct {
}

func (f *fakeCli) Mount(volumeName string, mountPath string, options []string) error {
	if volumeName == "" {
		return status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	if mountPath == "" {
		return status.Error(codes.InvalidArgument, "Target path not provided")
	}

	return nil
}

func (f *fakeCli) Umount(mountPath string) error {
	if mountPath == "" {
		return status.Error(codes.InvalidArgument, "Target path not provided")
	}

	return nil
}

func (f *fakeCli) Create(volumeName string, server string, size int64, onDeletePolicy string) (*sealfsVolume, error) {

	volume := &sealfsVolume{
		id:       volumeName,
		size:     size,
		uuid:     volumeName,
		onDelete: onDeletePolicy,
	}

	return volume, nil
}

func (f *fakeCli) Delete(name string) error {
	if name == "" {
		return fmt.Errorf("missing required volume name")
	}

	return nil
}

func (f *fakeCli) Probe() error {
	return nil
}
