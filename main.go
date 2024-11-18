package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/klauspost/compress/zstd"
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

var mkdirFlag = &cli.BoolFlag{
	Name:    "mkdir",
	Aliases: []string{"m"},
	Value:   false,
	Usage:   "сreate destdir if it does not exist",
}

var zstdFlag = &cli.BoolFlag{
	Name:    "zstd",
	Aliases: []string{"z"},
	Value:   false,
	Usage:   "enables compression when sending",
}

func main() {
	cmd := &cli.App{
		Name:      "pft",
		Usage:     "TCP file sender/receiver",
		UsageText: "pft [global options] command [command options] [files...]",
		Version:   "v0.5.0-develop",
		Flags:     []cli.Flag{bufferFlag, portFlag, zstdFlag},
		Commands: []*cli.Command{
			{
				Name:      "host-send",
				Aliases:   []string{"hs", "sh"},
				Usage:     "sending files as a host",
				UsageText: "pft hs [options] [files...]",
				Action:    HostSend,
				Flags:     []cli.Flag{},
			},
			{
				Name:      "host-receive",
				Aliases:   []string{"hr", "rh"},
				Usage:     "receiving files as a host",
				UsageText: "pft hr [options]",
				Action:    HostReceive,
				Flags:     []cli.Flag{destDirFlag, mkdirFlag},
			},
			{
				Name:      "host-pipe",
				Aliases:   []string{"pipe", "p"},
				Usage:     "pipe between 2 clients",
				UsageText: "pft p [options]",
				Action:    HostPipe,
				Flags:     []cli.Flag{},
			},
			{
				Name:      "client-send",
				Aliases:   []string{"cs", "sc"},
				Usage:     "sending files as a client",
				UsageText: "pft cs [options] [files...]",
				Action:    ClientSend,
				Flags:     []cli.Flag{addrFlag},
			},
			{
				Name:      "client-receive",
				Aliases:   []string{"cr", "rc"},
				Usage:     "receiving files as a client",
				UsageText: "pft cr [options]",
				Action:    ClientReceive,
				Flags:     []cli.Flag{addrFlag, destDirFlag, mkdirFlag},
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

	conn, err := connectHost(ctx.String("port"))
	if err != nil {
		return err
	}

	defer conn.Close()
	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}

	err = checkHeaders(SND_HEADER, conn)
	if err != nil {
		return err
	}

	sendFlags := uint32(0)
	if ctx.Bool("zstd") {
		sendFlags |= ZstdComressionFlag
	}

	flags, err := exchangeFlags(sendFlags, conn)
	if err != nil {
		return err
	}

	var writer io.Writer = conn
	if flags&ZstdComressionFlag|sendFlags&ZstdComressionFlag != 0 {
		zstdConn, err := zstd.NewWriter(conn, zstd.WithEncoderLevel(zstd.SpeedFastest))
		if err != nil {
			return err
		}
		fmt.Println("Use zstd")
		defer zstdConn.Close()
		writer = zstdConn
	}

	return sendFiles(files, writer, size)
}

func HostReceive(ctx *cli.Context) error {
	conn, err := connectHost(ctx.String("port"))
	if err != nil {
		return err
	}
	defer conn.Close()

	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}

	err = checkDirExist(ctx.String("destdir"), ctx.Bool("mkdir"))
	if err != nil {
		return err
	}

	err = checkHeaders(RCV_HEADER, conn)
	if err != nil {
		return err
	}

	sendFlags := uint32(0)
	if ctx.Bool("zstd") {
		sendFlags |= ZstdComressionFlag
	}

	flags, err := exchangeFlags(sendFlags, conn)
	if err != nil {
		return err
	}

	var reader io.Reader = conn
	if flags&ZstdComressionFlag|sendFlags&ZstdComressionFlag != 0 {
		zstdConn, err := zstd.NewReader(conn)
		if err != nil {
			return err
		}
		fmt.Println("Use zstd")
		defer func() { go zstdConn.Close() }() //may be blocked
		reader = zstdConn
	}

	return getFiles(ctx.String("destdir"), reader, size)

}

func ClientSend(ctx *cli.Context) error {
	files := make([]string, ctx.Args().Len())
	for i := 0; i < ctx.Args().Len(); i++ {
		files[i] = ctx.Args().Get(i)
	}

	conn, err := connectClient(ctx.String("address"), ctx.String("port"))
	if err != nil {
		return err
	}
	defer conn.Close()

	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}

	err = checkHeaders(SND_HEADER, conn)
	if err != nil {
		return err
	}

	sendFlags := uint32(0)
	if ctx.Bool("zstd") {
		sendFlags |= ZstdComressionFlag
	}

	flags, err := exchangeFlags(sendFlags, conn)
	if err != nil {
		return err
	}

	var writer io.Writer = conn
	if flags&ZstdComressionFlag|sendFlags&ZstdComressionFlag != 0 {
		zstdConn, err := zstd.NewWriter(conn, zstd.WithEncoderLevel(zstd.SpeedFastest))
		if err != nil {
			return err
		}
		fmt.Println("Use zstd")
		defer zstdConn.Close()
		writer = zstdConn
	}

	return sendFiles(files, writer, size)
}

func ClientReceive(ctx *cli.Context) error {

	conn, err := connectClient(ctx.String("address"), ctx.String("port"))
	if err != nil {
		return err
	}
	defer conn.Close()

	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}

	err = checkDirExist(ctx.String("destdir"), ctx.Bool("mkdir"))
	if err != nil {
		return err
	}

	err = checkHeaders(RCV_HEADER, conn)
	if err != nil {
		return err
	}

	sendFlags := uint32(0)
	if ctx.Bool("zstd") {
		sendFlags |= ZstdComressionFlag
	}

	flags, err := exchangeFlags(sendFlags, conn)
	if err != nil {
		return err
	}

	var reader io.Reader = conn
	if flags&ZstdComressionFlag|sendFlags&ZstdComressionFlag != 0 {
		zstdConn, err := zstd.NewReader(conn)
		if err != nil {
			return err
		}
		fmt.Println("Use zstd")
		defer func() { go zstdConn.Close() }() //may be blocked
		reader = zstdConn
	}

	return getFiles(ctx.String("destdir"), reader, size)
}

func HostPipe(ctx *cli.Context) error {
	size, err := bufSizeToNum(ctx.String("buffer-size"))
	if err != nil {
		return err
	}

	conn1, err := connectHost(ctx.String("port"))
	if err != nil {
		return err
	}
	defer conn1.Close()

	conn2, err := connectHost(ctx.String("port"))
	if err != nil {
		return err
	}
	defer conn2.Close()


	pipe := func (in net.Conn, out net.Conn, res chan error) {
		var err error = nil

		buf := make([]byte, size)
		for {
			n, err := in.Read(buf)
			if err != nil {
				break
			}
			_, err = out.Write(buf[:n])
			if err != nil {
				break
			}
		}
		if err == io.EOF {
			err = nil
		}
		res <- err
	}

	res1 := make(chan error)
	res2 := make(chan error)
	
	go pipe(conn1, conn2, res1)
	go pipe(conn2, conn1, res2)
	
	err = <-res1
	if err != nil {
		return err
	}
	err = <-res2
	return err
}
