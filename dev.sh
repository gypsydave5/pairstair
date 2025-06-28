#!/bin/sh
set -e

# PairStair Development Helper Script
# Common development tasks and checks

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Development helper commands:"
    echo ""
    echo "Commands:"
    echo "  docs        Regenerate man page documentation"
    echo "  version     Show current version info"
    echo ""
    echo "Examples:"
    echo "  $0 docs     # Regenerate man page from markdown"
    echo "  $0 version  # Show version info and git tags"
    echo ""
    echo "Note: Use 'go test ./...', 'go build', etc. directly for basic operations"
    echo ""
    exit 1
}

# Check arguments
if [ $# -ne 1 ]; then
    usage
fi

COMMAND="$1"

case "$COMMAND" in
    "docs")
        echo "📖 Regenerating documentation..."
        cd docs
        ./gen-man.sh
        cd ..
        echo "✓ Documentation updated"
        ;;
    
    "version")
        echo "📋 Version information:"
        echo ""
        if [ -f pairstair ]; then
            echo "Built binary version:"
            ./pairstair --version 2>/dev/null || echo "(version check failed)"
        else
            echo "No binary found - building temporarily..."
            go build
            echo "Version:"
            ./pairstair --version 2>/dev/null || echo "(version check failed)"
            rm -f pairstair
        fi
        echo ""
        echo "Latest git tag:"
        git tag --sort=-version:refname | head -1 2>/dev/null || echo "(no tags found)"
        echo ""
        echo "Current commit:"
        git rev-parse --short HEAD 2>/dev/null || echo "(not in git repo)"
        ;;
    
    *)
        echo "Error: Unknown command '$COMMAND'"
        usage
        ;;
esac
