package task

import (
	"github.com/RichardKnop/machinery/v1"
)

// ManagerTranscoder managment of the task
func ManagerTranscoder(kind, resourceID, resourceName, directory, tracker string, server *machinery.Server) error {
	var err error
	if kind == "livepeer" {
		LivepeerChain(resourceID, directory, tracker, resourceName, server)
	}

	if kind == "local" {
		Local(resourceID, resourceName, directory, tracker, server)
	}

	return err
}

// ManagerIPFS managment of the task
func ManagerIPFS(resourceID, directory, tracker string, server *machinery.Server) {
	IPFSAddDir(resourceID, directory, tracker, server)
}
