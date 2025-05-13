// Package engine provides the DAG execution engine for RunInk.
package engine

import (
        "context"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "runtime"
        "sync"
        "time"
)

// NodeMetric contains metrics for a single node execution
type NodeMetric struct {
        NodeID       string        `json:"nodeId"`
        NodeName     string        `json:"nodeName"`
        StartTime    time.Time     `json:"startTime"`
        EndTime      time.Time     `json:"endTime"`
        Duration     time.Duration `json:"duration"`
        MemoryStart  uint64        `json:"memoryStart"`
        MemoryEnd    uint64        `json:"memoryEnd"`
        MemoryDelta  int64         `json:"memoryDelta"`
        CPUTimeStart time.Duration `json:"cpuTimeStart"`
        CPUTimeEnd   time.Duration `json:"cpuTimeEnd"`
        Status       NodeState     `json:"status"`
        Error        string        `json:"error,omitempty"`
        RetryCount   int           `json:"retryCount"`
}

// MonitorSnapshot provides an immutable snapshot of the current execution state
type MonitorSnapshot struct {
        DAGName         string                 `json:"dagName"`
        StartTime       time.Time              `json:"startTime"`
        CurrentTime     time.Time              `json:"currentTime"`
        ElapsedTime     time.Duration          `json:"elapsedTime"`
        TotalNodes      int                    `json:"totalNodes"`
        CompletedNodes  int                    `json:"completedNodes"`
        RunningNodes    int                    `json:"runningNodes"`
        PendingNodes    int                    `json:"pendingNodes"`
        FailedNodes     int                    `json:"failedNodes"`
        SkippedNodes    int                    `json:"skippedNodes"`
        RetryingNodes   int                    `json:"retryingNodes"`
        NodeMetrics     map[string]NodeMetric  `json:"nodeMetrics"`
        ResourceMetrics ResourceMetrics        `json:"resourceMetrics"`
        Progress        float64                `json:"progress"`
}

// ResourceMetrics tracks overall resource usage for the DAG execution
type ResourceMetrics struct {
        PeakMemory     uint64        `json:"peakMemory"`
        CurrentMemory  uint64        `json:"currentMemory"`
        TotalCPUTime   time.Duration `json:"totalCPUTime"`
        GoroutineCount int           `json:"goroutineCount"`
}

// Monitor provides monitoring and observability for DAG execution
type Monitor struct {
        mu            sync.RWMutex
        dagName       string
        startTime     time.Time
        nodeMetrics   map[string]*NodeMetric
        totalNodes    int
        httpServer    *http.Server
        resourceTicks []ResourceMetrics
        tickInterval  time.Duration
        done          chan struct{}
}

// NewMonitor creates a new execution monitor
func NewMonitor(dagName string, totalNodes int) *Monitor {
        return &Monitor{
                dagName:      dagName,
                startTime:    time.Now(),
                nodeMetrics:  make(map[string]*NodeMetric),
                totalNodes:   totalNodes,
                tickInterval: time.Second,
                done:         make(chan struct{}),
        }
}

