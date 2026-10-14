#!/usr/bin/env bash
# Validate that a commit message follows the Conventional Commits format.
# Called by lefthook on the commit-msg hook.
#
# Format:
#   <type>(<scope>)<!>: <subject>
#
# Where:
#   - type is one of: feat, fix, docs, style, refactor, perf, test, build,
#     ci, chore, revert
#   - scope is optional, in parentheses
#   - ! is optional, marks a BREAKING CHANGE
#   - subject is a short description (no trailing period), max 72 chars
#
# Body and footer are optional but must be separated from the subject by a
# blank line if present. A "BREAKING CHANGE:" footer is also recognised
# as a breaking change marker.

set -euo pipefail

MSG_FILE="$1"
FIRST_LINE=$(head -1 "$MSG_FILE")

# Allowed types per conventionalcommits.org
PATTERN='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-zA-Z0-9_.-]+\))?!?: .+'

# Allow merge / revert / fixup commits that some workflows generate
# automatically — they're not user-authored.
if [[ "$FIRST_LINE" =~ ^Merge\  ]] || \
   [[ "$FIRST_LINE" =~ ^Revert\ \" ]] || \
   [[ "$FIRST_LINE" =~ ^(fixup|squash)!\  ]] || \
   [[ "$FIRST_LINE" =~ ^Initial\ commit$ ]]; then
    exit 0
fi

if ! [[ "$FIRST_LINE" =~ $PATTERN ]]; then
    cat >&2 <<EOF
ERROR: commit message does not follow Conventional Commits format.

Got:      ${FIRST_LINE:0:80}
Expected: <type>(<scope>)<!>: <subject>

Types:    feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert

Examples:
  feat: add Jira issue deep link
  fix(teams): preserve commas in users parameter
  docs: update README with new shortcode example
  feat!: drop Go 1.22 support
  feat(api)!: change IssueParams.IssueType to enum
EOF
    exit 1
fi

# Enforce subject length (best practice from the Conventional Commits spec)
SUBJECT="${FIRST_LINE#*: }"
if [[ ${#SUBJECT} -gt 72 ]]; then
    echo "ERROR: commit subject is ${#SUBJECT} characters; max is 72." >&2
    echo "  Subject: $SUBJECT" >&2
    exit 1
fi

# Subject must not end with a period.
if [[ "$SUBJECT" =~ \.$ ]] && [[ ${#SUBJECT} -gt 1 ]]; then
    # Allow "..." truncation, "etc.", and abbreviations.
    if ! [[ "$SUBJECT" =~ \.\.\.$ ]] && \
       ! [[ "$SUBJECT" =~ (etc|vs|i\.e|e\.g|approx)\.$ ]]; then
        echo "ERROR: commit subject should not end with a period." >&2
        echo "  Subject: $SUBJECT" >&2
        exit 1
    fi
fi

exit 0
