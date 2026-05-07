#!/bin/sh

# Exit immediately if a command fails
set -e

echo "Running migrations..."
./migrate up

echo "Starting API..."
exec ./main