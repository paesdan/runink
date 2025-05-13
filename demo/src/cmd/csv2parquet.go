package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/runink/runink/nodes"
	"github.com/spf13/cobra"
)

// csv2parquetCmd represents the csv2parquet command
var csv2parquetCmd = &cobra.Command{
	Use:   "csv2parquet",
	Short: "Convert CSV files to PARQUET format",
	Long: `Convert CSV files to PARQUET format using the RunInk DAG execution engine.
This command demonstrates the CSV reader and Parquet writer nodes.

Example usage:
  runink csv2parquet --csv examples/weather.csv --parquet output/weather.parquet
  runink csv2parquet --csv examples/weather.csv --use-dag`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		csvPath, err := cmd.Flags().GetString("csv")
		if err != nil {
			return err
		}

		parquetPath, err := cmd.Flags().GetString("parquet")
		if err != nil {
			return err
		}

		useDag, err := cmd.Flags().GetBool("use-dag")
		if err != nil {
			return err
		}

		// Validate flags
		if csvPath == "" {
			return fmt.Errorf("csv path is required")
		}

		// If parquet path is not provided, use the same path as the CSV file but with .parquet extension
		if parquetPath == "" {
			dir, file := filepath.Split(csvPath)
			ext := filepath.Ext(file)
			name := file[:len(file)-len(ext)]
			parquetPath = filepath.Join(dir, name+".parquet")
		}

		// Run the example
		fmt.Printf("Converting CSV file %s to PARQUET file %s\n", csvPath, parquetPath)
		
		if useDag {
			return nodes.RunCSVToParquetWithDAG(csvPath, parquetPath)
		}
		
		return nodes.RunCSVToParquetExample(csvPath, parquetPath)
	},
}

func init() {
	rootCmd.AddCommand(csv2parquetCmd)

	// Add flags
	csv2parquetCmd.Flags().String("csv", "", "Path to the CSV file (required)")
	csv2parquetCmd.Flags().String("parquet", "", "Path to the output PARQUET file (defaults to CSV path with .parquet extension)")
	csv2parquetCmd.Flags().Bool("use-dag", false, "Use the DAG execution engine (default: false)")
}
