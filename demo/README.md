# RunInk

RunInk is a command-line tool for parsing contract, conf, dsl, and herd files to build and execute Directed Acyclic Graphs (DAGs).

## Project Structure

```
runink_demo/
├── README.md
└── src/
    ├── cmd/
    │   ├── execute.go
    │   ├── root.go
    │   └── run.go
    ├── internal/
    │   ├── engine/
    │   │   ├── dag.go
    │   │   ├── graph.go
    │   │   └── parser.go
    │   └── parser/
    │       ├── conf.go
    │       ├── contract.go
    │       ├── dsl.go
    │       └── herd.go
    ├── go.mod
    ├── go.sum
    └── main.go
```

## Features

- Parse contract, conf, dsl, and herd files
- Build a Directed Acyclic Graph (DAG) based on the parsed files
- Execute the DAG in the correct order
- Handle errors and provide feedback
- Advanced error handling and recovery
- Flexible data passing between nodes
- Comprehensive monitoring and observability
- Multiple execution modes (sync/async)
- Isolation support for secure execution

## Running RunInk from the Command Line

RunInk provides a powerful command-line interface with numerous options to customize the execution of your DAGs. Below are detailed instructions on how to use the CLI.

### Basic Usage

```bash
# Run with required files
runink run --contract file.contract --conf file.conf --dsl file.dsl

# Run with optional herd file
runink run --contract file.contract --conf file.conf --dsl file.dsl --herd file.herd

# Run with verbose output
runink run --contract file.contract --conf file.conf --dsl file.dsl --verbose
```

### Advanced Options

RunInk's enhanced execution engine supports several advanced options that can be configured via command-line flags:

#### Error Handling

Control how RunInk handles errors during execution:

```bash
# Stop execution on first error (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl --error-mode stop

# Continue execution when errors occur in non-critical nodes
runink run --contract file.contract --conf file.conf --dsl file.dsl --error-mode continue
```

#### Data Passing Strategy

Configure how data is passed between nodes in the DAG:

```bash
# Use JSON for data passing (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl --data-pass json

# Use files for data passing (useful for large data sets)
runink run --contract file.contract --conf file.conf --dsl file.dsl --data-pass file

# Use stdout for data passing (useful for debugging)
runink run --contract file.contract --conf file.conf --dsl file.dsl --data-pass stdout
```

#### Monitoring Level

Set the level of execution monitoring:

```bash
# Basic monitoring (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl --monitoring basic

# Verbose monitoring with detailed execution information
runink run --contract file.contract --conf file.conf --dsl file.dsl --monitoring verbose

# No monitoring for maximum performance
runink run --contract file.contract --conf file.conf --dsl file.dsl --monitoring none
```

#### Execution Mode

Choose between synchronous and asynchronous execution:

```bash
# Synchronous execution (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl --execution-mode sync

# Asynchronous execution for parallel processing
runink run --contract file.contract --conf file.conf --dsl file.dsl --execution-mode async
```

#### Integrated Flow

Enable or disable integrated execution flow:

```bash
# Enable integrated flow (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl --integrated-flow=true

# Disable integrated flow
runink run --contract file.contract --conf file.conf --dsl file.dsl --integrated-flow=false
```

#### Isolation Level

Set the isolation level for secure execution:

```bash
# No isolation (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl --isolation-level none

# Process-level isolation
runink run --contract file.contract --conf file.conf --dsl file.dsl --isolation-level process

# Container-level isolation
runink run --contract file.contract --conf file.conf --dsl file.dsl --isolation-level container
```

#### Run ID

Specify a unique identifier for the execution:

```bash
# Auto-generated run ID (default)
runink run --contract file.contract --conf file.conf --dsl file.dsl

# Custom run ID
runink run --contract file.contract --conf file.conf --dsl file.dsl --run-id my-custom-run-id
```

### Complete Example

Here's an example that combines multiple advanced options:

```bash
runink run \
  --contract workflow.contract \
  --conf production.conf \
  --dsl pipeline.dsl \
  --herd resources.herd \
  --error-mode continue \
  --data-pass file \
  --monitoring verbose \
  --execution-mode async \
  --integrated-flow=true \
  --isolation-level process \
  --run-id production-run-2025-05-11 \
  --verbose
```

This command will:
- Use the specified contract, conf, dsl, and herd files
- Continue execution even if non-critical errors occur
- Use files for data passing between nodes
- Enable verbose monitoring for detailed execution information
- Run nodes asynchronously where possible
- Use the integrated execution flow
- Apply process-level isolation for security
- Use a custom run ID for tracking
- Enable verbose output for additional information

## Development

This project is currently in the skeleton implementation phase. The following components are planned:

1. File parsers for different file types
2. DAG builder based on parsed files
3. DAG execution engine
4. Error handling and reporting

## Building from Source

```bash
cd ~/runink_demo/src
go build -o runink
```

## Running Tests

```bash
cd ~/runink_demo/src
go test ./...
```

## Enhanced Execution Engine

The RunInk execution engine has been enhanced with several advanced features to improve reliability, performance, and observability:

### Advanced Error Handling and Recovery
- Panic recovery for individual nodes to prevent cascading failures
- Configurable error handling policies (fatal vs. non-fatal errors)
- Detailed error collection and reporting
- Graceful shutdown on critical failures with proper resource cleanup

### Typed Channel Data Passing
- Support for both generic (`any`) and strongly-typed data channels
- Automatic channel management (creation, closing)
- Buffered channels to prevent deadlocks
- Support for multiple input dependencies with proper synchronization

### Execution Modes
- **Sync Mode** (`engine.SyncMode`): Executes nodes in sequence, waiting for each to complete
- **Async Mode** (`engine.AsyncMode`): Executes nodes in parallel where dependencies allow
- Configurable at runtime to adapt to different workload requirements

### Monitoring and Observability
- Built-in monitoring hooks for node start/success/error events
- Execution timing and duration tracking
- Customizable monitoring implementations
- Aggregated execution statistics for performance analysis

### Isolation Support
- Process isolation for security
- Container isolation for complete environment separation
- Configurable isolation policies
- Isolation ID for tracking and management

### Usage Example (Go API)

```go
import (
    "context"
    "github.com/runink/internal/engine"
)

// Build your DAG
dag := engine.BuildDAG(files, true, "run-123")

// Create context, monitor, and error handler
ctx := context.Background()
monitor := engine.NewDefaultMonitor()
errHandler := engine.NewDefaultErrorHandler(true)

// Execute the DAG
err := engine.Execute(ctx, dag, engine.SyncMode, monitor, errHandler)
```

For more detailed information, see the [Engine Documentation](/src/internal/engine/README.md).
