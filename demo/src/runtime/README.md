# RunInk Runtime Package

This package provides isolation and resource control features for RunInk node execution using Linux kernel features:

1. **Chroot** - Filesystem isolation
2. **Namespaces** - Process isolation
3. **Cgroups** - Resource control

## Components

- **types.go** - Defines data structures for resource limits and executor configuration
- **chroot.go** - Provides functions for filesystem isolation using chroot
- **namespaces.go** - Provides functions for process isolation using Linux namespaces
- **cgroup.go** - Provides functions for resource control using cgroups v2
- **executor.go** - Provides a high-level API for executing commands in isolated environments

## Usage

### Basic Usage

```go
import "github.com/runink/runink/runtime"

// Define resource limits
limits := runtime.Limits{
    CPUQuota:  "100000 50000", // 50% CPU
    MemoryMax: "100M",         // 100MB memory
    IOWeight:  "100",          // Default I/O weight
}

// Execute a command in an isolated environment
err := runtime.ExecInSandbox(
    []string{"/bin/sh", "-c", "echo 'Hello from isolated environment'"},
    limits,
    "/tmp/runink_chroot", // Optional chroot directory
)
```

### Advanced Usage with Executor

```go
import "github.com/runink/runink/runtime"

// Create a new executor
executor := runtime.NewExecutor().
    SetCommand([]string{"/bin/sh", "-c", "echo 'Hello from isolated environment'"}).
    SetLimits(runtime.Limits{
        CPUQuota:  "100000 50000", // 50% CPU
        MemoryMax: "100M",         // 100MB memory
        IOWeight:  "100",          // Default I/O weight
    }).
    SetChrootDir("/tmp/runink_chroot").
    SetNamespaces(runtime.CLONE_NEWUTS | runtime.CLONE_NEWPID | runtime.CLONE_NEWNS)

// Execute the command
result, err := executor.Execute()
if err != nil {
    // Handle error
}

// Process the result
fmt.Println("Exit code:", result.ExitCode)
fmt.Println("Stdout:", string(result.Stdout))
fmt.Println("Stderr:", string(result.Stderr))
```

## Integration with DAG Execution

The runtime package is integrated with the DAG execution engine in the `internal/engine/execute.go` file. When a node in the DAG is executed, the engine extracts the command and resource limits from the node configuration and executes the command in an isolated environment using the runtime package.

### Node Configuration

Nodes in the DAG can specify the following configuration options for isolation and resource control:

- `command` - The command to execute
- `cpu_quota` - CPU quota in the format "period quota" (e.g., "100000 50000" for 50% CPU)
- `memory_max` - Maximum memory limit (e.g., "100M" for 100MB)
- `io_weight` - I/O weight (e.g., "100" for default I/O weight)
- `chroot` - Chroot directory for filesystem isolation

### Example DSL

```
# Example RunInk DSL with isolation and resource control
source: "file:///tmp/input.txt"

# Simple echo command with default isolation
transform echo_step (command: "echo 'Hello from isolated environment'", memory_max: "50M")

# CPU-intensive task with CPU limits
transform cpu_task (command: "for i in {1..1000}; do echo $i; done | sort -n", cpu_quota: "100000 30000", after: echo_step)
```

## Requirements

- Linux kernel with cgroups v2 support
- Root privileges for some operations (chroot, namespaces)
- The cgroup v2 filesystem must be mounted at `/sys/fs/cgroup`

## Limitations

This is a minimal viable implementation with the following limitations:

1. Limited error handling and recovery
2. No support for user namespaces (requires additional setup)
3. No support for network namespaces (requires additional setup)
4. No support for custom mount points in the chroot environment
5. No support for resource usage monitoring
6. No support for resource usage accounting

## Future Improvements

1. Add support for user namespaces
2. Add support for network namespaces
3. Add support for custom mount points
4. Add support for resource usage monitoring
5. Add support for resource usage accounting
6. Add support for resource usage limits enforcement
7. Add support for resource usage limits violation handling
