package apr

import (
	"os/exec"
)

type aprInterfacer interface {
	pull() ([]byte, error)
	build() ([]byte, error)
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

func (a *apr) build() ([]byte, error) {
	// building
	_, err := exec.Command("go", "build").Output()
	if err != nil {
		return nil, err
	}
	// return out, nil
	return []byte("build successful"), nil
}

func (a *apr) restart() ([]byte, error) {
	// restarting
	// out, err := exec.Command("/bin/sh", "cmd_restart.sh").Output()
	_, err := exec.Command("systemctl", "restart", "ajudge.service").Output()
	if err != nil {
		return nil, err
	}
	// return out, nil
	return []byte("server restarted successfully"), nil
}

func newAprService() aprInterfacer {
	return &apr{}
}
