package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/prince/internal/cmd"
)

func main() {
	app := &cli.App{
		Name:  "prince",
		Usage: "代理服务",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "启动服务端",
			},
			&cli.BoolFlag{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "启动客户端端",
			},
			&cli.StringFlag{
				Name:  "proxy_host",
				Usage: "代理服务监听地址",
			},
			&cli.StringFlag{
				Name:  "transfer_host",
				Usage: "转发服务监听地址",
			},
		},
		Action: func(c *cli.Context) error {
			isServer := c.Bool("server")
			isClient := c.Bool("client")
			if isServer && isClient {
				return fmt.Errorf("不能同时启动服务端和客户端")
			} else if !isServer && !isClient {
				return fmt.Errorf("请指定 --server 启动服务端，或者指定 --client 启动客户端")
			}

			proxyHost := c.String("proxy_host")
			transferHost := c.String("transfer_host")

			if isClient {
				return cmd.RunClient(transferHost)
			}

			if isServer {
				return cmd.RunServer(transferHost, proxyHost)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
