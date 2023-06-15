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
	"os"
	"os/exec"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

type Cli interface {
	Mount(volumeName string, mountPath string, options []string) error
	Umount(mountPath string) error
	Create(name string, server string, size int64, onDeletePolicy string) (*sealfsVolume, error)
	Delete(name string) error
	Probe() error
}

type SealfsCli struct {
}

func (f *SealfsCli) Mount(volumeName string, mountPath string, options []string) error {
	if volumeName == "" {
		return status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	if mountPath == "" {
		return status.Error(codes.InvalidArgument, "Target path not provided")
	}

	// mkdir -p mountPath
	if err := os.MkdirAll(mountPath, 0750); err != nil {
		return fmt.Errorf("failed to create mount path %s: %v", mountPath, err)
	}

	sealfsCmd := "sealfs-client"
	sealfsArgs := []string{"--log-level", "warn", "mount", mountPath, volumeName}
	mountArgsLogStr := ""
	for _, arg := range sealfsArgs {
		mountArgsLogStr = mountArgsLogStr + " " + arg
	}

	klog.V(4).Infof("Mounting cmd (%s) with arguments (%s)", sealfsCmd, sealfsArgs)
	command := exec.Command(sealfsCmd, sealfsArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		if err.Error() == "wait: no child processes" {
			if command.ProcessState.Success() {
				return nil
			}
			// Rewrite err with the actual exit error of the process.
			err = &exec.ExitError{ProcessState: command.ProcessState}
		}
		klog.Errorf("Mount failed: %v\nMounting command: %s\nMounting arguments: %s\nOutput: %s\n", err, sealfsCmd, mountArgsLogStr, string(output))
		return fmt.Errorf("mount failed: %v\nMounting command: %s\nMounting arguments: %s\nOutput: %s",
			err, sealfsCmd, mountArgsLogStr, string(output))
	}
	return nil
}

func (f *SealfsCli) Umount(mountPath string) error {
	if mountPath == "" {
		return status.Error(codes.InvalidArgument, "Target path not provided")
	}

	sealfsCmd := "sealfs-client"
	sealfsArgs := []string{"--log-level", "warn", "umount", mountPath}
	mountArgsLogStr := ""
	for _, arg := range sealfsArgs {
		mountArgsLogStr = mountArgsLogStr + " " + arg
	}

	klog.V(4).Infof("Umounting cmd (%s) with arguments (%s)", sealfsCmd, sealfsArgs)
	command := exec.Command(sealfsCmd, sealfsArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		if err.Error() == "wait: no child processes" {
			if command.ProcessState.Success() {
				return nil
			}
			// Rewrite err with the actual exit error of the process.
			err = &exec.ExitError{ProcessState: command.ProcessState}
		}
		klog.Errorf("Umount failed: %v\nUmounting command: %s\nUmounting arguments: %s\nOutput: %s\n", err, sealfsCmd, mountArgsLogStr, string(output))
		return fmt.Errorf("Umount failed: %v\nUmounting command: %s\nUmounting arguments: %s\nOutput: %s",
			err, sealfsCmd, mountArgsLogStr, string(output))
	}

	// rm -rf mountPath

	if err := os.RemoveAll(mountPath); err != nil {
		return fmt.Errorf("failed to remove mount path %s: %v", mountPath, err)
	}

	return nil
}

func (f *SealfsCli) Create(volumeName string, server string, size int64, onDeletePolicy string) (*sealfsVolume, error) {

	sealfsCmd := "sealfs-client"
	sealfsArgs := []string{"--log-level", "warn", "create", "-m", server, volumeName, "100000"}
	mountArgsLogStr := ""
	for _, arg := range sealfsArgs {
		mountArgsLogStr = mountArgsLogStr + " " + arg
	}

	klog.V(4).Infof("Creating cmd (%s) with arguments (%s)", sealfsCmd, sealfsArgs)
	command := exec.Command(sealfsCmd, sealfsArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		if err.Error() != "wait: no child processes" || !command.ProcessState.Success() {
			// Rewrite err with the actual exit error of the process.
			if err.Error() == "wait: no child processes" {
				err = &exec.ExitError{ProcessState: command.ProcessState}
			}
			klog.Errorf("Mount failed: %v\nMounting command: %s\nMounting arguments: %s\nOutput: %s\n", err, sealfsCmd, mountArgsLogStr, string(output))
			return nil, fmt.Errorf("mount failed: %v\nMounting command: %s\nMounting arguments: %s\nOutput: %s",
				err, sealfsCmd, mountArgsLogStr, string(output))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create volume %v: %v", volumeName, err)
	}

	volume := &sealfsVolume{
		id:       volumeName,
		size:     size,
		uuid:     volumeName,
		onDelete: onDeletePolicy,
	}

	return volume, nil
}

func (f *SealfsCli) Delete(name string) error {
	if name == "" {
		return fmt.Errorf("missing required volume name")
	}

	return nil
}

func (f *SealfsCli) Probe() error {
	sealfsCmd := "sealfs-client"
	sealfsArgs := []string{"--log-level", "warn", "probe"}
	mountArgsLogStr := ""
	for _, arg := range sealfsArgs {
		mountArgsLogStr = mountArgsLogStr + " " + arg
	}

	klog.V(4).Infof("Creating cmd (%s) with arguments (%s)", sealfsCmd, sealfsArgs)
	command := exec.Command(sealfsCmd, sealfsArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		if err.Error() != "wait: no child processes" || !command.ProcessState.Success() {
			// Rewrite err with the actual exit error of the process.
			if err.Error() == "wait: no child processes" {
				err = &exec.ExitError{ProcessState: command.ProcessState}
			}
			klog.Errorf("Mount failed: %v\nMounting command: %s\nMounting arguments: %s\nOutput: %s\n", err, sealfsCmd, mountArgsLogStr, string(output))
			return err
		}
	}

	return nil
}
