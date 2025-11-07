#!/bin/bash

# Usage: ./update_project.sh <new-import-path> <old-variable-name> <new-variable-name>
# The script assumes the target directory is the directory where the script is located.

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
NEW_IMPORT_PATH=$1
OLD_VARIABLE_NAME=$2
NEW_VARIABLE_NAME=$3

# Define old import path and module path
OLD_IMPORT_PATH="github.com/rendyfutsuy/base-go"
OLD_MODULE_PATH="go_base_project"
OLD_APP_NAME="go base project"
OLD_DB_NAME="go_base_project"

# Function to convert string to snake_case
to_snake_case() {
  echo "$1" | awk '{
    gsub(/([A-Z])/, "_\\1")
    gsub(/^[A-Z]/, "\\L&")
    gsub(/[^a-zA-Z0-9_]/, "_")
    print tolower($0)
  }'
}

# Function to capitalize the first letter of each word
capitalize_words() {
  echo "$1" | awk '{
    n = split($0, a, /[^a-zA-Z0-9]/);
    for (i = 1; i <= n; i++) {
      a[i] = toupper(substr(a[i], 1, 1)) tolower(substr(a[i], 2));
    }
    print a[1] (n > 1 ? (a[2] ? a[2] : "") : "");
  }'
}

if [ -z "$NEW_IMPORT_PATH" ] || [ -z "$OLD_VARIABLE_NAME" ] || [ -z "$NEW_VARIABLE_NAME" ]; then
  echo "Usage: $0 <new-import-path> <old-variable-name> <new-variable-name>"
  exit 1
fi

# Extract the last segment from the new import path and convert to snake_case
LAST_SEGMENT=$(echo "$NEW_IMPORT_PATH" | awk -F'/' '{print $NF}')
LAST_SEGMENT_SNAKE=$(to_snake_case "$LAST_SEGMENT")

# Capitalize the first letter of each word in NEW_VARIABLE_NAME
NEW_VARIABLE_NAME_CAPITALIZED=$(capitalize_words "$NEW_VARIABLE_NAME")

echo "Script Directory: $SCRIPT_DIR"
echo "New Import Path: $NEW_IMPORT_PATH"
echo "Old Variable Name: $OLD_VARIABLE_NAME"
echo "New Variable Name: $NEW_VARIABLE_NAME"
echo "New Variable Name Capitalized: $NEW_VARIABLE_NAME_CAPITALIZED"
echo "Last Segment: $LAST_SEGMENT"
echo "Last Segment in Snake Case: $LAST_SEGMENT_SNAKE"

# Replace the import path in all Go files in the script's directory
find "$SCRIPT_DIR" -type f -name "*.go" -print0 | xargs -0 sed -i '' "s|$OLD_IMPORT_PATH|$NEW_IMPORT_PATH|g"

# Replace the variable name in all Go files in the script's directory
find "$SCRIPT_DIR" -type f -name "*.go" -print0 | xargs -0 sed -i '' "s|$OLD_VARIABLE_NAME|$NEW_VARIABLE_NAME|g"

# Replace DatabaseGoBaseProject and setStringConnectionGoBaseProject patterns in all Go files
find "$SCRIPT_DIR" -type f -name "*.go" -print0 | xargs -0 sed -i '' "s|DatabaseGoBaseProject|$NEW_VARIABLE_NAME|g"
find "$SCRIPT_DIR" -type f -name "*.go" -print0 | xargs -0 sed -i '' "s|setStringConnectionGoBaseProject|setStringConnection$NEW_VARIABLE_NAME_CAPITALIZED|g"

# Replace the module path in the go.mod file
if [ -f "$SCRIPT_DIR/go.mod" ]; then
  sed -i '' "s|module $OLD_IMPORT_PATH|module $NEW_IMPORT_PATH|g" "$SCRIPT_DIR/go.mod"
else
  echo "No go.mod file found in $SCRIPT_DIR"
fi

# Replace go_base_project with the last segment in snake_case in all Go files
find "$SCRIPT_DIR" -type f -name "*.go" -print0 | xargs -0 sed -i '' "s|$OLD_MODULE_PATH|$LAST_SEGMENT_SNAKE|g"

# Update the app_name and db_name in the config file
CONFIG_FILE="$SCRIPT_DIR/config.json.example"

if [ -f "$CONFIG_FILE" ]; then
  # Replace app_name
  sed -i '' "s|\"app_name\": \"$OLD_APP_NAME\"|\"app_name\": \"$LAST_SEGMENT\"|g" "$CONFIG_FILE"

  # Replace db_name in all occurrences
  sed -i '' "s|\"$OLD_DB_NAME\"|\"$LAST_SEGMENT_SNAKE\"|g" "$CONFIG_FILE"

  echo "Config file updated:"
  echo " - app_name '$OLD_APP_NAME' replaced with '$LAST_SEGMENT'"
  echo " - Database names '$OLD_DB_NAME' replaced with '$LAST_SEGMENT_SNAKE'"
else
  echo "No config file found at $CONFIG_FILE"
fi

echo "Replacements completed in $SCRIPT_DIR:"
echo " - Import path '$OLD_IMPORT_PATH' replaced with '$NEW_IMPORT_PATH'"
echo " - Variable name '$OLD_VARIABLE_NAME' replaced with '$NEW_VARIABLE_NAME'"
echo " - Module path '$OLD_IMPORT_PATH' replaced with '$NEW_IMPORT_PATH'"
echo " - '$OLD_MODULE_PATH' replaced with '$LAST_SEGMENT_SNAKE'"
