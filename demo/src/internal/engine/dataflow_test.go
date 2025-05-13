package engine

import (
	"testing"
	"time"
)

func TestDataPacket(t *testing.T) {
	// Create a data packet
	packet := NewDataPacket("test data", "node1")
	
	// Add metadata
	packet.WithMetadata("key1", "value1").WithMetadata("key2", "value2")
	
	// Check values
	if packet.Payload != "test data" {
		t.Errorf("Expected payload 'test data', got '%v'", packet.Payload)
	}
	
	if packet.SourceNode != "node1" {
		t.Errorf("Expected source node 'node1', got '%s'", packet.SourceNode)
	}
	
	if packet.Metadata["key1"] != "value1" || packet.Metadata["key2"] != "value2" {
		t.Errorf("Metadata not set correctly")
	}
}

func TestJSONCodec(t *testing.T) {
	codec := JSONCodec{}
	
	// Test data
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	
	original := TestStruct{Name: "test", Value: 42}
	
	// Encode
	bytes, err := codec.Encode(original)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}
	
	// Decode
	var decoded TestStruct
	err = codec.Decode(bytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}
	
	// Check values
	if decoded.Name != original.Name || decoded.Value != original.Value {
		t.Errorf("Decoded value doesn't match original: %+v vs %+v", decoded, original)
	}
}

func TestChannelManager(t *testing.T) {
	cm := NewChannelManager()
	
	// Register an edge
	edge := &DataEdge{
		From:       "node1",
		To:         "node2",
		BufferSize: 5,
		CodecName:  "json",
		RetryPolicy: RetryPolicy{
			MaxRetries:    3,
			RetryDelay:    time.Millisecond * 10,
			BackoffFactor: 2.0,
		},
	}
	
	cm.RegisterEdge(edge)
	
	// Get the channel
	ch, err := cm.GetChannel("node1", "node2")
	if err != nil {
		t.Fatalf("Failed to get channel: %v", err)
	}
	
	// Check buffer size
	if cap(ch) != 5 {
		t.Errorf("Expected buffer size 5, got %d", cap(ch))
	}
	
	// Get the edge
	retrievedEdge, err := cm.GetEdge("node1", "node2")
	if err != nil {
		t.Fatalf("Failed to get edge: %v", err)
	}
	
	if retrievedEdge.From != "node1" || retrievedEdge.To != "node2" {
		t.Errorf("Retrieved edge doesn't match original")
	}
	
	// Test non-existent edge
	_, err = cm.GetEdge("node1", "node3")
	if err == nil {
		t.Errorf("Expected error for non-existent edge, got nil")
	}
}

func TestSendReceive(t *testing.T) {
	cm := NewChannelManager()
	
	// Register an edge
	edge := &DataEdge{
		From:       "node1",
		To:         "node2",
		BufferSize: 5,
		CodecName:  "json",
		RetryPolicy: RetryPolicy{
			MaxRetries:    3,
			RetryDelay:    time.Millisecond * 10,
			BackoffFactor: 2.0,
		},
	}
	
	cm.RegisterEdge(edge)
	
	// Get the channel
	ch, _ := cm.GetChannel("node1", "node2")
	
	// Send data directly to the channel
	data := "test data"
	packet := NewDataPacket(data, "node1")
	ch <- packet
	
	// Receive data
	var received string
	err := Receive(cm, "node1", "node2", &received)
	if err != nil {
		t.Fatalf("Failed to receive data: %v", err)
	}
	
	if received != data {
		t.Errorf("Expected '%s', got '%s'", data, received)
	}
}

func TestBatchProcessor(t *testing.T) {
	bp := NewBatchProcessor(3, time.Second)
	
	// Add items
	if bp.Add("item1") {
		t.Errorf("Batch should not be full after adding 1 item")
	}
	
	if bp.Add("item2") {
		t.Errorf("Batch should not be full after adding 2 items")
	}
	
	if !bp.Add("item3") {
		t.Errorf("Batch should be full after adding 3 items")
	}
	
	// Get batch
	batch := bp.GetBatch()
	if len(batch) != 3 {
		t.Errorf("Expected batch size 3, got %d", len(batch))
	}
	
	// Check batch is empty after getting
	if len(bp.GetBatch()) != 0 {
		t.Errorf("Batch should be empty after getting")
	}
}

func TestDataTransformer(t *testing.T) {
	// Create a transformer that doubles integers
	transformer := NewFunctionTransformer(func(input interface{}) (interface{}, error) {
		if val, ok := input.(int); ok {
			return val * 2, nil
		}
		return input, nil
	})
	
	// Transform data
	result, err := transformer.Transform(5)
	if err != nil {
		t.Fatalf("Failed to transform data: %v", err)
	}
	
	if result != 10 {
		t.Errorf("Expected 10, got %v", result)
	}
}
