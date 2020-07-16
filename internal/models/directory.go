package models

import (
	"encoding/json"
	"fmt"
	"log"
)

// Directory struct used to bind a directory
type Directory struct {
	ID        string    `json:"id"`
	CID       string    `json:"cid"`
	Resources Resources `json:"resources"`
}

// Resources array of resource
type Resources []Resource

// Resource struct used to bind a resource
type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CID  string `json:"cid"`
}

// MarshalBinary retrieve resource from binary
func (d *Directory) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

// UnmarshalBinary bind directory save on redis
func (d *Directory) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}

	return nil
}

// Save add directory to redis
func (d *Directory) Save() {
	m, err := d.MarshalBinary()

	if err != nil {
		log.Println("err", err)
	}

	if err := db.Redis.Set(fmt.Sprintf("directory_%s", d.ID), m, 0).Err(); err != nil {
		fmt.Printf("Unable to store example struct into redis due to: %s \n", err)
	}
}

// Get return a directory save on redis
func (d *Directory) Get() {
	InitDB()

	cacheData, cacheErr := db.Redis.Get(fmt.Sprintf("directory_%s", d.ID)).Result()

	if cacheErr == nil {
		if err := d.UnmarshalBinary([]byte(cacheData)); err != nil {
			fmt.Printf("Unable to unmarshal data into the new example struct due to: %s \n", err)
		}
	}

	defer db.Redis.Close()
}
