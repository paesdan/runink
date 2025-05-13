
// Package runtime provides isolation and resource control for RunInk node execution
package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// ApplyChroot changes the root directory to the specified path
// This provides filesystem isolation for the process
func ApplyChroot(root string) error {
	// Ensure the root directory exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return fmt.Errorf("chroot directory %s does not exist", root)
	}

	// Change to the new root directory
	if err := syscall.Chdir(root); err != nil {
		return fmt.Errorf("failed to chdir to %s: %v", root, err)
	}

	// Change the root directory
	if err := syscall.Chroot(root); err != nil {
		return fmt.Errorf("failed to chroot to %s: %v", root, err)
	}

	// Change to the root directory inside the chroot
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("failed to chdir to / after chroot: %v", err)
	}

	return nil
}

// PrepareChroot creates a minimal chroot environment
// This is a helper function to set up a basic chroot directory
func PrepareChroot(rootDir string) error {
	// Create the root directory if it doesn't exist
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return fmt.Errorf("failed to create chroot directory: %v", err)
	}

	// Create basic directories inside the chroot
	dirs := []string{
		"bin", "dev", "etc", "lib", "lib64",
		"proc", "sys", "tmp", "usr", "var",
	}

	for _, dir := range dirs {
		path := filepath.Join(rootDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Create /dev/null
	devNull := filepath.Join(rootDir, "dev", "null")
	if _, err := os.Stat(devNull); os.IsNotExist(err) {
		f, err := os.Create(devNull)
		if err != nil {
			return fmt.Errorf("failed to create /dev/null: %v", err)
		}
		f.Close()
	}

	return nil
}
