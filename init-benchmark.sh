#!/bin/bash

for port in $(seq 5001 "$1"); do
  go run src/main.go "$port" > "logs/node_${port}.log" 2>&1 &
done

go run src/main.go 5000
