ship() {
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        echo "Feature workflow CLI"
        echo ""
        echo "Usage: ship [options]"
        echo ""
        echo "Options:"
        echo "  -n <branch-name>              Start a new feature from main"
        echo "  prs [branch-name] [-s target] Create PRs (uses current branch if not specified)"
        echo "    -s <target-branch>          Create PR to specific branch only"
        echo ""
        echo "Examples:"
        echo "  ship -n mikun/my-feature       # creates from main"
        echo "  ship prs mikun/my-feature      # creates 3 PRs to main, stage, dev"
        echo "  ship prs                       # uses current branch, creates 3 PRs"
        echo "  ship prs -s stage              # PR current branch to stage only"
        echo "  ship prs mikun/my-feature -s stage  # PR specified branch to stage"
        return
    fi

    if [ "$1" = "prs" ]; then
        local branch_name=""
        local target_branch=""
        
        # Parse arguments
        shift  # skip "prs"
        while [ $# -gt 0 ]; do
            case "$1" in
                -s)
                    target_branch="$2"
                    shift 2
                    ;;
                *)
                    branch_name="$1"
                    shift
                    ;;
            esac
        done
        
        # If no branch name specified, use current branch
        if [ -z "$branch_name" ]; then
            branch_name=$(git rev-parse --abbrev-ref HEAD)
            if [ -z "$branch_name" ] || [ "$branch_name" = "HEAD" ]; then
                echo "Error: Could not determine branch name"
                return 1
            fi
        fi
        
        # If target specified, create single PR, otherwise create all PRs
        if [ -n "$target_branch" ]; then
            echo "Creating pull request: $branch_name -> $target_branch"
            if gh pr create \
                --base "$target_branch" \
                --head "$branch_name" \
                --title "$branch_name" \
                --body "$branch_name"; then
                echo "✓ PR to $target_branch created"
            else
                echo "✗ Failed to create PR to $target_branch"
                return 1
            fi
        else
            echo "Creating pull requests for: $branch_name"
            echo ""
            for target in main stage dev; do
                echo ""
                if gh pr create \
                    --base "$target" \
                    --head "$branch_name" \
                    --title "$branch_name" \
                    --body "$branch_name"; then
                    echo "✓ PR to $target created"
                else
                    echo "✗ Failed to create PR to $target"
                fi
            done
            echo ""
            echo "✓ All PRs created!"
        fi
        return
    fi

    local branch_name=""
    local base_branch="main"

    if [ "$1" = "-n" ]; then
        if [ -z "$2" ]; then
            echo "Error: -n flag requires a branch name"
            return 1
        fi
        branch_name="$2"
    elif [ -z "$1" ]; then
        read "?Enter branch name: " branch_name
    else
        echo "Error: Invalid option. Use -h for help."
        return 1
    fi

    if [ -z "$branch_name" ]; then
        echo "Error: Branch name required"
        return 1
    fi

    echo "📦 Starting feature workflow: $branch_name (from $base_branch)"
    echo ""

    git checkout "$base_branch" || return 1
    git pull origin "$base_branch" || return 1
    git checkout -b "$branch_name" || return 1
    git push -u origin "$branch_name" || return 1
    echo "✓ Branch created and pushed"
    echo ""
    echo "Make your changes, commit, and push. Then run:"
    echo "  ship prs"
}