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
	// Help as printed with -help option
	Help = `git changelog [-help] [-version] [-file changelog] [-release regexp]
Print markdown changelog from git logs:
-help           To print this help
-version        To print version
-file changelog To write changelog in given file
-tag regexp     To set regexp for release tags (defaults to "^(v|V)?\d+.*$")
-nodate         To omit dates in releases titles`
	dateFormat       = "Mon Jan 2 15:04:05 2006 -0700"
	defaultRegexpTag = `^(v|V)?\d+.*$`
)

// Version is the version
var Version = "UNKNOWN"

// RegexpCommit is the regexp for commits
var RegexpCommit = regexp.MustCompile(`commit\s+(\w{40})\s*(\(.*?\))?(\n.*)+?\nDate:\s+(.*?)\n\n\s+(.*?)\n`)

// RegexpTag is the regexp for version tags
var RegexpTag *regexp.Regexp

// commit hols a commit data
type commit struct {
	ID      string
	Tags    []string
	Date    string
	Message string
}

// Version returns version in commit of empty string
func (c *commit) Version() string {
	version := ""
	for _, tag := range c.Tags {
		if RegexpTag.MatchString(tag) {
			version = tag
		}
	}
	return version
}

// parseCommandLine parses command line and returns:
// - help: a boolean that tells if we print help
// - version: a boolean that tells if we print version
// - file: a string that is the output file
// - tag: release tag regexp
// - nodate: a boolean that tells if dates should be omitted
func parseCommandLine() (*bool, *bool, *string, *string, *bool) {
	help := flag.Bool("help", false, "Print help")
	version := flag.Bool("version", false, "Print version")
	file := flag.String("file", "", "Output file")
	tag := flag.String("tag", defaultRegexpTag, "Regexp for release tags")
	nodate := flag.Bool("nodate", false, "Omit date in release titles")
	flag.Parse()
	return help, version, file, tag, nodate
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
func generateMarkdown(commits []commit, nodate bool) string {
	var builder strings.Builder
	builder.WriteString("# Changelog\n")
	if len(commits) > 0 && commits[0].Version() == "" {
		builder.WriteString("\n")
	}
	for _, commit := range commits {
		version := commit.Version()
		if version != "" {
			if nodate {
				builder.WriteString("\n## Release " + version + "\n\n")
			} else {
				builder.WriteString("\n## Release " + version + " (" + commit.Date + ")\n\n")
			}
		}
		builder.WriteString("- " + commit.Message + "\n")
	}
	return builder.String()
}

func main() {
	help, version, file, release, nodate := parseCommandLine()
	if *help {
		fmt.Println(Help)
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	var err error
	RegexpTag, err = regexp.Compile(*release)
	if err != nil {
		println(fmt.Sprintf("Error compiling release tags regexp: %v", err))
		os.Exit(3)
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
	markdown := generateMarkdown(commits, *nodate)
	if *file == "" {
		fmt.Print(markdown)
	} else {
		ioutil.WriteFile(*file, []byte(markdown), 0644)
	}
}