// StartNode begins monitoring for a node execution and returns a function to be deferred
// for completing the monitoring when the node finishes execution
func (m *Monitor) StartNode(nodeID, nodeName string) func(status NodeState, err error, retryCount int) {
        m.mu.Lock()
        
        // Get current memory stats
        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)
        
        // Create new node metric
        metric := &NodeMetric{
                NodeID:       nodeID,
                NodeName:     nodeName,
                StartTime:    time.Now(),
                MemoryStart:  memStats.Alloc,
                CPUTimeStart: getCPUTime(),
                Status:       NodeRunning,
        }
        
        m.nodeMetrics[nodeID] = metric
        m.mu.Unlock()
        
        // Log node start
        log.Printf("▶ Starting node: %s (%s)", nodeName, nodeID)
        
        // Return function to be called when node completes
        return func(status NodeState, err error, retryCount int) {
                m.mu.Lock()
                defer m.mu.Unlock()
                
                // Get current memory stats
                runtime.ReadMemStats(&memStats)
                
                // Update metric
                metric.EndTime = time.Now()
                metric.Duration = metric.EndTime.Sub(metric.StartTime)
                metric.MemoryEnd = memStats.Alloc
                metric.MemoryDelta = int64(metric.MemoryEnd) - int64(metric.MemoryStart)
                metric.CPUTimeEnd = getCPUTime()
                metric.Status = status
                metric.RetryCount = retryCount
                
                if err != nil {
                        metric.Error = err.Error()
                }
                
                // Log node completion
                var statusSymbol string
                switch status {
                case NodeSucceeded:
                        statusSymbol = "✓"
                case NodeFailed:
                        statusSymbol = "✗"
                case NodeRetrying:
                        statusSymbol = "↻"
                case NodeSkipped:
                        statusSymbol = "→"
                default:
                        statusSymbol = "?"
                }
                
                errStr := ""
                if err != nil {
                        errStr = fmt.Sprintf(" error=%v", err)
                }
                log.Printf("%s Completed node: %s (%s) duration=%v status=%s retries=%d%s",
                        statusSymbol, nodeName, nodeID, metric.Duration, status, retryCount, errStr)
        }
}

// Snapshot returns an immutable snapshot of the current execution state
func (m *Monitor) Snapshot() MonitorSnapshot {
        m.mu.RLock()
        defer m.mu.RUnlock()
        
        // Get current memory stats
        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)
        
        // Count nodes by status
        var completed, running, pending, failed, skipped, retrying int
        nodeMetrics := make(map[string]NodeMetric)
        
        for id, metric := range m.nodeMetrics {
                // Make a copy of the metric
                nodeMetrics[id] = *metric
                
                // Count by status
                switch metric.Status {
                case NodeSucceeded:
                        completed++
                case NodeRunning:
                        running++
                case NodePending:
                        pending++
                case NodeFailed:
                        failed++
                case NodeSkipped:
                        skipped++
                case NodeRetrying:
                        retrying++
                }
        }
        
        // Calculate pending nodes
        pending = m.totalNodes - (completed + running + failed + skipped + retrying)
        if pending < 0 {
                pending = 0
        }
        
        // Calculate progress
        progress := 0.0
        if m.totalNodes > 0 {
                progress = float64(completed+failed+skipped) / float64(m.totalNodes)
        }
        
        // Create resource metrics
        resourceMetrics := ResourceMetrics{
                CurrentMemory:  memStats.Alloc,
                PeakMemory:     getPeakMemory(m.resourceTicks),
                TotalCPUTime:   getCPUTime(),
                GoroutineCount: runtime.NumGoroutine(),
        }
        
        return MonitorSnapshot{
                DAGName:        m.dagName,
                StartTime:      m.startTime,
                CurrentTime:    time.Now(),
                ElapsedTime:    time.Since(m.startTime),
                TotalNodes:     m.totalNodes,
                CompletedNodes: completed,
                RunningNodes:   running,
                PendingNodes:   pending,
                FailedNodes:    failed,
                SkippedNodes:   skipped,
                RetryingNodes:  retrying,
                NodeMetrics:    nodeMetrics,
                ResourceMetrics: resourceMetrics,
                Progress:       progress,
        }
}

// StartHTTPServer starts an HTTP server for monitoring on the specified address
func (m *Monitor) StartHTTPServer(addr string) error {
        mux := http.NewServeMux()
        
        // Handler for JSON metrics
        mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
                snapshot := m.Snapshot()
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(snapshot)
        })
        
        // Handler for HTML dashboard
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                snapshot := m.Snapshot()
                w.Header().Set("Content-Type", "text/html")
                renderHTMLDashboard(w, snapshot)
        })
        
        m.httpServer = &http.Server{
                Addr:    addr,
                Handler: mux,
        }
        
        log.Printf("Starting monitoring server at http://%s", addr)
        return m.httpServer.ListenAndServe()
}

