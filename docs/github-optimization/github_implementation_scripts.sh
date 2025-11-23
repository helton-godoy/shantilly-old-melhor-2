#!/bin/bash

# GitHub Structures Implementation Script for Shantilly
# Repository: helton-godoy/shantilly
# This script creates all necessary GitHub structures

set -e

echo "🚀 Shantilly GitHub Structures Implementation"
echo "============================================="

# Configuration
REPO="helton-godoy/shantilly"
REPO_OWNER="helton-godoy"
REPO_NAME="shantilly"

# Function to create labels
create_labels() {
    echo "📋 Creating GitHub Labels..."
    
    # Priority Labels
    gh label create "priority::critical" --color "d73a4a" --description "Critical priority - must fix immediately" --repo "$REPO" --force
    gh label create "priority::high" --color "fb8500" --description "High priority - important for next release" --repo "$REPO" --force
    gh label create "priority::medium" --color "fbbf24" --description "Medium priority - nice to have" --repo "$REPO" --force
    gh label create "priority::low" --color "10b981" --description "Low priority - future consideration" --repo "$REPO" --force
    
    # Type Labels
    gh label create "type::feature" --color "1f77b4" --description "New functionality" --repo "$REPO" --force
    gh label create "type::bug" --color "d62728" --description "Something isn't working" --repo "$REPO" --force
    gh label create "type::refactor" --color "ff7f0e" --description "Refactoring code" --repo "$REPO" --force
    gh label create "type::task" --color "2ca02c" --description "Non-code related tasks" --repo "$REPO" --force
    gh label create "type::security" --color "e377c2" --description "Security related" --repo "$REPO" --force
    gh label create "type::documentation" --color "9467bd" --description "Documentation changes" --repo "$REPO" --force
    gh label create "type::optimization" --color "bcbd22" --description "Performance improvements" --repo "$REPO" --force
    
    # Area Labels
    gh label create "area::core" --color "8c564b" --description "Core functionality" --repo "$REPO" --force
    gh label create "area::ui" --color "17becf" --description "User Interface" --repo "$REPO" --force
    gh label create "area::security" --color "9467bd" --description "Security" --repo "$REPO" --force
    gh label create "area::runtime" --color "bcbd22" --description "Runtime engine" --repo "$REPO" --force
    gh label create "area::integrations" --color "ffbb78" --description "Third-party integrations" --repo "$REPO" --force
    gh label create "area::infrastructure" --color "98df8a" --description "Infrastructure" --repo "$REPO" --force
    
    # Status Labels
    gh label create "status::triage" --color "c7c7c7" --description "Needs initial assessment" --repo "$REPO" --force
    gh label create "status::in-progress" --color "6f42c1" --description "Currently being worked on" --repo "$REPO" --force
    gh label create "status::review" --color "fd7e14" --description "Needs code review" --repo "$REPO" --force
    gh label create "status::blocked" --color "dc3545" --description "Cannot proceed" --repo "$REPO" --force
    gh label create "status::done" --color "28a745" --description "Completed" --repo "$REPO" --force
    
    # Epic Labels
    gh label create "epic::e1" --color "6f42c1" --description "E1 Foundation" --repo "$REPO" --force
    gh label create "epic::e2" --color "fd7e14" --description "E2 Advanced Features" --repo "$REPO" --force
    gh label create "epic::e3" --color "20c997" --description "E3 Multi-Panel" --repo "$REPO" --force
    gh label create "epic::e4" --color "ffc107" --description "E4 ScriptRunner" --repo "$REPO" --force
    gh label create "epic::e5" --color "17a2b8" --description "E5 Security" --repo "$REPO" --force
    
    echo "✅ Labels created successfully!"
}

# Function to create milestones
create_milestones() {
    echo "🎯 Creating Milestones..."
    
    # v1.0 Alpha
    gh api repos/$REPO/milestones -f title="v1.0 Alpha" \
        -f description="MVP funcional com foundations básicas - CLI Foundation, YAML Parsing, TUI Structure, Form Rendering" \
        -f due_on="2025-12-19T23:59:59Z" || echo "⚠️ Milestone v1.0 Alpha might already exist"

    # v1.0 Beta
    gh api repos/$REPO/milestones -f title="v1.0 Beta" \
        -f description="Funcionalidades avançadas e UX melhorada - Advanced Form Types, Enhanced Error Handling" \
        -f due_on="2026-01-15T23:59:59Z" || echo "⚠️ Milestone v1.0 Beta might already exist"

    # v2.0 Features
    gh api repos/$REPO/milestones -f title="v2.0 Features" \
        -f description="Runtime architecture completa - Event Engine, ScriptRunner, Modal Stack, Security Hardening" \
        -f due_on="2026-02-20T23:59:59Z" || echo "⚠️ Milestone v2.0 Features might already exist"

    # v2.0 Extensions
    gh api repos/$REPO/milestones -f title="v2.0 Extensions" \
        -f description="Integrações e extensões avançadas - Ansible Integration, SSH Server, Plugin Architecture" \
        -f due_on="2026-03-15T23:59:59Z" || echo "⚠️ Milestone v2.0 Extensions might already exist"
    
    echo "✅ Milestones created successfully!"
}

# Function to create project
create_project() {
    echo "📊 Creating GitHub Project 'Shantilly Roadmap'..."
    
    gh project create \
        --owner "$REPO_OWNER" \
        --title "Shantilly Roadmap"
    
    echo "✅ Project created successfully!"
    echo "ℹ️  Note: Project columns and automations need to be configured manually via web interface"
}

# Function to import issues
import_issues() {
    echo "📋 Importing structured issues..."
    
    # This would need to be implemented based on the issues structure
    # For now, we'll create a placeholder
    echo "ℹ️  Issue import requires manual implementation using github_roadmap_issues.md"
}

# Main execution
main() {
    echo "Starting GitHub structures implementation for $REPO"
    echo ""
    
    # Check if authenticated
    if ! gh auth status > /dev/null 2>&1; then
        echo "❌ Error: Not authenticated to GitHub"
        echo "Please run: gh auth login"
        exit 1
    fi
    
    # Execute functions
    create_labels
    echo ""
    create_milestones
    echo ""
    create_project
    echo ""
    import_issues
    
    echo ""
    echo "🎉 GitHub structures implementation completed!"
    echo ""
    echo "Next steps:"
    echo "1. Configure project columns and automations manually"
    echo "2. Import issues using github_roadmap_issues.md"
    echo "3. Test automation rules"
    echo "4. Validate all structures"
}

# Run main function
main
