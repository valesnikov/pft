package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
	"os"
)

const DEFAULT_PORT = "29192"

var PORT_FLAG = &cli.StringFlag{
	Name:    "port",
	Aliases: []string{"p"},
	Value:   DEFAULT_PORT,
	Usage:   "network port for transmission",
}

var ADDR_FLAG = &cli.StringFlag{
	Name:    "address",
	Aliases: []string{"a"},
	Value:   "",
	Usage:   "network address for transmission",
}

func main() {
	cmd := &cli.App{
		Name:                 "pft",
		Usage:                "TCP file sender/receiver",
		EnableBashCompletion: true,

		Commands: []*cli.Command{
			{
				Name:    "hs",
				Aliases: []string{"sh"},
				Usage:   "sending files as a host",
				Action:  hostSend,
				Flags:   []cli.Flag{PORT_FLAG},
			},

			{
				Name:    "hr",
				Aliases: []string{"rh"},
				Usage:   "receiving files as a host",
				Action:  hostReceive,
				Flags:   []cli.Flag{PORT_FLAG},
			},

			{
				Name:    "cs",
				Aliases: []string{"sc"},
				Usage:   "sending files as a client",
				Action:  clientSend,
				Flags:   []cli.Flag{PORT_FLAG, ADDR_FLAG},
			},

			{
				Name:    "cr",
				Aliases: []string{"rc"},
				Usage:   "receiving files as a client",
				Action:  clientReceive,
				Flags:   []cli.Flag{PORT_FLAG, ADDR_FLAG},
			},
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func hostSend(ctx *cli.Context) error {
	files := make([]string, ctx.Args().Len())
	for i := 0; i < ctx.Args().Len(); i++ {
		files[i] = ctx.Args().Get(i)
	}

	ln, err := net.Listen("tcp", ":"+ctx.String("port"))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v port\n", ctx.String("port"))

	conn, err := ln.Accept()
	if err != nil {
		return err
	} else {
		return sendFiles(files, conn)
	}
}

func hostReceive(ctx *cli.Context) error {
	ln, err := net.Listen("tcp", ":"+ctx.String("port"))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v port\n", ctx.String("port"))

	conn, err := ln.Accept()
	if err != nil {
		return err
	} else {
		return getFiles(ctx.Args().Get(0), conn)
	}
}

func clientSend(ctx *cli.Context) error {
	files := make([]string, ctx.Args().Len())
	for i := 0; i < ctx.Args().Len(); i++ {
		files[i] = ctx.Args().Get(i)
	}

	conn, err := net.Dial("tcp", ctx.String("address")+":"+ctx.String("port"))
	if err != nil {
		return err
	}
	return sendFiles(files, conn)
}

func clientReceive(ctx *cli.Context) error {
	conn, err := net.Dial("tcp", ctx.String("address")+":"+"29192")
	if err != nil {
		return err
	}
	return getFiles(ctx.Args().Get(0), conn)
}
