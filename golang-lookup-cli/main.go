package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	//.Run(os.Args)
	app.Name = "Website Lookup CLI"
	app.Usage = "Let's you query IPs, CNAMEs, MX records and Name Servers!"
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// We'll be using the same flag for all our commands
	// so we'll define it up here
	hostFlag := cli.StringFlag{
		Name:  "host",
		Value: "google.com",
	}
	hostFlags := []cli.Flag{
		hostFlag,
	}
	// we create our commands
	app.Commands = []cli.Command{
		{
			Name:  "ns",
			Usage: "Looks Up the NameServers for a Particular Host",
			Flags: hostFlags,
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				// a simple lookup function
				ns, err := net.LookupNS(c.String(hostFlag.Name))
				if err != nil {
					return err
				}
				// we log the results to our console
				// using a trusty fmt.Println statement
				for i := 0; i < len(ns); i++ {
					fmt.Println(ns[i].Host)
				}
				return nil
			},
		},
		{
			Name:  "ip",
			Usage: "Looks up the IP addresses for a particular host",
			Flags: hostFlags,
			Action: func(c *cli.Context) error {
				ip, err := net.LookupIP(c.String(hostFlag.Name))
				if err != nil {
					fmt.Println(err)
				}
				for i := 0; i < len(ip); i++ {
					fmt.Println(ip[i])
				}
				return nil
			},
		},
		{
			Name:  "cname",
			Usage: "Looks up the CNAME for a particular host",
			Flags: hostFlags,
			Action: func(c *cli.Context) error {
				cname, err := net.LookupCNAME(c.String(hostFlag.Name))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(cname)
				return nil
			},
		},
		{
			Name:  "mx",
			Usage: "Looks up the MX records for a particular host",
			Flags: hostFlags,
			Action: func(c *cli.Context) error {
				mx, err := net.LookupMX(c.String(hostFlag.Name))
				if err != nil {
					fmt.Println(err)
				}
				for i := 0; i < len(mx); i++ {
					fmt.Println(mx[i].Host, mx[i].Pref)
				}
				return nil
			},
		},
	}

	// start our application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
