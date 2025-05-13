package runtime

import (
	"fmt"
	"os"
	"testing"
)

// TestExecInSandbox demonstrates how to use the ExecInSandbox function
func TestExecInSandbox(t *testing.T) {
	// Skip this test if not running as root
	if os.Geteuid() != 0 {
		t.Skip("This test requires root privileges")
	}

	// Define resource limits
	limits := Limits{
		CPUQuota:  "100000 50000", // 50% CPU
		MemoryMax: "100M",         // 100MB memory
		IOWeight:  "100",          // Default I/O weight
	}

	// Create a temporary directory for the chroot
	tempDir, err := os.MkdirTemp("", "runink-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Prepare the chroot environment
	if err := PrepareChroot(tempDir); err != nil {
		t.Fatalf("Failed to prepare chroot: %v", err)
	}

	// Execute a simple command in the sandbox
	err = ExecInSandbox(
		[]string{"/bin/sh", "-c", "echo 'Hello from isolated environment'"},
		limits,
		tempDir,
	)

	if err != nil {
		t.Fatalf("Failed to execute command in sandbox: %v", err)
	}
}

// TestExecutor demonstrates how to use the Executor
func TestExecutor(t *testing.T) {
	// Skip this test if not running as root
	if os.Geteuid() != 0 {
		t.Skip("This test requires root privileges")
	}

	// Create a new executor
	executor := NewExecutor().
		SetCommand([]string{"/bin/sh", "-c", "echo 'Hello from isolated environment'"}).
		SetLimits(Limits{
			CPUQuota:  "100000 50000", // 50% CPU
			MemoryMax: "100M",         // 100MB memory
			IOWeight:  "100",          // Default I/O weight
		})

	// Execute the command
	result, err := executor.Execute()
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check the result
	if result.ExitCode != 0 {
		t.Fatalf("Command failed with exit code %d: %s", result.ExitCode, string(result.Stderr))
	}

	fmt.Printf("Command output: %s\n", string(result.Stdout))
}

// ExampleExecInSandbox demonstrates how to use the ExecInSandbox function
func ExampleExecInSandbox() {
	// Define resource limits
	limits := Limits{
		CPUQuota:  "100000 50000", // 50% CPU
		MemoryMax: "100M",         // 100MB memory
		IOWeight:  "100",          // Default I/O weight
	}

	// Execute a simple command in the sandbox
	err := ExecInSandbox(
		[]string{"/bin/sh", "-c", "echo 'Hello from isolated environment'"},
		limits,
		"", // Empty string means create a temporary directory
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Command executed successfully")
}

// ExampleExecutor demonstrates how to use the Executor
func ExampleExecutor() {
	// Create a new executor
	executor := NewExecutor().
		SetCommand([]string{"/bin/sh", "-c", "echo 'Hello from isolated environment'"}).
		SetLimits(Limits{
			CPUQuota:  "100000 50000", // 50% CPU
			MemoryMax: "100M",         // 100MB memory
			IOWeight:  "100",          // Default I/O weight
		})

	// Execute the command
	result, err := executor.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Check the result
	if result.ExitCode != 0 {
		fmt.Printf("Command failed with exit code %d: %s\n", result.ExitCode, string(result.Stderr))
		return
	}

	fmt.Printf("Command output: %s\n", string(result.Stdout))
}
