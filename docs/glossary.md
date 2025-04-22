# Runink Glossary

This glossary defines key terms, acronyms, and concepts used throughout the Runink documentation and codebase.

---

### .feature File

A human-readable file written in Gherkin syntax used to describe a data pipeline scenario using Given/When/Then structure and tags like @source, @step, and @sink.

---

### BDD (Behavior-Driven Development)

A software development approach that describes application behavior in plain language, often used with .feature files.

---

### Golden File

A snapshot of the expected output from a pipeline or transformation, used to assert correctness in automated tests.

---

### Schema Contract

A versioned definition of a data structure (e.g., JSON, Protobuf, Go struct) used to validate pipeline input/output and detect schema drift.

---

### Schema Drift

An unintended or unexpected change in a schema that may break pipeline compatibility.

---

### DAG (Directed Acyclic Graph)

A graph of pipeline stages where each edge represents a dependency, and there are no cycles. Used for orchestrating non-linear workflows.

---

### DLQ (Dead Letter Queue)

A place to store invalid or failed records so they can be analyzed and retried later without interrupting the rest of the pipeline.

---

### REPL (Read-Eval-Print Loop)

An interactive interface that lets users type commands and immediately see the results. Runink’s REPL supports loading data, applying steps, and viewing outputs.

---

### Lineage

A traceable path showing how data flows through each pipeline stage, from source to sink, including what transformed it and which contract was applied.

---

### Prometheus

A monitoring system used to collect and store metrics from pipeline executions.

---

### OpenTelemetry

An observability framework for collecting traces and metrics, helping to visualize the execution path and performance of pipelines.

---

### gRPC

A high-performance, open-source universal RPC framework used for running distributed pipeline stages.

---

### Protobuf (Protocol Buffers)

A method for serializing structured data, used in gRPC communication and schema definitions.

---

### Bigmachine

A Go library for orchestrating distributed, stateless workers. Used by Runink to scale pipelines across multiple machines.

---

### Contract Hash

A hash value generated from a schema contract to uniquely identify its version. Used for detecting changes and tracking usage.

---

### SBOM (Software Bill of Materials)

A manifest of all dependencies and components included in a software release, used for compliance and security auditing.

---

### FDC3 (Financial Desktop Connectivity and Collaboration Consortium)

A standard for interop between financial applications. Runink can integrate with FDC3 schemas and messaging models.

---

### CDM (Common Domain Model)

A standardized schema used in finance and trading to represent products, trades, and events. Supported natively by Runink.

---

Have a term you’d like added? Open an issue or suggest a change in the docs!