// StopHTTPServer stops the HTTP server
func (m *Monitor) StopHTTPServer(ctx context.Context) error {
        if m.httpServer != nil {
                return m.httpServer.Shutdown(ctx)
        }
        return nil
}

// StartResourceMonitoring begins periodic collection of resource metrics
func (m *Monitor) StartResourceMonitoring() {
        ticker := time.NewTicker(m.tickInterval)
        
        go func() {
                for {
                        select {
                        case <-ticker.C:
                                m.collectResourceMetrics()
                        case <-m.done:
                                ticker.Stop()
                                return
                        }
                }
        }()
}

// StopResourceMonitoring stops the resource monitoring goroutine
func (m *Monitor) StopResourceMonitoring() {
        close(m.done)
}

// collectResourceMetrics collects current resource usage metrics
func (m *Monitor) collectResourceMetrics() {
        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)
        
        m.mu.Lock()
        defer m.mu.Unlock()
        
        m.resourceTicks = append(m.resourceTicks, ResourceMetrics{
                CurrentMemory:  memStats.Alloc,
                PeakMemory:     memStats.Alloc, // Current value, we'll calculate peak separately
                TotalCPUTime:   getCPUTime(),
                GoroutineCount: runtime.NumGoroutine(),
        })
        
        // Limit the number of ticks we store to avoid memory growth
        if len(m.resourceTicks) > 3600 { // Store up to 1 hour of 1-second ticks
                m.resourceTicks = m.resourceTicks[1:]
        }
}

// PrintStatus prints the current execution status to stdout
func (m *Monitor) PrintStatus() {
        snapshot := m.Snapshot()
        
        fmt.Printf("\n=== RunInk DAG Execution Status ===\n")
        fmt.Printf("DAG: %s\n", snapshot.DAGName)
        fmt.Printf("Elapsed Time: %v\n", snapshot.ElapsedTime.Round(time.Millisecond))
        fmt.Printf("Progress: %.1f%% (%d/%d nodes)\n", 
                snapshot.Progress*100, 
                snapshot.CompletedNodes+snapshot.FailedNodes+snapshot.SkippedNodes,
                snapshot.TotalNodes)
        
        fmt.Printf("\nNode Status:\n")
        fmt.Printf("  Completed: %d\n", snapshot.CompletedNodes)
        fmt.Printf("  Running:   %d\n", snapshot.RunningNodes)
        fmt.Printf("  Pending:   %d\n", snapshot.PendingNodes)
        fmt.Printf("  Failed:    %d\n", snapshot.FailedNodes)
        fmt.Printf("  Skipped:   %d\n", snapshot.SkippedNodes)
        fmt.Printf("  Retrying:  %d\n", snapshot.RetryingNodes)
        
        fmt.Printf("\nResource Usage:\n")
        fmt.Printf("  Current Memory: %.2f MB\n", float64(snapshot.ResourceMetrics.CurrentMemory)/1024/1024)
        fmt.Printf("  Peak Memory:    %.2f MB\n", float64(snapshot.ResourceMetrics.PeakMemory)/1024/1024)
        fmt.Printf("  Goroutines:     %d\n", snapshot.ResourceMetrics.GoroutineCount)
        
        // Print recently completed nodes
        fmt.Printf("\nRecent Node Completions:\n")
        count := 0
        for _, metric := range snapshot.NodeMetrics {
                if metric.Status != NodeRunning && metric.Status != NodePending && !metric.EndTime.IsZero() {
                        if count < 5 { // Show only the 5 most recent completions
                                statusStr := string(metric.Status)
                                if metric.Error != "" {
                                        statusStr += " (ERROR)"
                                }
                                fmt.Printf("  %s: %s - %s in %v\n", 
                                        metric.NodeName, 
                                        statusStr,
                                        metric.EndTime.Format("15:04:05"),
                                        metric.Duration.Round(time.Millisecond))
                                count++
                        }
                }
        }
        
        fmt.Printf("\nRunning Nodes:\n")
        for _, metric := range snapshot.NodeMetrics {
                if metric.Status == NodeRunning {
                        fmt.Printf("  %s - running for %v\n", 
                                metric.NodeName,
                                time.Since(metric.StartTime).Round(time.Millisecond))
                }
        }
        fmt.Println("=====================================")
}

