package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	
	#TODO: Create the following packages
	"github.com/runink/pkg/contract"
	"github.com/runink/pkg/dsl"
	"github.com/runink/pkg/envloader"
	"github.com/runink/pkg/generator"
	"github.com/runink/pkg/render"
)

func main() {
	// Parse CLI flags
	featurePath := flag.String("feature", "", "Path to the feature.dsl file")
	contractPath := flag.String("contract", "", "Path to the Go contract file")
	envPath := flag.String("env", "", "Path to the feature.env file")
	outputDir := flag.String("out", "rendered", "Directory to write rendered output")
	flag.Parse()

	// Mandatory parameters check
	if *featurePath == "" || *contractPath == "" || *envPath == "" {
		log.Fatal("❌ Missing required parameters: --feature, --contract, --env")
	}

	// Load environment variables for interpolation
	if err := envloader.LoadFromFile(*envPath); err != nil {
		log.Fatalf("❌ Failed to load env file: %v", err)
	}

	// Parse DSL Scenario
	scenario, err := dsl.ParseFile(*featurePath)
	if err != nil {
		log.Fatalf("❌ Failed to parse DSL: %v", err)
	}

	// Parse Contract Steps
	steps, err := contract.ExtractAnnotatedSteps(*contractPath)
	if err != nil {
		log.Fatalf("❌ Failed to parse contract: %v", err)
	}

	// Preallocate DAG with rough estimate based on step count
	dag := generator.NewDAG(len(steps) + 2)

	// Build DAG from scenario and contract
	if err := generator.BuildDAGInto(dag, scenario, steps); err != nil {
		log.Fatalf("❌ DAG build failed: %v", err)
	}

	// Apply DAG optimizations:
	// - Enable zero-copy pipelines
	// - Enforce struct reuse
	// - Honor affinity and placement constraints
	dag.Optimize(generator.Options{
		EnableStreaming:    true,
		UseStructReuse:     true,
		ApplyAffinityTags:  true,
		ApplyCheckpointing: true,
	})

	// Create output directory if missing
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create output dir: %v", err)
	}

	// Setup streaming output writer
	outPath := filepath.Join(*outputDir, sanitize(scenario.Metadata.Feature)+".go")
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("❌ Failed to open output file: %v", err)
	}
	defer outFile.Close()

	// Pipe DAG rendering through an io.Pipe to avoid buffering large DAGs in memory
	reader, writer := io.Pipe()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(outFile, reader); err != nil {
			log.Fatalf("❌ Failed streaming DAG to file: %v", err)
		}
	}()

	if err := render.StreamToGoPipe(dag, writer); err != nil {
		log.Fatalf("❌ Failed to render DAG: %v", err)
	}
	writer.Close()
	wg.Wait()

	log.Printf("✅ DAG generated successfully: %s", outPath)
}

// sanitize ensures safe filenames
func sanitize(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
