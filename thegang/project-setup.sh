#!/bin/bash

# Get the original module name
OLD_MODULE_NAME=$(go list -m)

# Get the remote URL of the repository
remote_url=$(git remote get-url origin)

# Extract the GitHub repository path from the remote URL
github_repo_path=$(echo "$remote_url" | sed -E 's/.*github.com[:/](.*).git/\1/')

NEW_MODULE_NAME="github.com/"$github_repo_path

# Change the module name in go.mod
go mod edit -module "$NEW_MODULE_NAME";

# Escape /
OLD_MODULE_NAME=$(echo "$OLD_MODULE_NAME" | sed 's/\//\\\//g')
NEW_MODULE_NAME=$(echo "$NEW_MODULE_NAME" | sed 's/\//\\\//g')

# Change the module name in the source code
find . -type f -name '*.go' \
  -exec sed -i -e "s/$OLD_MODULE_NAME/$NEW_MODULE_NAME/g" {} \;

# Find all files with a '.go-e' extension
find . -type f -name "*.go-e" | while read file; do
    # Get the original file name by removing the '.go-e' extension
    original_file=$(echo $file | sed 's/\.go-e$//')
    echo "$original_file"
    # If the original file exists
    if [ -f "$original_file.go" ]; then
        # Remove the backup file
        rm "$file"
    fi
done

# Copy docker-compose to root
cp thegang/docker-compose.yaml .