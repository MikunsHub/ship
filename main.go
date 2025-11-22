package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"google.golang.org/genai"
)

type CommandFunc func() (string, error)

type Commands struct {
	CurrentBranch       CommandFunc
	CreatePR            func(base, head, title, body string) (string, error)
	CheckOutMain        CommandFunc
	PullMainBranch      CommandFunc
	CreateFeatureBranch func(branchName string) (string, error)
	PushBranchToOrigin  func(branchName string) (string, error)
	GetCommits          func(base, head string) ([]string, error)
}

var commands = Commands{
	CurrentBranch: func() (string, error) {
		cmd := exec.Command("git", "branch", "--show-current")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	},
	CreatePR: func(base, head, title, body string) (string, error) {
		cmd := exec.Command(
			"gh", "pr", "create",
			"--base", base,
			"--head", head,
			"--title", title,
			"--body", body,
		)
		output, err := cmd.CombinedOutput()
		return string(output), err
	},
	CheckOutMain: func() (string, error) {
		cmd := exec.Command(
			"git", "checkout", MAIN,
		)
		output, err := cmd.CombinedOutput()
		return string(output), err
	},
	PullMainBranch: func() (string, error) {
		cmd := exec.Command(
			"git", "pull", "origin", MAIN,
		)
		output, err := cmd.CombinedOutput()
		return string(output), err
	},
	CreateFeatureBranch: func(branchName string) (string, error) {
		cmd := exec.Command(
			"git", "checkout", "-b", branchName,
		)
		output, err := cmd.CombinedOutput()
		return string(output), err
	},
	PushBranchToOrigin: func(branchName string) (string, error) {
		cmd := exec.Command(
			"git", "push", "-u", "origin", branchName,
		)
		output, err := cmd.CombinedOutput()
		return string(output), err
	},
	GetCommits: func(base, head string) ([]string, error) {
		cmd := exec.Command(
			"git", "log", "--oneline", fmt.Sprintf("%s..%s", base, head),
		)
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		commits := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(commits) == 1 && commits[0] == "" {
			return []string{}, nil
		}
		return commits, nil
	},
}

const MAIN string = "main"

var baseBranches = []string{"main", "stage", "dev"}
var newFeatureOpts = []string{"-f", "feature"}
var helpOpts = []string{"-h", "help"}

func main() {
	ship()
}

func ship() {
	if len(os.Args) < 2 {
		showUsageMessage()
		return
	}

	first_arg := os.Args[1]

	if slices.Contains(helpOpts, first_arg) {
		showUsageMessage()
		return
	}

	if slices.Contains(newFeatureOpts, first_arg) {
		var userBranchName string
		fmt.Print("Enter branch name:")
		fmt.Scan(&userBranchName)
		fmt.Println("your branch name is: ", userBranchName)
		createFeatureBranch(userBranchName)
	}

	if first_arg == "prs" {
		var targetBranch string
		currentBranch, err := commands.CurrentBranch()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		if len(os.Args) > 2 && os.Args[2] == "-s" {
			targetBranch = os.Args[3]
		}

		fmt.Println(currentBranch)
		fmt.Println(targetBranch)

		if targetBranch == "" {
			for _, base := range baseBranches {
				createPrs(base, currentBranch)
			}
		} else {
			createPrs(targetBranch, currentBranch)
		}
		fmt.Println("All Prs created")
	}

}

func showUsageMessage() {
	fmt.Println("usage: ship [options]")

	fmt.Println("To see help, run `ship -h` or `ship help`")
	fmt.Println("Options:")
	fmt.Println("  -f <branch-name>              Start a new feature from main")
	fmt.Println("  prs [branch-name] [-s target] Create PRs (uses current branch if not specified)")
	fmt.Println("    -s <target-branch>          Create PR to specific branch only")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  ship -f mikun/my-feature       # creates from main")
	fmt.Println("  ship prs mikun/my-feature      # creates 3 PRs to main, stage, dev")
	fmt.Println("  ship prs                       # uses current branch, creates 3 PRs")
	fmt.Println("  ship prs -s stage              # PR current branch to stage only")
	fmt.Println("  ship prs mikun/my-feature -s stage  # PR specified branch to stage")
}

