
// Package runtime provides isolation and resource control for RunInk node execution
package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CgroupV2Path is the path to the cgroup v2 filesystem
const CgroupV2Path = "/sys/fs/cgroup"

// ApplyCgroup creates a cgroup for the process and applies resource limits
// This provides resource control for CPU, memory, and I/O
func ApplyCgroup(name string, pid int, limits Limits) error {
	// Check if cgroup v2 is available
	if _, err := os.Stat(CgroupV2Path); os.IsNotExist(err) {
		return fmt.Errorf("cgroup v2 filesystem not mounted at %s", CgroupV2Path)
	}

	// Create cgroup directory
	cgroupPath := filepath.Join(CgroupV2Path, name)
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory: %v", err)
	}

	// Add process to cgroup
	if err := os.WriteFile(
		filepath.Join(cgroupPath, "cgroup.procs"),
		[]byte(strconv.Itoa(pid)),
		0644,
	); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %v", err)
	}

	// Apply CPU limits if specified
	if limits.CPUQuota != "" {
		if err := os.WriteFile(
			filepath.Join(cgroupPath, "cpu.max"),
			[]byte(limits.CPUQuota),
			0644,
		); err != nil {
			return fmt.Errorf("failed to set CPU quota: %v", err)
		}
	}

	// Apply memory limits if specified
	if limits.MemoryMax != "" {
		if err := os.WriteFile(
			filepath.Join(cgroupPath, "memory.max"),
			[]byte(limits.MemoryMax),
			0644,
		); err != nil {
			return fmt.Errorf("failed to set memory limit: %v", err)
		}
	}

	// Apply I/O weight if specified
	if limits.IOWeight != "" {
		if err := os.WriteFile(
			filepath.Join(cgroupPath, "io.weight"),
			[]byte(limits.IOWeight),
			0644,
		); err != nil {
			return fmt.Errorf("failed to set I/O weight: %v", err)
		}
	}

	return nil
}

// CleanupCgroup removes the cgroup directory
func CleanupCgroup(name string) error {
	cgroupPath := filepath.Join(CgroupV2Path, name)
	
	// Check if the cgroup exists
	if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
		return nil
	}

	// Remove the cgroup directory
	if err := os.RemoveAll(cgroupPath); err != nil {
		return fmt.Errorf("failed to remove cgroup directory: %v", err)
	}

	return nil
}

// ParseMemoryString converts a human-readable memory string to bytes
// Examples: "100M", "1G", "512K"
func ParseMemoryString(memStr string) (string, error) {
	memStr = strings.TrimSpace(memStr)
	if memStr == "" {
		return "", nil
	}

	// If it's already a number, return it
	if _, err := strconv.Atoi(memStr); err == nil {
		return memStr, nil
	}

	// Parse the unit
	unit := memStr[len(memStr)-1:]
	value := memStr[:len(memStr)-1]

	// Convert to bytes
	val, err := strconv.Atoi(value)
	if err != nil {
		return "", fmt.Errorf("invalid memory value: %s", value)
	}

	switch strings.ToUpper(unit) {
	case "K":
		return strconv.Itoa(val * 1024), nil
	case "M":
		return strconv.Itoa(val * 1024 * 1024), nil
	case "G":
		return strconv.Itoa(val * 1024 * 1024 * 1024), nil
	default:
		return "", fmt.Errorf("unknown memory unit: %s", unit)
	}
}
