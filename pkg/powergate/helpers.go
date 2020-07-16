package powergate

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/textileio/powergate/api/client"
)

func checkErr(e error) {
	if e != nil {
		Fatal(e)
	}
}

// Fatal prints a fatal error to stdout, and exits immediately with
// error code 1.
func Fatal(err error, args ...interface{}) {
	log.Println(err.Error())
	words := strings.SplitN(err.Error(), " ", 2)
	words[0] = strings.Title(words[0])
	msg := strings.Join(words, " ")
	fmt.Println(aurora.Sprintf(aurora.Red("> Error! %s"),
		aurora.Sprintf(aurora.BrightBlack(msg), args...)))
	os.Exit(1)
}

func authCtx(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, client.AuthKey, token)
}
