
// CSV to PARQUET Conversion DAG
// This DAG demonstrates loading a CSV file and converting it to PARQUET format

DAG csv_to_parquet {
  // Define the nodes in our DAG
  node load_csv {
    type = "source"
    source_type = "csv"
    path = "./source/sample.csv"
    contract = "csv_contract.yaml"
    options {
      header = true
      inferSchema = false
    }
  }

  node save_parquet {
    type = "sink"
    sink_type = "parquet"
    path = "./target/sample.parquet"
    options {
      compression = "snappy"
      overwrite = true
    }
  }

  // Define the edges (data flow)
  edge {
    from = "load_csv"
    to = "save_parquet"
  }
}
