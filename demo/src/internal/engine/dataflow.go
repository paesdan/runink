// Package engine provides the execution engine for RunInk DAGs.
package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// DataPacket represents a structured data packet that flows between DAG nodes.
// It contains metadata about the data and the actual payload.
type DataPacket struct {
	// Metadata contains information about the data packet
	Metadata map[string]string
	// Payload contains the actual data
	Payload interface{}
	// Timestamp when the packet was created
	Timestamp time.Time
	// SourceNode is the ID of the node that created this packet
	SourceNode string
}

// NewDataPacket creates a new data packet with the given payload and source node ID.
func NewDataPacket(payload interface{}, sourceNodeID string) *DataPacket {
	return &DataPacket{
		Metadata:   make(map[string]string),
		Payload:    payload,
		Timestamp:  time.Now(),
		SourceNode: sourceNodeID,
	}
}

// WithMetadata adds metadata to the data packet and returns the packet for chaining.
func (dp *DataPacket) WithMetadata(key, value string) *DataPacket {
	dp.Metadata[key] = value
	return dp
}

// Codec defines an interface for encoding and decoding data.
type Codec interface {
	// Encode serializes data into bytes
	Encode(data interface{}) ([]byte, error)
	// Decode deserializes bytes into the provided data structure
	Decode(bytes []byte, target interface{}) error
	// Name returns the name of the codec
	Name() string
}

// JSONCodec implements the Codec interface using JSON serialization.
type JSONCodec struct{}

// Encode serializes data into JSON bytes.
func (c JSONCodec) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// Decode deserializes JSON bytes into the provided data structure.
func (c JSONCodec) Decode(bytes []byte, target interface{}) error {
	return json.Unmarshal(bytes, target)
}

// Name returns the name of the codec.
func (c JSONCodec) Name() string {
	return "json"
}

// DefaultCodec is the default codec used for serialization/deserialization.
var DefaultCodec Codec = JSONCodec{}

// CodecRegistry stores registered codecs by name.
var CodecRegistry = map[string]Codec{
	"json": DefaultCodec,
}

// RegisterCodec registers a codec with the given name.
func RegisterCodec(name string, codec Codec) {
	CodecRegistry[name] = codec
}

// GetCodec returns a codec by name, or the default codec if not found.
func GetCodec(name string) Codec {
	if codec, exists := CodecRegistry[name]; exists {
		return codec
	}
	return DefaultCodec
}

// DataEdge extends the basic Edge with data flow configuration.
type DataEdge struct {
	// From is the ID of the source node
	From string
	// To is the ID of the target node
	To string
	// BufferSize is the size of the channel buffer (0 for unbuffered)
	BufferSize int
	// CodecName is the name of the codec to use for serialization
	CodecName string
	// RetryPolicy defines how to handle failures in data transmission
	RetryPolicy RetryPolicy
}

// RetryPolicy defines how to handle failures in data transmission.
type RetryPolicy struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int
	// RetryDelay is the delay between retries
	RetryDelay time.Duration
	// BackoffFactor is the factor by which to increase the delay after each retry
	BackoffFactor float64
}

// DefaultRetryPolicy is the default retry policy.
var DefaultRetryPolicy = RetryPolicy{
	MaxRetries:    3,
	RetryDelay:    time.Second,
	BackoffFactor: 2.0,
}

// ChannelManager manages channels between nodes.
type ChannelManager struct {
	// channels maps edge IDs to channels
	channels map[string]chan *DataPacket
	// edges maps edge IDs to edge configurations
	edges map[string]*DataEdge
	// mutex protects the maps
	mutex sync.RWMutex
}

// NewChannelManager creates a new channel manager.
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]chan *DataPacket),
		edges:    make(map[string]*DataEdge),
	}
}