// StartPeriodicStatusPrinting starts printing status at regular intervals
func (m *Monitor) StartPeriodicStatusPrinting(interval time.Duration) {
        ticker := time.NewTicker(interval)
        
        go func() {
                for {
                        select {
                        case <-ticker.C:
                                m.PrintStatus()
                        case <-m.done:
                                ticker.Stop()
                                return
                        }
                }
        }()
}

// Helper functions

// getCPUTime is a placeholder for getting CPU time
// In a real implementation, this would use OS-specific methods to get actual CPU time
func getCPUTime() time.Duration {
        // This is a simplified version that doesn't actually measure CPU time
        // In a real implementation, you would use OS-specific methods
        return time.Duration(0)
}

// getPeakMemory returns the peak memory usage from resource ticks
func getPeakMemory(ticks []ResourceMetrics) uint64 {
        var peak uint64
        for _, tick := range ticks {
                if tick.CurrentMemory > peak {
                        peak = tick.CurrentMemory
                }
        }
        return peak
}

// renderHTMLDashboard renders a simple HTML dashboard
func renderHTMLDashboard(w http.ResponseWriter, snapshot MonitorSnapshot) {
        html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>RunInk DAG Execution Monitor</title>
    <meta http-equiv="refresh" content="2">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .progress-bar { 
            width: 100%%; 
            background-color: #f3f3f3; 
            border-radius: 5px; 
            margin: 10px 0; 
        }
        .progress-bar-fill { 
            height: 30px; 
            background-color: #4CAF50; 
            border-radius: 5px; 
            width: %f%%; 
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
        }
        .stats { display: flex; flex-wrap: wrap; }
        .stat-box { 
            background-color: #f9f9f9; 
            border: 1px solid #ddd; 
            border-radius: 5px;
            padding: 10px; 
            margin: 5px; 
            flex: 1; 
            min-width: 200px;
        }
        .node-table { 
            width: 100%%; 
            border-collapse: collapse; 
            margin-top: 20px; 
        }
        .node-table th, .node-table td { 
            border: 1px solid #ddd; 
            padding: 8px; 
            text-align: left; 
        }
        .node-table th { background-color: #f2f2f2; }
        .node-table tr:nth-child(even) { background-color: #f9f9f9; }
        .status-succeeded { color: green; }
        .status-failed { color: red; }
        .status-running { color: blue; }
        .status-pending { color: gray; }
        .status-skipped { color: orange; }
        .status-retrying { color: purple; }
    </style>
</head>
<body>
    <h1>RunInk DAG Execution Monitor</h1>
    <p><strong>DAG:</strong> %s</p>
    <p><strong>Started:</strong> %s (elapsed: %s)</p>
    
    <div class="progress-bar">
        <div class="progress-bar-fill">%.1f%%</div>
    </div>
    
    <div class="stats">
        <div class="stat-box">
            <h3>Node Status</h3>
            <p>Total: %d</p>
            <p>Completed: %d</p>
            <p>Running: %d</p>
            <p>Pending: %d</p>
            <p>Failed: %d</p>
            <p>Skipped: %d</p>
            <p>Retrying: %d</p>
        </div>
        
        <div class="stat-box">
            <h3>Resource Usage</h3>
            <p>Current Memory: %.2f MB</p>
            <p>Peak Memory: %.2f MB</p>
            <p>Goroutines: %d</p>
        </div>
    </div>
    
    <h2>Node Details</h2>
    <table class="node-table">
        <tr>
            <th>Node</th>
            <th>Status</th>
            <th>Start Time</th>
            <th>Duration</th>
            <th>Memory Delta</th>
            <th>Retries</th>
            <th>Error</th>
        </tr>
`,
                snapshot.Progress*100,
                snapshot.DAGName,
                snapshot.StartTime.Format("2006-01-02 15:04:05"),
                snapshot.ElapsedTime.Round(time.Millisecond),
                snapshot.Progress*100,
                snapshot.TotalNodes,
                snapshot.CompletedNodes,
                snapshot.RunningNodes,
                snapshot.PendingNodes,
                snapshot.FailedNodes,
                snapshot.SkippedNodes,
                snapshot.RetryingNodes,
                float64(snapshot.ResourceMetrics.CurrentMemory)/1024/1024,
                float64(snapshot.ResourceMetrics.PeakMemory)/1024/1024,
                snapshot.ResourceMetrics.GoroutineCount,
        )
        
        // Add node rows
        for _, metric := range snapshot.NodeMetrics {
                var statusClass string
                switch metric.Status {
                case NodeSucceeded:
                        statusClass = "status-succeeded"
                case NodeFailed:
                        statusClass = "status-failed"
                case NodeRunning:
                        statusClass = "status-running"
                case NodePending:
                        statusClass = "status-pending"
                case NodeSkipped:
                        statusClass = "status-skipped"
                case NodeRetrying:
                        statusClass = "status-retrying"
                }
                
                var duration string
                if metric.Status == NodeRunning {
                        duration = time.Since(metric.StartTime).Round(time.Millisecond).String()
                } else if !metric.EndTime.IsZero() {
                        duration = metric.Duration.Round(time.Millisecond).String()
                } else {
                        duration = "-"
                }
                
                var memoryDelta string
                if metric.MemoryDelta != 0 {
                        memoryDelta = fmt.Sprintf("%.2f MB", float64(metric.MemoryDelta)/1024/1024)
                } else {
                        memoryDelta = "-"
                }
                
                html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td class="%s">%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%d</td>
            <td>%s</td>
        </tr>
`,
                        metric.NodeName,
                        statusClass, metric.Status,
                        metric.StartTime.Format("15:04:05"),
                        duration,
                        memoryDelta,
                        metric.RetryCount,
                        metric.Error,
                )
        }
        
        html += `
    </table>
    
    <script>
        // Auto-refresh the page every 2 seconds
        setTimeout(function() {
            location.reload();
        }, 2000);
    </script>
</body>
</html>
`
        
        fmt.Fprint(w, html)
}

