package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
)

type CommandFunc func() (string, error)

type Commands struct {
	CurrentBranch       CommandFunc
	CreatePR            func(base, head, title, body string) (string, error)
	CheckOutMain        CommandFunc
	PullMainBranch      CommandFunc
	CreateFeatureBranch func(branchName string) (string, error)
	PushBranchToOrigin  func(branchName string) (string, error)
}

var commands = Commands{
	CurrentBranch: func() (string, error) {
		cmd := exec.Command("git", "branch", "--show-current")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(output), nil
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
}

const MAIN string = "main"

var baseBranches = []string{"main", "stage", "dev"}
var newFeatureOpts = []string{"-f", "feature"}
var helpOpts = []string{"-h", "help"}

func main() {
	ship()
}

func ship() {
	first_arg := os.Args[1]

	if slices.Contains(helpOpts, first_arg) {
		usageMessage()
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
			createPrs(currentBranch, targetBranch)
		}
		fmt.Println("All Prs created")
	}

}

func usageMessage() {
	fmt.Println("usage: ship [options]")

	fmt.Println("To see help, run `ship -h` or `ship help`")
	fmt.Println("Options:")
	fmt.Println("  -f <branch-name>              Start a new feature from main")
	fmt.Println("  prs [branch-name] [-s target] Create PRs (uses current branch if not specified)")
	fmt.Println("    -s <target-branch>          Create PR to specific branch only")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  ship -n mikun/my-feature       # creates from main")
	fmt.Println("  ship prs mikun/my-feature      # creates 3 PRs to main, stage, dev")
	fmt.Println("  ship prs                       # uses current branch, creates 3 PRs")
	fmt.Println("  ship prs -s stage              # PR current branch to stage only")
	fmt.Println("  ship prs mikun/my-feature -s stage  # PR specified branch to stage")
}

func createPrs(baseBranch string, currentBranch string) {
	fmt.Printf("Pr created: %s<=>%s \n", baseBranch, currentBranch)
}

func createFeatureBranch(branchName string) {
	fmt.Printf("📦 Starting feature workflow: %s (from %s)\n", branchName, MAIN)
	fmt.Println("")

	// Checkout main branch
	output, err := commands.CheckOutMain()
	if err != nil {
		fmt.Printf("Error checking out main: %s\n", output)
		return
	}

	// Pull latest from main
	output, err = commands.PullMainBranch()
	if err != nil {
		fmt.Printf("Error pulling main: %s\n", output)
		return
	}

	// Create feature branch
	output, err = commands.CreateFeatureBranch(branchName)
	if err != nil {
		fmt.Printf("Error creating feature branch: %s\n", output)
		return
	}

	// Push branch to origin
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