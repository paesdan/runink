
# CSV to PARQUET Conversion Demo

This example demonstrates how to use the RunInk DAG execution engine to convert a CSV file to PARQUET format. The demo showcases several advanced features of the RunInk engine including data contract validation, resource allocation, and configurable execution settings.

## Directory Structure

```
csv_to_parquet/
├── conf.yaml           # Execution configuration
├── csv_contract.yaml   # Data schema definition
├── csv_to_parquet.dsl  # DAG definition
├── herd.yaml           # Resource allocation
├── source/             # Source directory
│   └── sample.csv      # Sample CSV data
└── target/             # Target directory for output
    └── sample.parquet  # (will be created when DAG runs)
```

## Files Description

- **csv_to_parquet.dsl**: Defines the DAG pipeline with nodes for loading CSV and saving as PARQUET
- **csv_contract.yaml**: Defines the schema for the CSV data
- **conf.yaml**: Contains configuration settings for execution, logging, monitoring, and error handling
- **herd.yaml**: Specifies resource allocation for the DAG execution
- **source/sample.csv**: Sample CSV file with user data
- **target/**: Directory where the PARQUET file will be saved

## Running the Demo

To run the CSV to PARQUET conversion demo, use the Cobra CLI command:

```bash
# Navigate to the example directory
cd ~/runink_demo/examples/csv_to_parquet

# Run the DAG using the RunInk Cobra CLI
runink dag run --dsl csv_to_parquet.dsl --conf conf.yaml --herd herd.yaml
```

## Expected Output

After successful execution, you should see:

1. Console output showing the progress and completion of the DAG execution
2. A new PARQUET file created at `./target/sample.parquet`
3. Metrics data saved at `./metrics/csv_to_parquet_metrics.json` (if monitoring is enabled)

You can verify the PARQUET file was created correctly using:

```bash
# List the target directory
ls -la ./target

# If you have parquet-tools installed, you can view the content
# parquet-tools show ./target/sample.parquet
```

## Advanced Usage

### Changing Execution Mode

You can modify the `conf.yaml` file to change the execution mode from sequential to parallel:

```yaml
execution:
  mode: "parallel"
```

### Adjusting Resource Allocation

Modify the `herd.yaml` file to adjust CPU, memory, or disk resources based on your requirements.

### Custom Error Handling

Change the error handling strategy in `conf.yaml`:

```yaml
error_handling:
  strategy: "retry"  # Options: "stop_on_error", "continue", or "retry"
```

## Troubleshooting

If you encounter any issues:

1. Check the logs for detailed error messages
2. Verify that all required files are in the correct locations
3. Ensure the RunInk engine is properly installed and configured
4. Confirm that the CSV file matches the schema defined in the contract

For more information, refer to the RunInk documentation.