// RegisterEdge registers an edge with the channel manager.
func (cm *ChannelManager) RegisterEdge(edge *DataEdge) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	edgeID := fmt.Sprintf("%s->%s", edge.From, edge.To)
	cm.edges[edgeID] = edge

	// Create a channel with the specified buffer size
	cm.channels[edgeID] = make(chan *DataPacket, edge.BufferSize)
}

// GetChannel returns the channel for the given edge.
func (cm *ChannelManager) GetChannel(fromNodeID, toNodeID string) (chan *DataPacket, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	edgeID := fmt.Sprintf("%s->%s", fromNodeID, toNodeID)
	channel, exists := cm.channels[edgeID]
	if !exists {
		return nil, fmt.Errorf("no channel exists for edge %s", edgeID)
	}

	return channel, nil
}

// GetEdge returns the edge configuration for the given edge.
func (cm *ChannelManager) GetEdge(fromNodeID, toNodeID string) (*DataEdge, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	edgeID := fmt.Sprintf("%s->%s", fromNodeID, toNodeID)
	edge, exists := cm.edges[edgeID]
	if !exists {
		return nil, fmt.Errorf("no edge exists for %s", edgeID)
	}

	return edge, nil
}

// CloseAllChannels closes all channels managed by the channel manager.
func (cm *ChannelManager) CloseAllChannels() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, ch := range cm.channels {
		close(ch)
	}
}

// Send sends data from a source node to a target node.
// It handles serialization and retries according to the edge's retry policy.
func Send(cm *ChannelManager, fromNodeID, toNodeID string, data interface{}) error {
	// Get the channel for this edge
	ch, err := cm.GetChannel(fromNodeID, toNodeID)
	if err != nil {
		return err
	}

	// Get the edge configuration
	edge, err := cm.GetEdge(fromNodeID, toNodeID)
	if err != nil {
		return err
	}

	// Create a data packet
	packet := NewDataPacket(data, fromNodeID)

	// Send the packet to the channel
	select {
	case ch <- packet:
		return nil
	default:
		// Channel is full, apply retry policy
		return sendWithRetry(ch, packet, edge.RetryPolicy)
	}
}

// sendWithRetry attempts to send a packet with retries according to the retry policy.
func sendWithRetry(ch chan *DataPacket, packet *DataPacket, policy RetryPolicy) error {
	delay := policy.RetryDelay

	for i := 0; i < policy.MaxRetries; i++ {
		// Wait before retrying
		time.Sleep(delay)

		// Try to send again
		select {
		case ch <- packet:
			return nil
		default:
			// Increase delay for next retry
			delay = time.Duration(float64(delay) * policy.BackoffFactor)
		}
	}

	return errors.New("failed to send data after maximum retries")
}

// Receive receives data from a channel and decodes it into the target.
// The target must be a pointer to the expected data type.
func Receive[T any](cm *ChannelManager, fromNodeID, toNodeID string, target *T) error {
	// Get the channel for this edge
	ch, err := cm.GetChannel(fromNodeID, toNodeID)
	if err != nil {
		return err
	}

	// Receive the packet from the channel
	packet, ok := <-ch
	if !ok {
		return errors.New("channel is closed")
	}

	// Check if the payload is already of the target type
	if payload, ok := packet.Payload.(*T); ok {
		*target = *payload
		return nil
	}

	// Try to convert the payload to the target type
	if reflect.TypeOf(packet.Payload).AssignableTo(reflect.TypeOf(*target)) {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(packet.Payload))
		return nil
	}

	// Get the edge configuration
	edge, err := cm.GetEdge(fromNodeID, toNodeID)
	if err != nil {
		return err
	}

	// Get the codec for this edge
	codec := GetCodec(edge.CodecName)

	// Encode and then decode to convert between types
	bytes, err := codec.Encode(packet.Payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	err = codec.Decode(bytes, target)
	if err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	return nil
}

