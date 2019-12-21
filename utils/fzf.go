package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SelectProfileByFzf(combinedProfiles []string, pattern string) ([]byte, error) {
	joinedProfiles := strings.Join(combinedProfiles, "\n")

	fzfCommand := fmt.Sprintf("echo -e '%s' | fzf-tmux --height 30%% --reverse -1 -0 --with-nth=1 --delimiter=: --header 'Select AWS profile' --query '%s'",
		joinedProfiles,
		pattern)
	shellCommand := exec.Command("bash", "-c", fzfCommand)
	shellCommand.Stdin = os.Stdin
	shellCommand.Stderr = os.Stderr
	return shellCommand.Output()
}

