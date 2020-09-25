package powergate

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/ipfs/go-cid"
	"github.com/textileio/powergate/api/client"
	pow "github.com/textileio/powergate/api/client"
	"github.com/textileio/powergate/ffs"
	"google.golang.org/grpc"
)

var (
	fcClient *client.Client
)

// FFSCreate return an instance of pow
func FFSCreate() (string, string, error) {
	client, err := pow.NewClient("127.0.0.1:5002", grpc.WithInsecure())
	if err != nil {
		// log error
		return "", "", err
	}

	ffsID, ffsToken, err := client.FFS.Create(context.Background())
	if err != nil {
		// log error
		return "", "", err
	}
	err = client.Close()
	if err != nil {
		// log error
	}
	return ffsID, ffsToken, err
}

// FFSAuthenticate return a authenticate client
func FFSAuthenticate(token, address string) *pow.Client {
	var opts []grpc.DialOption
	auth := pow.TokenAuth{}
	opts = append(opts, grpc.WithPerRPCCredentials(auth))
	client, err := pow.NewClient(address, grpc.WithPerRPCCredentials(auth))
	checkErr(err)
	return client
}

// FFSDefaultConfig show the default config
func FFSDefaultConfig(token, addr string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	client := FFSAuthenticate(token, addr)
	conf, err := client.FFS.DefaultStorageConfig(authCtx(ctx, token))
	checkErr(err)
	log.Println(fmt.Sprintf("Configuration %v", conf))
}

// FFSPush send the cid to powergate setup the hot and cold configuration
func FFSPush(cidHash, token, address string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	fClient := FFSAuthenticate(token, address)

	c, err := cid.Parse(cidHash)
	checkErr(err)

	options := []client.PushStorageConfigOption{}
	config := ffs.StorageConfig{}
	options = append(options, client.WithOverride(true))
	options = append(options, client.WithStorageConfig(config))

	jid, err := fClient.FFS.PushStorageConfig(authCtx(ctx, token), c, options...)
	checkErr(err)
	return jid.String()
}
