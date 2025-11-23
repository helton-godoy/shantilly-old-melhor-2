#!/bin/bash
set -e

REPO="helton-godoy/shantilly"

# Function to protect a branch
protect_branch() {
    BRANCH=$1
    echo "🛡️ Protecting branch: $BRANCH..."

    gh api -X PUT "repos/$REPO/branches/$BRANCH/protection" \
        -H "Accept: application/vnd.github+json" \
        --input - <<< '{
            "required_status_checks": {
                "strict": true,
                "contexts": ["Build and Test", "Lint"]
            },
            "enforce_admins": true,
            "required_pull_request_reviews": {
                "dismiss_stale_reviews": true,
                "require_code_owner_reviews": false,
                "required_approving_review_count": 0
            },
            "restrictions": null,
            "allow_force_pushes": false,
            "allow_deletions": false
        }' > /dev/null

    echo "✅ Branch $BRANCH protected successfully!"
}

# Main
echo "🚀 Configuring Branch Protection for $REPO"

# Check auth
if ! gh auth status > /dev/null 2>&1; then
    echo "❌ Error: Not authenticated to GitHub"
    exit 1
fi

protect_branch "main"
protect_branch "develop"

echo "🎉 All branches protected!"