// ReceiveWithTimeout receives data from a channel with a timeout.
// Returns an error if the timeout expires before data is available.
func ReceiveWithTimeout[T any](cm *ChannelManager, fromNodeID, toNodeID string, target *T, timeout time.Duration) error {
	// Get the channel for this edge
	ch, err := cm.GetChannel(fromNodeID, toNodeID)
	if err != nil {
		return err
	}

	// Set up a timeout
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// Wait for data or timeout
	select {
	case packet, ok := <-ch:
		if !ok {
			return errors.New("channel is closed")
		}

		// Check if the payload is already of the target type
		if payload, ok := packet.Payload.(*T); ok {
			*target = *payload
			return nil
		}

		// Try to convert the payload to the target type
		if reflect.TypeOf(packet.Payload).AssignableTo(reflect.TypeOf(*target)) {
			reflect.ValueOf(target).Elem().Set(reflect.ValueOf(packet.Payload))
			return nil
		}

		// Get the edge configuration
		edge, err := cm.GetEdge(fromNodeID, toNodeID)
		if err != nil {
			return err
		}

		// Get the codec for this edge
		codec := GetCodec(edge.CodecName)

		// Encode and then decode to convert between types
		bytes, err := codec.Encode(packet.Payload)
		if err != nil {
			return fmt.Errorf("failed to encode payload: %w", err)
		}

		err = codec.Decode(bytes, target)
		if err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}

		return nil

	case <-timer.C:
		return errors.New("timeout waiting for data")
	}
}

// TryReceive attempts to receive data from a channel without blocking.
// Returns false if no data is available.
func TryReceive[T any](cm *ChannelManager, fromNodeID, toNodeID string, target *T) (bool, error) {
	// Get the channel for this edge
	ch, err := cm.GetChannel(fromNodeID, toNodeID)
	if err != nil {
		return false, err
	}

	// Try to receive without blocking
	select {
	case packet, ok := <-ch:
		if !ok {
			return false, errors.New("channel is closed")
		}

		// Check if the payload is already of the target type
		if payload, ok := packet.Payload.(*T); ok {
			*target = *payload
			return true, nil
		}

		// Try to convert the payload to the target type
		if reflect.TypeOf(packet.Payload).AssignableTo(reflect.TypeOf(*target)) {
			reflect.ValueOf(target).Elem().Set(reflect.ValueOf(packet.Payload))
			return true, nil
		}

		// Get the edge configuration
		edge, err := cm.GetEdge(fromNodeID, toNodeID)
		if err != nil {
			return false, err
		}

		// Get the codec for this edge
		codec := GetCodec(edge.CodecName)

		// Encode and then decode to convert between types
		bytes, err := codec.Encode(packet.Payload)
		if err != nil {
			return false, fmt.Errorf("failed to encode payload: %w", err)
		}

		err = codec.Decode(bytes, target)
		if err != nil {
			return false, fmt.Errorf("failed to decode payload: %w", err)
		}

		return true, nil

	default:
		return false, nil
	}
}

// DataFlowExecutor extends the basic executor with data flow capabilities.
type DataFlowExecutor struct {
	// ChannelManager manages channels between nodes
	ChannelManager *ChannelManager
	// ExecutionConfig holds the configuration for DAG execution
	Config ExecutionConfig
}

// NewDataFlowExecutor creates a new data flow executor.
func NewDataFlowExecutor(config ExecutionConfig) *DataFlowExecutor {
	return &DataFlowExecutor{
		ChannelManager: NewChannelManager(),
		Config:         config,
	}
}

