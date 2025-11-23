import re
import subprocess
import sys

def parse_issues(file_path):
    with open(file_path, 'r') as f:
        content = f.read()

    issues = []
    current_milestone = None
    
    # Regex patterns
    milestone_pattern = re.compile(r'### (.*?) \(.*?\)')
    issue_header_pattern = re.compile(r'\*\*Issue #(\d+): (.*?)\*\*')
    meta_pattern = re.compile(r'- Type: (.*?) \| Priority: (.*?) \| Area: (.*?) \| Story Points: (\d+) \| Epic: (.*)')
    desc_pattern = re.compile(r'- Description: (.*)')
    criteria_pattern = re.compile(r'- Acceptance Criteria: (.*)')

    lines = content.split('\n')
    current_issue = {}
    
    for line in lines:
        line = line.strip()
        
        # Match Milestone
        m_match = milestone_pattern.match(line)
        if m_match:
            current_milestone = m_match.group(1).strip()
            continue

        # Match Issue Header
        h_match = issue_header_pattern.match(line)
        if h_match:
            if current_issue:
                issues.append(current_issue)
            current_issue = {
                'number': h_match.group(1),
                'title': h_match.group(2),
                'milestone': current_milestone
            }
            continue

        # Match Metadata
        meta_match = meta_pattern.match(line)
        if meta_match and current_issue:
            current_issue['type'] = meta_match.group(1).strip()
            current_issue['priority'] = meta_match.group(2).strip()
            current_issue['area'] = meta_match.group(3).strip()
            current_issue['points'] = meta_match.group(4).strip()
            current_issue['epic'] = meta_match.group(5).strip()
            continue

        # Match Description
        d_match = desc_pattern.match(line)
        if d_match and current_issue:
            current_issue['description'] = d_match.group(1).strip()
            continue

        # Match Acceptance Criteria
        c_match = criteria_pattern.match(line)
        if c_match and current_issue:
            current_issue['criteria'] = c_match.group(1).strip()
            continue

    # Milestone Mapping
    milestone_map = {
        "Foundation Sprint": "v1.0 Alpha",
        "Advanced Features Sprint": "v1.0 Beta",
        "Runtime Architecture Sprint": "v2.0 Features",
        "Extensions & Integrations": "v2.0 Extensions",
        "Additional Enhancement Issues": "v2.0 Extensions" # Mapping extras to extensions for now
    }

    if current_issue:
        issues.append(current_issue)

    # Apply mapping
    for issue in issues:
        if issue['milestone'] in milestone_map:
            issue['milestone'] = milestone_map[issue['milestone']]
            
    return issues

def create_github_issue(issue):
    # Construct labels
    labels = [
        f"type::{issue['type']}",
        f"priority::{issue['priority']}",
        f"area::{issue['area']}",
        f"epic::{issue['epic'].split(' ')[0].lower()}" # Extract epic code e.g. e1
    ]
    
    # Construct body
    body = f"""## Description
{issue['description']}

## Acceptance Criteria
{issue['criteria'].replace(' | ', '\n- ')}

## Metadata
- **Story Points**: {issue['points']}
- **Epic**: {issue['epic']}
"""

    cmd = [
        'gh', 'issue', 'create',
        '--title', issue['title'],
        '--body', body,
        '--milestone', issue['milestone'],
        '--label', ','.join(labels)
    ]

    print(f"Creating issue: {issue['title']} (Milestone: {issue['milestone']})")
    try:
        subprocess.run(cmd, check=True, capture_output=True, text=True)
        print("✅ Created successfully")
    except subprocess.CalledProcessError as e:
        print(f"❌ Failed to create issue: {e.stderr}")

def main():
    issues_file = 'github_roadmap_issues.md'
    print(f"Parsing issues from {issues_file}...")
    issues = parse_issues(issues_file)
    print(f"Found {len(issues)} issues.")
    
    for issue in issues:
        create_github_issue(issue)

if __name__ == '__main__':
    main()
