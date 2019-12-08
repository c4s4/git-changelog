package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	help = `git-changelog [-help] [-version] [-file changelog]
Print markdown changelog from git logs:
-help           To print this help
-version        To print version
-file changelog To write changelog in given file`
	dateFormat = "Mon Jan 2 15:04:05 2006 -0700"
)

// Version is the version
var Version = "UNKNOWN"

// RegexpCommit is the regexp for commits
var RegexpCommit = regexp.MustCompile(`commit\s+(\w{40})\s*(\(.*?\))?(\n.*)+?\nDate:\s+(.*?)\n\n\s+(.*?)\n`)

// RegexpVersion is the regexp for versions
var RegexpVersion = regexp.MustCompile(`^(v|V)?\d+.*$`)

// commit hols a commit data
type commit struct {
	ID      string
	Tags    []string
	Date    string
	Message string
}

// parseCommandLine parses command line and returns:
// - help: a boolean that tells if we print help
// - version: a boolean that tells if we print version
// - file: a string that is the output file
func parseCommandLine() (*bool, *bool, *string) {
	help := flag.Bool("help", false, "Print help")
	version := flag.Bool("version", false, "Print version")
	file := flag.String("file", "", "Output file")
	flag.Parse()
	return help, version, file
}

// gitLogsDecorate returns git logs with tags:
// - git logs as a string
// - error if any
func gitLogsDecorate() ([]byte, error) {
	cmd := exec.Command("git", "log", "--decorate")
	return cmd.CombinedOutput()
}

// parseGitLogs parses logs:
// - logs as []byte
// Return:
// - list of commits
// - error if any
func parseGitLogs(logs []byte) ([]commit, error) {
	entries := RegexpCommit.FindAllSubmatch(logs, -1)
	var commits = make([]commit, len(entries))
	for index, entry := range entries {
		var tags []string
		if len(entry[2]) > 0 {
			noParens := string(entry[2][1 : len(entry[2])-1])
			parts := strings.Split(noParens, ", ")
			for _, part := range parts {
				if strings.HasPrefix(part, "tag: ") {
					tags = append(tags, part[5:])
				}
			}
		}
		date, err := time.Parse(dateFormat, string(entry[4]))
		if err != nil {
			return nil, err
		}
		iso := date.Format(time.RFC3339)
		commits[index] = commit{
			ID:      string(entry[1]),
			Tags:    tags,
			Date:    iso[:10],
			Message: string(entry[5]),
		}
	}
	return commits, nil
}

// generateMarkdown generated markdown changelog from commits:
// - commits
// Return
// - changelog in markdown format as string
func generateMarkdown(commits []commit) string {
	var builder strings.Builder
	builder.WriteString("# Changelog\n\n")
	for _, commit := range commits {
		version := ""
		for _, tag := range commit.Tags {
			if RegexpVersion.MatchString(tag) {
				version = tag
			}
		}
		if version != "" {
			builder.WriteString("\n## Release " + version + " (" + commit.Date + ")\n\n")
		}
		builder.WriteString("- " + commit.Message + "\n")
	}
	return builder.String()
}

func main() {
	help, version, file := parseCommandLine()
	if *help {
		fmt.Println(help)
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	logs, err := gitLogsDecorate()
	if err != nil {
		println(strings.TrimSpace(string(logs)))
		os.Exit(1)
	}
	commits, err := parseGitLogs(logs)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}
	markdown := generateMarkdown(commits)
	if *file == "" {
		fmt.Print(markdown)
	} else {
		ioutil.WriteFile(*file, []byte(markdown), 0644)
	}
}
