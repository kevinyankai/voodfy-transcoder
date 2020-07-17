package task

import (
	"fmt"
	"strings"

	"github.com/RichardKnop/machinery/v1"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	"github.com/Voodfy/voodfy-transcoder/pkg/voodfyapi"
	"github.com/google/uuid"
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

// ManagerSetupAccountVoodfy setup an account at Voodfy
func ManagerSetupAccountVoodfy() string {
	api := voodfyapi.NewClient()
	id := uuid.New()

	device := models.Device{
		UUID: id.String(),
	}

	isCreated := device.Save()

	if isCreated {
		api.Payload = device.ToSignup()
		api.Signup()
		return fmt.Sprintf("Store with safety the secret hash %s", device.SecretHash)
	}

	return "Device exist!\nPlease use the command *login* instead *signup*"

}

// ManagerLoginVoodfy setup an account at Voodfy
func ManagerLoginVoodfy(secret string) string {
	api := voodfyapi.NewClient()
	device, ok := models.GetBySecretHash(secret)

	if !ok {
		id, err := api.Retrieve(strings.TrimSpace(secret))
		if err != nil {
			utils.SendError("voodfycli.tasks.manager.ManagerLoginVoodfy", err)
			return "Secret invalid, try again!"
		}
		device.UUID = id
		device.SecretHash = strings.TrimSpace(secret)
	}

	api.Payload = device.ToSignup()
	token, err := api.Token()

	device.Token = token
	device.Update()

	if err != nil {
		utils.SendError("voodfycli.tasks.manager.ManagerLoginVoodfy", err)
	}

	return "Logged!"
}
