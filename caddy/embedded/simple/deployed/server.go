package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/mholt/caddy"
	_ "github.com/mholt/caddy/caddyhttp"
)

func init() {
	// configure default caddyfile
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(defaultLoader))
}

func main() {
	caddy.AppName = "CaddyTest"
	caddy.AppVersion = "0.1.1"

	// load caddyfile
	caddyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		log.Fatal(err)
	}

	// start caddy server
	instance, err := caddy.Start(caddyfile)
	if err != nil {
		log.Fatal(err)
	}

	instance.Wait()
}

// provide loader function
func defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}
