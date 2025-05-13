// Package engine provides the execution engine for RunInk DAGs.
package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ExecutionMode defines the mode in which a node should be executed.
type ExecutionMode int

const (
	// Batch mode processes all input data at once and produces output after processing is complete.
	Batch ExecutionMode = iota
	
	// Streaming mode processes each input item as it arrives and produces output continuously.
	Streaming
)

// String returns a string representation of the execution mode.
func (m ExecutionMode) String() string {
	switch m {
	case Batch:
		return "Batch"
	case Streaming:
		return "Streaming"
	default:
		return "Unknown"
	}
}

// ModeHandler defines the interface for execution mode handlers.
// Each handler implements a specific execution strategy for processing data.
type ModeHandler interface {
	// Run executes the node's logic with the given input and output channels.
	// The context can be used to signal cancellation.
	Run(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket, processor func(*DataPacket) (*DataPacket, error)) error
}

// BatchHandler implements the ModeHandler interface for batch processing.
// It collects all input data before processing and then produces output.
type BatchHandler struct {
	// MaxBatchSize is the maximum number of items to process in a batch.
	// If set to 0 or negative, all available items will be processed in a single batch.
	MaxBatchSize int
	
	// Timeout is the maximum time to wait for input before processing the current batch.
	// If set to 0, it will wait indefinitely until all input is received.
	Timeout time.Duration
}

// NewBatchHandler creates a new batch handler with the given configuration.
func NewBatchHandler(maxBatchSize int, timeout time.Duration) *BatchHandler {
	return &BatchHandler{
		MaxBatchSize: maxBatchSize,
		Timeout:      timeout,
	}
}

