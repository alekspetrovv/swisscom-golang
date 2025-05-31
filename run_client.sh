#!/bin/bash

DEFAULT_PARALLEL_REQUESTS=100
DEFAULT_STEPS=10
DEFAULT_TARGET_URL="http://localhost:4005/api/services"

PARALLEL_REQUESTS=${1:-$DEFAULT_PARALLEL_REQUESTS}
STEPS=${2:-$DEFAULT_STEPS}
TARGET_URL=${3:-$DEFAULT_TARGET_URL}
GO_CLIENT_FILE=load-test.go

echo "Launching Go HTTP Batch Client..."
echo "------------------------------------------"
echo "Effective Configuration:"
echo "  Parallel Requests per Step: $PARALLEL_REQUESTS"
echo "  Number of Steps:            $STEPS"
echo "  Target URL:                 $TARGET_URL"
echo "------------------------------------------"

go run "$GO_CLIENT_FILE" \
    -parallel="$PARALLEL_REQUESTS" \
    -steps="$STEPS" \
    -url="$TARGET_URL"

echo "------------------------------------------"
echo "Script finished."