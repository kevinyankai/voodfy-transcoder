package powergate

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multiaddr"
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
	ma, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/5002")
	if err != nil {
		// log error
		return "", "", err
	}
	client, err := pow.NewClient(ma, grpc.WithInsecure())
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
	ma, err := multiaddr.NewMultiaddr(address)
	checkErr(err)

	var opts []grpc.DialOption
	auth := pow.TokenAuth{}
	opts = append(opts, grpc.WithPerRPCCredentials(auth))
	client, err := pow.NewClient(ma, grpc.WithInsecure(), grpc.WithPerRPCCredentials(auth))
	checkErr(err)
	return client
}

// FFSDefaultConfig show the default config
func FFSDefaultConfig(token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	client := FFSAuthenticate(token, "/ip4/127.0.0.1/tcp/5002")
	conf, err := client.FFS.DefaultConfig(authCtx(ctx, token))
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

	options := []client.PushConfigOption{}
	config := ffs.CidConfig{Cid: c}
	options = append(options, client.WithCidConfig(config))

	jid, err := fClient.FFS.PushConfig(authCtx(ctx, token), c, options...)
	checkErr(err)
	return jid.String()
}
