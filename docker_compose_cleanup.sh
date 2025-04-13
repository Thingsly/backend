#!/bin/bash

echo "Starting full Docker Compose project cleanup process..."

# Check if docker-compose.yml file exists
if [ ! -f "docker-compose.yml" ]; then
    echo "Error: docker-compose.yml file not found in the current directory."
    exit 1
fi

# Function: Execute a command and continue even if it fails
run_command() {
    echo "Running: $1"
    $1 || echo "Warning: Command '$1' failed, but continuing script execution."
}

# Stop and remove all containers
run_command "docker-compose down --remove-orphans"

# Remove project volumes
run_command "docker-compose down -v"

# Remove project-related images
run_command "docker-compose down --rmi all"

# Remove local host-mounted data directories
echo "Removing local host-mounted data directories..."
data_dirs=$(grep -E '^\s*-\s*\./.+:/.*' docker-compose.yml | awk '{print $2}' | cut -d':' -f1)
for dir in $data_dirs; do
    if [ -d "$dir" ]; then
        echo "Deleting directory: $dir"
        run_command "rm -rf $dir"
    fi
done

# Special handling for ./data directory
if [ -d "./data" ]; then
    echo "Deleting ./data directory"
    run_command "rm -rf ./data"
fi

echo "Checking for any remaining project resources..."
echo "Project containers:"
docker-compose ps
echo "Project volumes:"
docker volume ls --filter name=$(basename $(pwd))
echo "Project networks:"
docker network ls --filter name=$(basename $(pwd))

echo "Cleanup completed."
echo "Warning: This script has deleted local data directories defined in docker-compose.yml. Make sure you no longer need this data."
echo "If any resources remain, you may need to remove them manually."
echo "Note: If errors occurred during cleanup, it may be due to issues in the docker-compose.yml configuration."
echo "Please review and fix any errors in the docker-compose.yml file if necessary."
