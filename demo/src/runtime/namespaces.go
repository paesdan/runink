
// Package runtime provides isolation and resource control for RunInk node execution
package runtime

import (
	"fmt"
	"syscall"
)

// Namespace constants for Linux
const (
	// CLONE_NEWUTS creates a new UTS namespace (hostname, domainname)
	CLONE_NEWUTS = syscall.CLONE_NEWUTS

	// CLONE_NEWPID creates a new PID namespace
	CLONE_NEWPID = syscall.CLONE_NEWPID

	// CLONE_NEWNS creates a new mount namespace
	CLONE_NEWNS = syscall.CLONE_NEWNS

	// CLONE_NEWIPC creates a new IPC namespace
	CLONE_NEWIPC = syscall.CLONE_NEWIPC

	// CLONE_NEWNET creates a new network namespace
	CLONE_NEWNET = syscall.CLONE_NEWNET

	// CLONE_NEWUSER creates a new user namespace
	// Note: This requires special privileges and is not used by default
	CLONE_NEWUSER = syscall.CLONE_NEWUSER

	// DefaultNamespaces is the default set of namespaces to unshare
	DefaultNamespaces = CLONE_NEWUTS | CLONE_NEWPID | CLONE_NEWNS | CLONE_NEWIPC
)

// EnterNamespaces creates new namespaces for the current process
// This provides isolation for various system resources
func EnterNamespaces(namespaces int) error {
	if namespaces == 0 {
		namespaces = DefaultNamespaces
	}

	// Unshare the specified namespaces
	if err := syscall.Unshare(namespaces); err != nil {
		return fmt.Errorf("failed to unshare namespaces: %v", err)
	}

	return nil
}

// SetHostname sets the hostname in the current UTS namespace
// This is useful after entering a new UTS namespace
func SetHostname(hostname string) error {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return fmt.Errorf("failed to set hostname: %v", err)
	}
	return nil
}

// MountProc mounts the proc filesystem
// This is necessary after entering a new PID namespace
func MountProc() error {
	// Ensure /proc exists
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("failed to mount /proc: %v", err)
	}
	return nil
}