func createPrs(baseBranch string, currentBranch string) {
	fmt.Printf("\n🔄 Creating pull request: %s <- %s\n", baseBranch, currentBranch)

	// Get commits between branches
	fmt.Println("📝 Fetching commits...")
	commits, err := commands.GetCommits(baseBranch, currentBranch)
	if err != nil {
		fmt.Printf("Error getting commits: %v\n", err)
		return
	}

	if len(commits) == 0 {
		fmt.Printf("⚠️  No commits found between %s and %s. Skipping PR creation.\n", baseBranch, currentBranch)
		return
	}

	fmt.Printf("Found %d commit(s)\n", len(commits))

	// Generate PR body using AI
	fmt.Println("🤖 Generating PR description with AI...")
	body, err := generatePRBody(commits)
	if err != nil {
		fmt.Printf("⚠️  AI generation failed: %v\n", err)
		body = "Changes:\n" + strings.Join(commits, "\n")
	}

	// Show generated body to user
	fmt.Println("\n┌─────────────────────────────────────────────┐")
	fmt.Println("│         GENERATED PR DESCRIPTION             │")
	fmt.Println("└─────────────────────────────────────────────┘")
	fmt.Println(body)
	fmt.Println("─────────────────────────────────────────────")

	// Ask user to accept or edit
	fmt.Print("\nAccept this description? (y/n/e for edit): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "n" {
		fmt.Println("❌ PR creation cancelled")
		return
	}

	if response == "e" {
		fmt.Println("\nEnter your custom PR description (press Ctrl+D when done):")
		fmt.Println("─────────────────────────────────────────────")
		customBody := ""
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			customBody += scanner.Text() + "\n"
		}
		if customBody != "" {
			body = strings.TrimSpace(customBody)
		}
	}

	// Create PR
	title := fmt.Sprintf("Merge %s into %s", currentBranch, baseBranch)
	fmt.Printf("\n🚀 Creating PR...\n")
	output, err := commands.CreatePR(baseBranch, currentBranch, title, body)
	if err != nil {
		fmt.Printf("❌ Error creating PR: %s\n", output)
		return
	}
	fmt.Printf("✅ PR created successfully!\n%s\n", output)
}

func createFeatureBranch(branchName string) {
	fmt.Printf("📦 Starting feature workflow: %s (from %s)\n", branchName, MAIN)
	fmt.Println("")

	output, err := commands.CheckOutMain()
	if err != nil {
		fmt.Printf("Error checking out main: %s\n", output)
		return
	}

	output, err = commands.PullMainBranch()
	if err != nil {
		fmt.Printf("Error pulling main: %s\n", output)
		return
	}

	output, err = commands.CreateFeatureBranch(branchName)
	if err != nil {
		fmt.Printf("Error creating feature branch: %s\n", output)
		return
	}

	output, err = commands.PushBranchToOrigin(branchName)
	if err != nil {
		fmt.Printf("Error pushing branch to origin: %s\n", output)
		return
	}

	fmt.Println("✓ Branch created and pushed")
	fmt.Println("")
	fmt.Println("Make your changes, commit, and push. Then run:")
	fmt.Println("  ship prs")
}

func generatePRBody(commits []string) (string, error) {
	ctx := context.Background()

	// Get API key from environment variable for security
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = "AIzaSyCJ1YF56wCLfhK-FQu0u2704XbeUbUhukA" // Fallback (should use env var)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", err
	}

	commitsText := strings.Join(commits, "\n")
	prompt := fmt.Sprintf(
		`Analyze these git commits and generate a concise, professional PR description.

Commits:
%s

Please provide:
1. A brief summary (2-3 sentences) of what changed
2. Key changes as bullet points
3. Any notable implementation details

Keep it concise and focus on the "why" and "what", not the "how".`,
		commitsText,
	)

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
