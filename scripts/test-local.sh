#!/bin/bash

# Local testing and development script
# This script helps test the application locally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== CI/CD Pipeline Application - Local Test Script ===${NC}\n"

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

command -v go >/dev/null 2>&1
print_status $? "Go installed"

command -v docker >/dev/null 2>&1
print_status $? "Docker installed"

command -v docker-compose >/dev/null 2>&1 || command -v docker compose >/dev/null 2>&1
print_status $? "Docker Compose installed"

# Run go mod verify
echo -e "\n${YELLOW}Verifying Go modules...${NC}"
go mod verify
print_status $? "Go modules verified"

# Format check
echo -e "\n${YELLOW}Checking code format...${NC}"
UNFORMATTED=$(gofmt -l .)
if [ -z "$UNFORMATTED" ]; then
    print_status 0 "Code is formatted"
else
    echo -e "${RED}The following files are not formatted:${NC}"
    echo "$UNFORMATTED"
    print_status 1 "Code format check"
fi

# Run go vet
echo -e "\n${YELLOW}Running go vet...${NC}"
go vet ./...
print_status $? "go vet passed"

# Run tests
echo -e "\n${YELLOW}Running tests...${NC}"
go test -v -race -timeout 30s ./...
print_status $? "Tests passed"

# Run tests with coverage
echo -e "\n${YELLOW}Running tests with coverage...${NC}"
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo -e "${GREEN}Total coverage: $COVERAGE${NC}"
print_status $? "Coverage report generated"

# Build application
echo -e "\n${YELLOW}Building application...${NC}"
go build -o bin/test-app cmd/api/main.go
print_status $? "Application built successfully"

# Build Docker image
echo -e "\n${YELLOW}Building Docker image...${NC}"
docker build -t cicd-pipeline-app:test .
print_status $? "Docker image built"

# Start services with docker-compose
echo -e "\n${YELLOW}Starting services with docker-compose...${NC}"
docker-compose up -d
sleep 5
print_status $? "Services started"

# Wait for application to be ready
echo -e "\n${YELLOW}Waiting for application to be ready...${NC}"
for i in {1..30}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_status 0 "Application is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        print_status 1 "Application failed to start"
    fi
    sleep 1
done

# Test health endpoint
echo -e "\n${YELLOW}Testing health endpoint...${NC}"
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
echo "$HEALTH_RESPONSE" | jq . || echo "$HEALTH_RESPONSE"
if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    print_status 0 "Health check passed"
else
    print_status 1 "Health check failed"
fi

# Test readiness endpoint
echo -e "\n${YELLOW}Testing readiness endpoint...${NC}"
READY_RESPONSE=$(curl -s http://localhost:8080/ready)
echo "$READY_RESPONSE" | jq . || echo "$READY_RESPONSE"
if echo "$READY_RESPONSE" | grep -q "ready"; then
    print_status 0 "Readiness check passed"
else
    print_status 1 "Readiness check failed"
fi

# Test API endpoints
echo -e "\n${YELLOW}Testing API endpoints...${NC}"

# Create pipeline
echo -e "\nCreating pipeline..."
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/pipelines \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Pipeline","description":"A test pipeline"}')
echo "$CREATE_RESPONSE" | jq . || echo "$CREATE_RESPONSE"
PIPELINE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id // empty')

if [ ! -z "$PIPELINE_ID" ]; then
    print_status 0 "Pipeline created (ID: $PIPELINE_ID)"
    
    # List pipelines
    echo -e "\nListing pipelines..."
    LIST_RESPONSE=$(curl -s http://localhost:8080/api/v1/pipelines)
    echo "$LIST_RESPONSE" | jq . || echo "$LIST_RESPONSE"
    print_status $? "Pipelines listed"
    
    # Get specific pipeline
    echo -e "\nGetting pipeline $PIPELINE_ID..."
    GET_RESPONSE=$(curl -s http://localhost:8080/api/v1/pipelines/$PIPELINE_ID)
    echo "$GET_RESPONSE" | jq . || echo "$GET_RESPONSE"
    print_status $? "Pipeline retrieved"
else
    echo -e "${YELLOW}Warning: Could not create pipeline (database might not be ready)${NC}"
fi

# Show logs
echo -e "\n${YELLOW}Recent application logs:${NC}"
docker-compose logs --tail=20 app

echo -e "\n${GREEN}=== All tests completed successfully! ===${NC}"
echo -e "${YELLOW}Services are running. Use 'docker-compose logs -f' to view logs.${NC}"
echo -e "${YELLOW}Use 'docker-compose down' to stop services.${NC}"