/*
Example integration with the Execute function:

In execute.go, modify the Execute function to use the monitor:

```go
func Execute(config ExecutionConfig) (*ExecutionResult, error) {
    // Create a new monitor
    monitor := NewMonitor(config.DAG.Name, len(config.DAG.Nodes))
    
    // Start resource monitoring
    monitor.StartResourceMonitoring()
    defer monitor.StopResourceMonitoring()
    
    // Optionally start HTTP server for dashboard
    go func() {
        if err := monitor.StartHTTPServer("localhost:8080"); err != nil && err != http.ErrServerClosed {
            log.Printf("Monitor HTTP server error: %v", err)
        }
    }()
    
    // Optionally start periodic status printing
    monitor.StartPeriodicStatusPrinting(5 * time.Second)
    
    // ... existing code ...
    
    // When executing a node, add monitoring:
    for _, node := range orderedNodes {
        wg.Add(1)
        
        go func(n *dag.Node) {
            defer wg.Done()
            
            // Start monitoring this node
            finishNode := monitor.StartNode(n.ID, n.Name)
            
            // ... existing node execution code ...
            
            // When node completes, call the finish function
            finishNode(nodeStatus, err, retryCount)
            
            // ... rest of node execution code ...
        }(node)
    }
    
    // ... rest of execution code ...
    
    // Before returning, stop the HTTP server
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    monitor.StopHTTPServer(ctx)
    
    return result, nil
}
```
*/
