package task

import (
	"fmt"
	"log"
	"strings"

	"github.com/RichardKnop/machinery/v1"
	"github.com/Voodfy/voodfy-transcoder/internal/models"
	"github.com/Voodfy/voodfy-transcoder/internal/utils"
	"github.com/Voodfy/voodfy-transcoder/pkg/powergate"
	"github.com/Voodfy/voodfy-transcoder/pkg/voodfyapi"
	"github.com/google/uuid"
)

// ManagerTranscoder managment of the task
func ManagerTranscoder(kind, resourceID, resourceName, directory, tracker string, server *machinery.Server) error {
	var err error
	if kind == "remote" {
		log.Println("not available")
	}

	if kind == "local" {
		Local(resourceID, resourceName, directory, tracker, server)
	}

	return err
}

// ManagerIPFS managment of the task
func ManagerIPFS(directory, resourceID string, server *machinery.Server) {
	IPFSAddDir(directory, resourceID, server)
}

// ManagerPowergate managment of the task that will use the powergate
func ManagerPowergate(directoryID string) string {
	api := voodfyapi.NewClient()
	pow, err := api.Powergate("", false)

	if err != nil {
		utils.SendError("voodfycli.tasks.manager.ManagerPowergate", err)
		return "Error to retrieve powergate instance, try again!"
	}

	directory := models.Directory{
		ID: directoryID,
	}

	directory.Get()

	for idx, r := range directory.Resources {
		jid := powergate.FFSPush(r.CID, pow.Token, pow.Address)
		r.Jid = jid
		directory.Resources[idx] = r
	}

	directory.Save()

	return "Stored, now you can verify the status of the job!"
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
