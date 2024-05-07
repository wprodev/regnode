package node

import (
	"github.com/rs/zerolog/log"
)

func RegisterWithCommand(cmd string, controlPlane, etcd, worker bool) (string, error) {
	if controlPlane {
		cmd = cmd + " --controlplane"
	}
	if etcd {
		cmd = cmd + " --etcd"
	}
	if worker {
		cmd = cmd + " --worker"
	}
	log.Info().Msgf("Executing shell command: %s", cmd)
	return "", nil
	// shellCmd := exec.Command("bash", "-c", preparedCmd)
	// output, err := shellCmd.CombinedOutput()
	// if err != nil {
	// 	return string(output), err
	// }
	//return string(output), nil
}
