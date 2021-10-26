package apr

import (
	"os/exec"
)

type aprInterfacer interface {
	pull() ([]byte, error)
	restart() ([]byte, error)
}

type apr struct{}

func (a *apr) pull() ([]byte, error) {
	// pulling
	out, err := exec.Command("git", "pull").Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (a *apr) restart() ([]byte, error) {
	// restarting
	out, err := exec.Command("/bin/sh", "cmd.sh").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func newAprService() aprInterfacer {
	return &apr{}
}
