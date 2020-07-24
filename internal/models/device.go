package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Voodfy/voodfy-transcoder/internal/utils"
)

// Device struct used to store informations during the use of voodfycli
type Device struct {
	UUID       string `json:"uuid"`
	Token      string `json:"token"`
	SecretHash string `json:"secretHash"`
}

// existDevice verify if exist device registered
func existDevice() bool {
	cacheData, _ := db.Redis.Keys("*").Result()

	for _, key := range cacheData {
		return strings.Contains(key, "device")
	}

	return false
}

// Save add device to redis
func (d *Device) Save() bool {
	InitDB()

	if existDevice() {
		return false
	}

	d.SecretHash = createHash(d.UUID)
	m, err := d.MarshalBinary()

	if err != nil {
		log.Println("err", err)
	}

	if err := db.Redis.Set(fmt.Sprintf("device_%s", string(d.UUID)), m, 0).Err(); err != nil {
		fmt.Printf("Unable to store example struct into redis due to: %s \n", err)
	}
	return true
}

// Update update device to redis
func (d *Device) Update() {
	InitDB()

	m, err := d.MarshalBinary()

	if err != nil {
		log.Println("err", err)
	}

	if err := db.Redis.Set(fmt.Sprintf("device_%s", string(d.UUID)), m, 0).Err(); err != nil {
		fmt.Printf("Unable to store example struct into redis due to: %s \n", err)
	}
}

// Get return a device save on redis
func (d *Device) Get() {
	var key string
	InitDB()

	keys, _ := db.Redis.Keys("*").Result()

	for _, k := range keys {
		if strings.Contains(k, "device_") {
			key = k
		}
	}

	cacheData, cacheErr := db.Redis.Get(key).Result()

	if cacheErr == nil {
		if err := d.UnmarshalBinary([]byte(cacheData)); err != nil {
			fmt.Printf("Unable to unmarshal data into the new example struct due to: %s \n", err)
		}
	}
}

// GetBySecretHash return a device save on redis
func GetBySecretHash(secret string) (device Device, ok bool) {
	var key string
	InitDB()
	keys, _ := db.Redis.Keys("*").Result()

	for _, k := range keys {
		if strings.Contains(k, "device_") {
			key = k
		}
	}

	cacheData, cacheErr := db.Redis.Get(key).Result()

	if cacheErr == nil {
		if err := device.UnmarshalBinary([]byte(cacheData)); err != nil {
			fmt.Printf("Unable to unmarshal data into the new example struct due to: %s \n", err)
		}
	}

	if strings.TrimSpace(secret) == device.SecretHash {
		ok = true
		return device, true
	}

	return Device{}, ok
}

// MarshalBinary retrieve resource from binary
func (d *Device) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

// UnmarshalBinary bind device save on redis
func (d *Device) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}

	return nil
}

// ToSignup return map to used on signup
func (d *Device) ToSignup() map[string]interface{} {
	return map[string]interface{}{
		"device":   d.UUID,
		"password": d.SecretHash,
	}
}

func createHash(id string) string {
	key := fmt.Sprintf("%s%s", id, utils.RandSeq(256))
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
