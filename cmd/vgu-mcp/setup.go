package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/tuanha/vgu-mcp/internal/config"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

const defaultMoodleURL = "https://moodle.vgu.edu.vn"

// runSetup runs the interactive one-time credential setup flow:
//  1. Prompts for Moodle URL, username, and password (hidden).
//  2. Attempts automatic token exchange via login/token.php.
//  3. Falls back to manual token paste if automatic login fails.
//  4. Verifies the token with GetSiteInfo and saves to config.json.
func runSetup() {
	ctx := context.Background()
	in := bufio.NewReader(os.Stdin)

	// Step 1: Moodle URL
	fmt.Fprintf(os.Stderr, "Moodle URL [%s]: ", defaultMoodleURL)
	moodleURL := readLine(in)
	if moodleURL == "" {
		moodleURL = defaultMoodleURL
	}

	// Step 2: Username
	fmt.Fprint(os.Stderr, "Username: ")
	username := readLine(in)
	if username == "" {
		log.Fatal("username cannot be empty")
	}

	// Step 3: Password (hidden — does not echo to terminal)
	fmt.Fprint(os.Stderr, "Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr) // newline after hidden input
	if err != nil {
		log.Fatalf("read password: %v", err)
	}
	password := string(passwordBytes)

	// Step 4a: Attempt automatic token exchange
	var token string
	token, err = moodle.Login(ctx, moodleURL, username, password)
	if err != nil {
		// Step 4b: Fallback — prompt user to paste token manually
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Automatic login failed: %v\n", err)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Please paste your Moodle web service token instead.")
		fmt.Fprintln(os.Stderr, "(Moodle → Profile → Preferences → Security keys)")
		fmt.Fprint(os.Stderr, "Token: ")
		token = readLine(in)
		if token == "" {
			log.Fatal("no token provided — setup aborted")
		}
	}

	// Step 5: Verify token and retrieve user info
	client := moodle.NewClient(moodleURL, token)
	info, err := client.GetSiteInfo(ctx)
	if err != nil {
		log.Fatalf("token verification failed: %v\nDouble-check your credentials and try again.", err)
	}
	fmt.Fprintf(os.Stderr, "\n✓ Login successful. Welcome, %s!\n", info.Fullname)

	// Step 6: Persist config — no password ever written to disk
	cfg := &config.Config{
		MoodleURL:   moodleURL,
		MoodleToken: token,
		UserID:      info.UserID,
	}
	if err := config.Save(cfg); err != nil {
		log.Fatalf("save config: %v", err)
	}
	path, _ := config.Path()
	fmt.Fprintf(os.Stderr, "✓ Credentials saved to %s\n", path)
	fmt.Fprintln(os.Stderr, "\nYou're all set. Configure your AI client to run: vgu-mcp")
}

// readLine reads one line from r, trimming whitespace.
func readLine(r *bufio.Reader) string {
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}
