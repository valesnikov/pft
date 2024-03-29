package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var portFlag = &cli.StringFlag{
	Name:    "port",
	Aliases: []string{"p"},
	Value:   "29192", //randomly selected port, usually free
	Usage:   "network port for transmission",
}

var addrFlag = &cli.StringFlag{
	Name:    "address",
	Aliases: []string{"a"},
	Value:   "",
	Usage:   "network address for transmission",
}

var destDirFlag = &cli.StringFlag{
	Name:    "destdir",
	Aliases: []string{"d"},
	Value:   ".",
	Usage:   "saving directory",
}

var bufferFlag = &cli.StringFlag{
	Name:    "buffer-size",
	Aliases: []string{"b"},
	Value:   "256K",
	Usage:   "r/w buffer size",
}

func main() {
	cmd := &cli.App{
		Name:      "pft",
		Usage:     "TCP file sender/receiver",
		UsageText: "pft command [command options] [files...]",
		Version:   "v0.4.1",
		Flags:     []cli.Flag{bufferFlag},
		Commands: []*cli.Command{
			{
				Name:      "hs",
				Aliases:   []string{"sh", "host-send"},
				Usage:     "sending files as a host",
				UsageText: "pft hs [options] [files...]",
				Action:    HostSend,
				Flags:     []cli.Flag{portFlag},
			},
			{
				Name:      "hr",
				Aliases:   []string{"rh", "host-receive"},
				Usage:     "receiving files as a host",
				UsageText: "pft hr [options]",
				Action:    HostReceive,
				Flags:     []cli.Flag{portFlag, destDirFlag},
			},
			{
				Name:      "cs",
				Aliases:   []string{"sc", "client-send"},
				Usage:     "sending files as a client",
				UsageText: "pft cs [options] [files...]",
				Action:    ClientSend,
				Flags:     []cli.Flag{portFlag, addrFlag},
			},
			{
				Name:      "cr",
				Aliases:   []string{"rc", "client-receive"},
				Usage:     "receiving files as a client",
				UsageText: "pft cr [options]",
				Action:    ClientReceive,
				Flags:     []cli.Flag{portFlag, addrFlag, destDirFlag},
			},
		},
	}
	if err := cmd.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func HostSend(ctx *cli.Context) error {
	files := make([]string, ctx.Args().Len())
	for i := 0; i < ctx.Args().Len(); i++ {
		files[i] = ctx.Args().Get(i)
	}

	ln, err := net.Listen("tcp", ":"+ctx.String("port"))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v:%v\n", getLocalIPs(), ctx.String("port"))

	conn, err := ln.Accept()
	if err != nil {
		return err
	} else {
		defer conn.Close()
		size, err := bufSizeToNum(ctx.String("buffer-size"))
		if err != nil {
			return err
		}
		return sendFiles(files, conn, size)
	}

}

func HostReceive(ctx *cli.Context) error {
	ln, err := net.Listen("tcp", ":"+ctx.String("port"))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Start listener on %v:%v\n", getLocalIPs(), ctx.String("port"))

	conn, err := ln.Accept()
	if err != nil {
		return err
	} else {
		defer conn.Close()
		size, err := bufSizeToNum(ctx.String("buffer-size"))
		if err != nil {
			return err
		}
		return getFiles(ctx.String("destdir"), conn, size)
	}
}

func ClientSend(ctx *cli.Context) error {
	files := make([]string, ctx.Args().Len())
	for i := 0; i < ctx.Args().Len(); i++ {
		files[i] = ctx.Args().Get(i)
	}

	fmt.Printf("Awaiting connection to %v:%v", ctx.String("address"), ctx.String("port"))
	fmt.Println("")
	RETRY:
	conn, err := net.Dial("tcp", ctx.String("address")+":"+ctx.String("port"))
	if err != nil {
		cleanLine()
		fmt.Print(err)
		time.Sleep(250 * time.Millisecond)
		goto RETRY
	}
	defer conn.Close()

	defer conn.Close()
	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}
	return sendFiles(files, conn, size)
}

func ClientReceive(ctx *cli.Context) error {
	fmt.Printf("Awaiting connection to %v:%v", ctx.String("address"), ctx.String("port"))
	fmt.Println("")
	RETRY:
	conn, err := net.Dial("tcp", ctx.String("address")+":"+ctx.String("port"))
	if err != nil {
		cleanLine()
		fmt.Print(err)
		time.Sleep(250 * time.Millisecond)
		goto RETRY
	}
	defer conn.Close()

	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}
	return getFiles(ctx.String("destdir"), conn, size)
}
