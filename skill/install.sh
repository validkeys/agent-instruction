#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_SOURCE="${SCRIPT_DIR}/agent-instruction"

# Default Claude Code skills directory
CLAUDE_SKILLS_DIR="${CLAUDE_SKILLS_DIR:-$HOME/.local/share/claude/skills}"
SKILL_DEST="${CLAUDE_SKILLS_DIR}/agent-instruction"

# Dry run mode
DRY_RUN=false
if [[ "$1" == "--dry-run" ]]; then
    DRY_RUN=true
fi

# Print functions
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo "ℹ $1"
}

# Check if source skill directory exists
if [[ ! -d "${SKILL_SOURCE}" ]]; then
    print_error "Skill source directory not found: ${SKILL_SOURCE}"
    exit 1
fi

# Check if skill.yaml exists
if [[ ! -f "${SKILL_SOURCE}/skill.yaml" ]]; then
    print_error "skill.yaml not found in ${SKILL_SOURCE}"
    exit 1
fi

print_info "Installing agent-instruction skill for Claude Code"
echo

# Create skills directory if it doesn't exist
if [[ ! -d "${CLAUDE_SKILLS_DIR}" ]]; then
    print_warning "Skills directory does not exist: ${CLAUDE_SKILLS_DIR}"

    if [[ "$DRY_RUN" == true ]]; then
        print_info "Would create: ${CLAUDE_SKILLS_DIR}"
    else
        print_info "Creating skills directory..."
        mkdir -p "${CLAUDE_SKILLS_DIR}" || {
            print_error "Failed to create skills directory"
            print_info "Try: mkdir -p ${CLAUDE_SKILLS_DIR}"
            exit 1
        }
        print_success "Created ${CLAUDE_SKILLS_DIR}"
    fi
fi

# Check if skill is already installed
if [[ -d "${SKILL_DEST}" ]]; then
    print_warning "Skill already installed at: ${SKILL_DEST}"

    if [[ "$DRY_RUN" == true ]]; then
        print_info "Would overwrite existing installation"
    else
        read -p "Overwrite existing installation? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled"
            exit 0
        fi

        print_info "Removing existing installation..."
        rm -rf "${SKILL_DEST}" || {
            print_error "Failed to remove existing installation"
            exit 1
        }
    fi
fi

# Copy skill directory
if [[ "$DRY_RUN" == true ]]; then
    print_info "Would copy: ${SKILL_SOURCE} -> ${SKILL_DEST}"
else
    print_info "Installing skill..."
    cp -r "${SKILL_SOURCE}" "${SKILL_DEST}" || {
        print_error "Failed to copy skill directory"
        exit 1
    }
    print_success "Copied skill to ${SKILL_DEST}"
fi

# Validate installation
if [[ "$DRY_RUN" == false ]]; then
    if [[ -f "${SKILL_DEST}/skill.yaml" ]]; then
        print_success "Validated skill.yaml exists"
    else
        print_error "Validation failed: skill.yaml not found after installation"
        exit 1
    fi

    # Check if agent-instruction CLI is available
    if command -v agent-instruction &> /dev/null; then
        AGENT_VERSION=$(agent-instruction --version 2>&1 || echo "unknown")
        print_success "agent-instruction CLI found: ${AGENT_VERSION}"
    else
        print_warning "agent-instruction CLI not found in PATH"
        print_info "Install with: go install github.com/yourusername/agent-instruction@latest"
    fi
fi

# Print success message
echo
if [[ "$DRY_RUN" == true ]]; then
    print_success "Dry run completed successfully"
else
    print_success "Installation completed successfully"
    echo
    print_info "Next steps:"
    echo "  1. Verify installation: ls ${SKILL_DEST}"
    echo "  2. Check Claude Code recognizes it: claude skills list"
    echo "  3. Test the skill in a Claude Code session"

    if ! command -v agent-instruction &> /dev/null; then
        echo
        print_warning "Don't forget to install the agent-instruction CLI"
    fi
fi
