
// Package runtime provides isolation and resource control for RunInk node execution
package runtime

// Limits defines resource limits for a node execution
type Limits struct {
	// CPUQuota defines the CPU quota in the format "period quota"
	// e.g., "100000 50000" means 50% of CPU time
	CPUQuota string

	// MemoryMax defines the maximum memory limit in bytes
	// e.g., "100M" for 100 megabytes
	MemoryMax string

	// DiskQuota defines the maximum disk space in bytes
	// e.g., "1G" for 1 gigabyte
	DiskQuota string

	// IOWeight defines the I/O weight for the cgroup
	// Valid values are 10-1000
	IOWeight string
}

// ExecutorConfig holds configuration for the isolated execution environment
type ExecutorConfig struct {
	// Command to execute
	Command []string

	// Resource limits
	Limits Limits

	// Root directory for chroot
	ChrootDir string

	// Working directory inside the chroot
	WorkDir string

	// Environment variables
	Env []string

	// Namespaces to unshare
	// Default: CLONE_NEWUTS | CLONE_NEWPID | CLONE_NEWNS | CLONE_NEWIPC
	Namespaces int

	// CgroupName is the name of the cgroup to create
	// Default: "runink-<pid>"
	CgroupName string
}

// ExecutorResult represents the result of an execution
type ExecutorResult struct {
	// Exit code of the command
	ExitCode int

	// Standard output
	Stdout []byte

	// Standard error
	Stderr []byte

	// Error if any occurred during execution
	Error error
}
