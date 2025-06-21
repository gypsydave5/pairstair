#!/bin/sh
set -e

# Move to the script's directory
cd "$(dirname "$0")"

# Check if pandoc is installed
if ! command -v pandoc >/dev/null 2>&1; then
  echo "Error: pandoc is not installed. Please install pandoc first."
  echo "See https://pandoc.org/installing.html for installation instructions."
  exit 1
fi

# Convert markdown to man page
pandoc -s -t man pairstair.1.md -o pairstair.1
echo "Generated man page: pairstair.1"

# Make the script executable
chmod +x "$(basename "$0")"

echo ""
echo "To view the man page, run:"
echo "man -l pairstair.1"
echo ""
echo "To install the man page system-wide (requires sudo):"
echo "sudo cp pairstair.1 /usr/local/share/man/man1/"
echo "sudo mandb"
