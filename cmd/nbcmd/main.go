package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mattn/go-nostrbuild"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/urfave/cli/v2"
)

var nsec string

const name = "nbcmd"

const version = "0.0.10"

var revision = "HEAD"

func init() {
	nsec = os.Getenv("NBCMD_NSEC")
}

func sign(ev *nostr.Event) error {
	var sk string
	if _, s, err := nip19.Decode(nsec); err != nil {
		return err
	} else {
		sk = s.(string)
	}
	if pub, err := nostr.GetPublicKey(sk); err == nil {
		if _, err := nip19.EncodePublicKey(pub); err != nil {
			return err
		}
		ev.PubKey = pub
	} else {
		return err
	}
	return ev.Sign(sk)
}

func doUpload(cCtx *cli.Context) error {
	verbose := cCtx.Bool("v")
	for _, arg := range cCtx.Args().Slice() {
		b, err := ioutil.ReadFile(arg)
		if err != nil {
			return err
		}
		result, err := nostrbuild.Upload(bytes.NewBuffer(b), sign)
		if err != nil {
			return err
		}
		if verbose {
			json.NewEncoder(os.Stdout).Encode(result)
		} else {
			fmt.Println(result.Data[0].URL)
		}
	}
	return nil
}

func doDelete(cCtx *cli.Context) error {
	verbose := cCtx.Bool("v")
	for _, arg := range cCtx.Args().Slice() {
		result, err := nostrbuild.Delete(arg, sign)
		if err != nil {
			return err
		}
		if verbose {
			json.NewEncoder(os.Stdout).Encode(result)
		} else {
			fmt.Println(result.Message)
		}
	}
	return nil
}

func main() {
	app := &cli.App{
		Usage:       "A cli application for nostr.build",
		Description: "A cli application for nostr.build",
		Commands: []*cli.Command{
			{
				Name:  "upload",
				Usage: "upload image files",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "v", Usage: "verbose"},
				},
				Action: doUpload,
			},
			{
				Name:  "delete",
				Usage: "delete image files",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "v", Usage: "verbose"},
				},
				Action: doDelete,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
