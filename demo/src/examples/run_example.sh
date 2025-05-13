#!/bin/bash

# Build the RunInk tool
cd ~/runink_demo/src
go build -o runink

# Create a directory for example files
mkdir -p ~/runink_demo/examples/files

# Copy the sample files from ~/Uploads
cp ~/Uploads/fdc3events.contract ~/runink_demo/examples/files/
cp ~/Uploads/fdc3events.conf ~/runink_demo/examples/files/
cp ~/Uploads/fdc3events.dsl ~/runink_demo/examples/files/
cp ~/Uploads/finance.herd ~/runink_demo/examples/files/

# Run the tool with the example files
./runink run \
  --contract ~/runink_demo/examples/files/fdc3events.contract \
  --conf ~/runink_demo/examples/files/fdc3events.conf \
  --dsl ~/runink_demo/examples/files/fdc3events.dsl \
  --herd ~/runink_demo/examples/files/finance.herd \
  --verbose