// SetupDataFlow initializes the data flow for the DAG.
// It creates channels for all edges with appropriate buffer sizes.
func (e *DataFlowExecutor) SetupDataFlow() error {
	// Create data edges from the DAG edges
	for _, edge := range e.Config.DAG.Edges {
		// Create a data edge with default settings
		dataEdge := &DataEdge{
			From:        edge.From,
			To:          edge.To,
			BufferSize:  1, // Default buffer size
			CodecName:   "json",
			RetryPolicy: DefaultRetryPolicy,
		}

		// Check if the edge has custom buffer size in the node config
		fromNode := e.Config.DAG.GetNode(edge.From)
		if fromNode != nil {
			if bufferSize, ok := fromNode.Config["buffer_size"]; ok {
				if bufferSizeInt, ok := bufferSize.(int); ok {
					dataEdge.BufferSize = bufferSizeInt
				} else if bufferSizeStr, ok := bufferSize.(string); ok {
					fmt.Sscanf(bufferSizeStr, "%d", &dataEdge.BufferSize)
				}
			}

			// Check if the edge has custom codec in the node config
			if codecName, ok := fromNode.Config["codec"]; ok {
				if codecNameStr, ok := codecName.(string); ok {
					dataEdge.CodecName = codecNameStr
				}
			}
		}

		// Register the edge with the channel manager
		e.ChannelManager.RegisterEdge(dataEdge)
	}

	return nil
}

// ExecuteWithDataFlow executes the DAG with data flow between nodes.
// It extends the basic Execute function with typed data passing.
func (e *DataFlowExecutor) ExecuteWithDataFlow() (*ExecutionResult, error) {
	// Set up data flow
	if err := e.SetupDataFlow(); err != nil {
		return nil, err
	}

	// Create a modified execution config that uses our data flow
	config := e.Config

	// Execute the DAG
	result, err := Execute(config)

	// Clean up channels
	e.ChannelManager.CloseAllChannels()

	return result, err
}

// BatchProcessor provides utilities for batch processing of data.
type BatchProcessor struct {
	// BatchSize is the number of items to process in a batch
	BatchSize int
	// Timeout is the maximum time to wait for a batch to fill
	Timeout time.Duration
	// buffer holds the current batch
	buffer []interface{}
	// mutex protects the buffer
	mutex sync.Mutex
}

// NewBatchProcessor creates a new batch processor.
func NewBatchProcessor(batchSize int, timeout time.Duration) *BatchProcessor {
	return &BatchProcessor{
		BatchSize: batchSize,
		Timeout:   timeout,
		buffer:    make([]interface{}, 0, batchSize),
	}
}

// Add adds an item to the batch.
// Returns true if the batch is full and ready to be processed.
func (bp *BatchProcessor) Add(item interface{}) bool {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	bp.buffer = append(bp.buffer, item)
	return len(bp.buffer) >= bp.BatchSize
}

// GetBatch returns the current batch and resets the buffer.
func (bp *BatchProcessor) GetBatch() []interface{} {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	batch := bp.buffer
	bp.buffer = make([]interface{}, 0, bp.BatchSize)
	return batch
}

// ProcessBatch processes a batch of data with the given function.
// The function is called for each item in the batch.
func ProcessBatch(batch []interface{}, processor func(interface{}) error) error {
	for _, item := range batch {
		if err := processor(item); err != nil {
			return err
		}
	}
	return nil
}

// DataTransformer defines an interface for transforming data between nodes.
type DataTransformer interface {
	// Transform transforms data from one type to another
	Transform(input interface{}) (interface{}, error)
}

// FunctionTransformer implements DataTransformer using a function.
type FunctionTransformer struct {
	// TransformFunc is the function to use for transformation
	TransformFunc func(interface{}) (interface{}, error)
}

// Transform transforms data using the transform function.
func (ft FunctionTransformer) Transform(input interface{}) (interface{}, error) {
	return ft.TransformFunc(input)
}

// NewFunctionTransformer creates a new function transformer.
func NewFunctionTransformer(fn func(interface{}) (interface{}, error)) *FunctionTransformer {
	return &FunctionTransformer{
		TransformFunc: fn,
	}
}

// TransformAndSend transforms data and sends it to the target node.
func TransformAndSend(cm *ChannelManager, fromNodeID, toNodeID string, data interface{}, transformer DataTransformer) error {
	// Transform the data
	transformed, err := transformer.Transform(data)
	if err != nil {
		return err
	}

	// Send the transformed data
	return Send(cm, fromNodeID, toNodeID, transformed)
}
