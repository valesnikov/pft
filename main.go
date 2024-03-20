package main

import (
	"fmt"
	"net"
	"os"

	"github.com/urfave/cli/v2"
)

const TransmissionBufferSize int = 1024 * 1024 //1MiB

const DefaultPort = "29192" //randomly selected port, usually free

var portFlag = &cli.StringFlag{
	Name:    "port",
	Aliases: []string{"p"},
	Value:   DefaultPort,
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

var fileFlag = &cli.StringFlag{
	Name:     "file",
	Aliases:  []string{"f"},
	Required: true,
	Usage:    "file for writing or reading",
}

func main() {
	cmd := &cli.App{
		Name:      "pft",
		Usage:     "TCP file sender/receiver",
		UsageText: "pft command [command options] [files...]",
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

			{
				Name:      "fs",
				Aliases:   []string{"sf", "file-send"},
				Usage:     "sending files to file",
				UsageText: "pft fs [options] [files...]",
				Action:    FileSend,
				Flags:     []cli.Flag{fileFlag},
			},
			{
				Name:      "fr",
				Aliases:   []string{"rf", "file-receive"},
				Usage:     "receiving files from file",
				UsageText: "pft fr [options]",
				Action:    FileReceive,
				Flags:     []cli.Flag{fileFlag, destDirFlag},
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

	fmt.Printf("Start listener on %v port\n", ctx.String("port"))

	conn, err := ln.Accept()
	if err != nil {
		return err
	} else {
		return sendFiles(files, conn)
	}

}

func HostReceive(ctx *cli.Context) error {
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
		return getFiles(ctx.String("destdir"), conn)
	}
}

func ClientSend(ctx *cli.Context) error {
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

func ClientReceive(ctx *cli.Context) error {
	conn, err := net.Dial("tcp", ctx.String("address")+":"+ctx.String("port"))
	if err != nil {
		return err
	}
	return getFiles(ctx.String("destdir"), conn)
}

func FileSend(ctx *cli.Context) error {
	files := make([]string, ctx.Args().Len())
	for i := 0; i < ctx.Args().Len(); i++ {
		files[i] = ctx.Args().Get(i)
	}
	file, err := os.Create(ctx.String("file"))
	if err != nil {
		return err
	}
	return sendFiles(files, newPftWriter(file))
}

func FileReceive(ctx *cli.Context) error {
	file, err := os.Open(ctx.String("file"))
	if err != nil {
		return err
	}
	return getFiles(ctx.String("destdir"), newPftReader(file))
}