// Run implements the ModeHandler interface for batch processing.
// It collects all input data, processes it as a batch, and then produces output.
func (h *BatchHandler) Run(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket, processor func(*DataPacket) (*DataPacket, error)) error {
	// Collect all input data
	var batch []*DataPacket
	var timer *time.Timer
	var timerCh <-chan time.Time
	
	if h.Timeout > 0 {
		timer = time.NewTimer(h.Timeout)
		timerCh = timer.C
		defer timer.Stop()
	}
	
	// Function to process the current batch
	processBatch := func() error {
		if len(batch) == 0 {
			return nil
		}
		
		// Process each item in the batch
		for _, packet := range batch {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Process the packet
				result, err := processor(packet)
				if err != nil {
					return err
				}
				
				// Send the result if not nil
				if result != nil {
					select {
					case out <- result:
						// Successfully sent
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
		}
		
		// Clear the batch
		batch = nil
		return nil
	}
	
	for {
		select {
		case <-ctx.Done():
			// Context canceled, process any remaining items and exit
			if err := processBatch(); err != nil {
				return err
			}
			return ctx.Err()
			
		case packet, ok := <-in:
			if !ok {
				// Input channel closed, process any remaining items and exit
				if err := processBatch(); err != nil {
					return err
				}
				return nil
			}
			
			// Add the packet to the batch
			batch = append(batch, packet)
			
			// Process the batch if it reaches the maximum size
			if h.MaxBatchSize > 0 && len(batch) >= h.MaxBatchSize {
				if err := processBatch(); err != nil {
					return err
				}
				
				// Reset the timer if it's active
				if timer != nil {
					timer.Reset(h.Timeout)
				}
			}
			
		case <-timerCh:
			// Timeout reached, process the current batch
			if err := processBatch(); err != nil {
				return err
			}
			
			// Reset the timer
			timer.Reset(h.Timeout)
		}
	}
}

// StreamingHandler implements the ModeHandler interface for streaming processing.
// It processes each input item as it arrives and produces output continuously.
type StreamingHandler struct {
	// BufferSize is the size of the internal buffer for processing.
	// A larger buffer can improve throughput but may increase latency.
	BufferSize int
}

// NewStreamingHandler creates a new streaming handler with the given configuration.
func NewStreamingHandler(bufferSize int) *StreamingHandler {
	return &StreamingHandler{
		BufferSize: bufferSize,
	}
}

// Run implements the ModeHandler interface for streaming processing.
// It processes each input item as it arrives and produces output continuously.
func (h *StreamingHandler) Run(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket, processor func(*DataPacket) (*DataPacket, error)) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
			
		case packet, ok := <-in:
			if !ok {
				// Input channel closed, exit
				return nil
			}
			
			// Process the packet
			result, err := processor(packet)
			if err != nil {
				return err
			}
			
			// Send the result if not nil
			if result != nil {
				select {
				case out <- result:
					// Successfully sent
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}

// ModeConfig holds configuration for execution mode.
type ModeConfig struct {
	// Mode is the execution mode to use.
	Mode ExecutionMode
	
	// BatchSize is the maximum number of items to process in a batch.
	// Only used in Batch mode.
	BatchSize int
	
	// BatchTimeout is the maximum time to wait for input before processing the current batch.
	// Only used in Batch mode.
	BatchTimeout time.Duration
	
	// StreamingBufferSize is the size of the internal buffer for streaming processing.
	// Only used in Streaming mode.
	StreamingBufferSize int
}

// DefaultModeConfig returns the default mode configuration.
func DefaultModeConfig() ModeConfig {
	return ModeConfig{
		Mode:                Batch,
		BatchSize:           100,
		BatchTimeout:        time.Second * 5,
		StreamingBufferSize: 10,
	}
}

// GetModeHandler returns the appropriate mode handler for the given mode configuration.
func GetModeHandler(config ModeConfig) ModeHandler {
	switch config.Mode {
	case Batch:
		return NewBatchHandler(config.BatchSize, config.BatchTimeout)
	case Streaming:
		return NewStreamingHandler(config.StreamingBufferSize)
	default:
		// Default to batch mode
		return NewBatchHandler(config.BatchSize, config.BatchTimeout)
	}
}

// NodeWithMode extends the basic Node with execution mode capabilities.
type NodeWithMode struct {
	// Node is the underlying DAG node.
	Node *dag.Node
	
	// ModeConfig is the execution mode configuration for this node.
	ModeConfig ModeConfig
	
	// Processor is the function that processes data packets.
	Processor func(*DataPacket) (*DataPacket, error)
}

// NewNodeWithMode creates a new node with the given execution mode.
func NewNodeWithMode(node *dag.Node, mode ExecutionMode) *NodeWithMode {
	config := DefaultModeConfig()
	config.Mode = mode
	
	return &NodeWithMode{
		Node:       node,
		ModeConfig: config,
	}
}

// WithProcessor sets the processor function for this node.
func (n *NodeWithMode) WithProcessor(processor func(*DataPacket) (*DataPacket, error)) *NodeWithMode {
	n.Processor = processor
	return n
}

// WithBatchConfig sets the batch configuration for this node.
func (n *NodeWithMode) WithBatchConfig(batchSize int, timeout time.Duration) *NodeWithMode {
	n.ModeConfig.BatchSize = batchSize
	n.ModeConfig.BatchTimeout = timeout
	return n
}

// WithStreamingConfig sets the streaming configuration for this node.
func (n *NodeWithMode) WithStreamingConfig(bufferSize int) *NodeWithMode {
	n.ModeConfig.StreamingBufferSize = bufferSize
	return n
}

// Execute executes this node with the given input and output channels.
func (n *NodeWithMode) Execute(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket) error {
	if n.Processor == nil {
		return errors.New("processor function not set")
	}
	
	handler := GetModeHandler(n.ModeConfig)
	return handler.Run(ctx, in, out, n.Processor)
}

// ModeTransition handles the transition between different execution modes.
// It can be used to connect nodes with different execution modes.
type ModeTransition struct {
	// FromMode is the execution mode of the source node.
	FromMode ExecutionMode
	
	// ToMode is the execution mode of the target node.
	ToMode ExecutionMode
	
	// BufferSize is the size of the buffer for the transition.
	BufferSize int
}

// NewModeTransition creates a new mode transition.
func NewModeTransition(fromMode, toMode ExecutionMode, bufferSize int) *ModeTransition {
	return &ModeTransition{
		FromMode:   fromMode,
		ToMode:     toMode,
		BufferSize: bufferSize,
	}
}

// CreateTransitionChannel creates a channel for the transition between modes.
func (t *ModeTransition) CreateTransitionChannel() chan *DataPacket {
	return make(chan *DataPacket, t.BufferSize)
}

// Run runs the transition between modes.
// It reads from the input channel and writes to the output channel,
// applying any necessary transformations between modes.
func (t *ModeTransition) Run(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket) error {
	// If the modes are the same, just pass through
	if t.FromMode == t.ToMode {
		return passThrough(ctx, in, out)
	}
	
	// Handle specific transitions
	switch {
	case t.FromMode == Batch && t.ToMode == Streaming:
		// Batch to Streaming: Split batch into individual items
		return batchToStreaming(ctx, in, out)
		
	case t.FromMode == Streaming && t.ToMode == Batch:
		// Streaming to Batch: Collect items into batches
		return streamingToBatch(ctx, in, out, t.BufferSize)
		
	default:
		// Unknown transition, just pass through
		return passThrough(ctx, in, out)
	}
}

// passThrough passes data from input to output without modification.
func passThrough(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
			
		case packet, ok := <-in:
			if !ok {
				// Input channel closed, exit
				return nil
			}
			
			// Pass through to output
			select {
			case out <- packet:
				// Successfully sent
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// batchToStreaming splits batch data into individual items for streaming.
func batchToStreaming(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
			
		case packet, ok := <-in:
			if !ok {
				// Input channel closed, exit
				return nil
			}
			
			// Check if the payload is a slice or array
			if items, ok := packet.Payload.([]interface{}); ok {
				// Split the batch into individual items
				for _, item := range items {
					// Create a new packet for each item
					itemPacket := &DataPacket{
						Metadata:   packet.Metadata,
						Payload:    item,
						Timestamp:  time.Now(),
						SourceNode: packet.SourceNode,
					}
					
					// Send the item
					select {
					case out <- itemPacket:
						// Successfully sent
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			} else {
				// Not a batch, just pass through
				select {
				case out <- packet:
					// Successfully sent
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}

// streamingToBatch collects streaming items into batches.
func streamingToBatch(ctx context.Context, in <-chan *DataPacket, out chan<- *DataPacket, batchSize int) error {
	var batch []*DataPacket
	
	for {
		select {
		case <-ctx.Done():
			// Context canceled, send any remaining items and exit
			if len(batch) > 0 {
				sendBatch(ctx, batch, out)
			}
			return ctx.Err()
			
		case packet, ok := <-in:
			if !ok {
				// Input channel closed, send any remaining items and exit
				if len(batch) > 0 {
					sendBatch(ctx, batch, out)
				}
				return nil
			}
			
			// Add the packet to the batch
			batch = append(batch, packet)
			
			// Send the batch if it reaches the maximum size
			if len(batch) >= batchSize {
				if err := sendBatch(ctx, batch, out); err != nil {
					return err
				}
				batch = nil
			}
		}
	}
}

// sendBatch sends a batch of packets as a single packet.
func sendBatch(ctx context.Context, batch []*DataPacket, out chan<- *DataPacket) error {
	if len(batch) == 0 {
		return nil
	}
	
	// Extract payloads from the batch
	payloads := make([]interface{}, len(batch))
	for i, packet := range batch {
		payloads[i] = packet.Payload
	}
	
	// Create a new packet with the batch of payloads
	batchPacket := &DataPacket{
		Metadata:   batch[0].Metadata,
		Payload:    payloads,
		Timestamp:  time.Now(),
		SourceNode: batch[0].SourceNode,
	}
	
	// Send the batch
	select {
	case out <- batchPacket:
		// Successfully sent
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ModeAwareExecutor extends the DataFlowExecutor with mode-aware execution.
type ModeAwareExecutor struct {
	// DataFlowExecutor is the underlying executor.
	*DataFlowExecutor
	
	// NodeModes maps node IDs to their execution modes.
	NodeModes map[string]ModeConfig
	
	// Transitions maps edge IDs to mode transitions.
	Transitions map[string]*ModeTransition
	
	// TransitionChannels maps edge IDs to transition channels.
	TransitionChannels map[string]chan *DataPacket
	
	// mutex protects the maps.
	mutex sync.RWMutex
}

// NewModeAwareExecutor creates a new mode-aware executor.
func NewModeAwareExecutor(config ExecutionConfig) *ModeAwareExecutor {
	return &ModeAwareExecutor{
		DataFlowExecutor:   NewDataFlowExecutor(config),
		NodeModes:          make(map[string]ModeConfig),
		Transitions:        make(map[string]*ModeTransition),
		TransitionChannels: make(map[string]chan *DataPacket),
	}
}

// SetNodeMode sets the execution mode for a node.
func (e *ModeAwareExecutor) SetNodeMode(nodeID string, config ModeConfig) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.NodeModes[nodeID] = config
}

// GetNodeMode returns the execution mode for a node.
// If the node doesn't have a specific mode set, it returns the default mode.
func (e *ModeAwareExecutor) GetNodeMode(nodeID string) ModeConfig {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	if mode, exists := e.NodeModes[nodeID]; exists {
		return mode
	}
	
	// Check if the node has a mode specified in its config
	node := e.Config.DAG.GetNode(nodeID)
	if node != nil {
		if modeStr, ok := node.Config["execution_mode"]; ok {
			if modeStr == "streaming" {
				config := DefaultModeConfig()
				config.Mode = Streaming
				return config
			}
		}
	}
	
	return DefaultModeConfig()
}

// SetupModeTransitions sets up the transitions between nodes with different modes.
func (e *ModeAwareExecutor) SetupModeTransitions() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	// Create transitions for all edges
	for _, edge := range e.Config.DAG.Edges {
		fromMode := e.GetNodeMode(edge.From).Mode
		toMode := e.GetNodeMode(edge.To).Mode
		
		// Create a transition if the modes are different
		if fromMode != toMode {
			edgeID := fmt.Sprintf("%s->%s", edge.From, edge.To)
			
			// Default buffer size
			bufferSize := 10
			
			// Check if the edge has a custom buffer size
			fromNode := e.Config.DAG.GetNode(edge.From)
			if fromNode != nil {
				if bufferSizeConfig, ok := fromNode.Config["buffer_size"]; ok {
					if bufferSizeInt, ok := bufferSizeConfig.(int); ok {
						bufferSize = bufferSizeInt
					} else if bufferSizeStr, ok := bufferSizeConfig.(string); ok {
						fmt.Sscanf(bufferSizeStr, "%d", &bufferSize)
					}
				}
			}
			
			// Create the transition
			transition := NewModeTransition(fromMode, toMode, bufferSize)
			e.Transitions[edgeID] = transition
			
			// Create the transition channel
			e.TransitionChannels[edgeID] = transition.CreateTransitionChannel()
		}
	}
	
	return nil
}

// ExecuteWithModeAwareness executes the DAG with mode-aware execution.
func (e *ModeAwareExecutor) ExecuteWithModeAwareness() (*ExecutionResult, error) {
	// Set up data flow
	if err := e.SetupDataFlow(); err != nil {
		return nil, err
	}
	
	// Set up mode transitions
	if err := e.SetupModeTransitions(); err != nil {
		return nil, err
	}
	
	// Create a context for execution
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Start transition goroutines
	var wg sync.WaitGroup
	for edgeID, transition := range e.Transitions {
		wg.Add(1)
		
		go func(id string, t *ModeTransition) {
			defer wg.Done()
			
			// Get the source and target channels
			parts := splitEdgeID(id)
			if len(parts) != 2 {
				return
			}
			
			fromNodeID, toNodeID := parts[0], parts[1]
			
			// Get the channels
			sourceChannel, err := e.ChannelManager.GetChannel(fromNodeID, toNodeID)
			if err != nil {
				return
			}
			
			targetChannel := e.TransitionChannels[id]
			
			// Run the transition
			t.Run(ctx, sourceChannel, targetChannel)
		}(edgeID, transition)
	}
	
	// Execute the DAG
	result, err := Execute(e.Config)
	
	// Wait for all transitions to complete
	cancel()
	wg.Wait()
	
	// Clean up channels
	e.ChannelManager.CloseAllChannels()
	for _, ch := range e.TransitionChannels {
		close(ch)
	}
	
	return result, err
}

// splitEdgeID splits an edge ID into source and target node IDs.
func splitEdgeID(edgeID string) []string {
	var fromNodeID, toNodeID string
	fmt.Sscanf(edgeID, "%s->%s", &fromNodeID, &toNodeID)
	return []string{fromNodeID, toNodeID}
}

// GetModeFromString converts a string mode name to an ExecutionMode.
func GetModeFromString(mode string) ExecutionMode {
	switch mode {
	case "streaming":
		return Streaming
	case "batch":
		return Batch
	default:
		return Batch
	}
}

// GetModeString converts an ExecutionMode to a string.
func GetModeString(mode ExecutionMode) string {
	switch mode {
	case Streaming:
		return "streaming"
	case Batch:
		return "batch"
	default:
		return "batch"
	}
}

// NodeModeExtension extends a node with execution mode capabilities.
// This can be used to add mode information to existing nodes.
type NodeModeExtension struct {
	// NodeID is the ID of the node to extend.
	NodeID string
	
	// Mode is the execution mode for this node.
	Mode ExecutionMode
	
	// Config is the mode-specific configuration.
	Config ModeConfig
}

// NewNodeModeExtension creates a new node mode extension.
func NewNodeModeExtension(nodeID string, mode ExecutionMode) *NodeModeExtension {
	config := DefaultModeConfig()
	config.Mode = mode
	
	return &NodeModeExtension{
		NodeID: nodeID,
		Mode:   mode,
		Config: config,
	}
}

// WithBatchConfig sets the batch configuration for this node extension.
func (e *NodeModeExtension) WithBatchConfig(batchSize int, timeout time.Duration) *NodeModeExtension {
	e.Config.BatchSize = batchSize
	e.Config.BatchTimeout = timeout
	return e
}

// WithStreamingConfig sets the streaming configuration for this node extension.
func (e *NodeModeExtension) WithStreamingConfig(bufferSize int) *NodeModeExtension {
	e.Config.StreamingBufferSize = bufferSize
	return e
}

// Apply applies this extension to a mode-aware executor.
func (e *NodeModeExtension) Apply(executor *ModeAwareExecutor) {
	executor.SetNodeMode(e.NodeID, e.Config)
}

// ApplyNodeModeExtensions applies a list of node mode extensions to an executor.
func ApplyNodeModeExtensions(executor *ModeAwareExecutor, extensions []*NodeModeExtension) {
	for _, ext := range extensions {
		ext.Apply(executor)
	}
}
