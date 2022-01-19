// Copyright 2022 the System Transparency Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/system-transparency/stmgr/keygen"
	"github.com/system-transparency/stmgr/ospkg"
	"github.com/system-transparency/stmgr/provision"
	"github.com/system-transparency/stmgr/sign"
)

const (
	usage = `Usage: stmgr <COMMAND> <SUBCOMMAND> [flags...]
COMMANDS:
	ospkg:
		Set of commands related to OS packages. This includes
		creating, signing and analyzing them.

	provision:
		Set of commands to provision a node for system-transparency
		usage, like creating and writing a host configuration.

	keygen:
		Commands to generate different keys and certificates for
		system-transparency.

	build:
		Not yet implemented!

Use 'stmgr <COMMAND> -help' for more info.
`

	ospkgUsage = `SUBCOMMANDS:
	create:
		Create an OS package from the provided operating
		system files.

	sign:
		Sign the provided OS package with your private key.

Use 'stmgr ospkg <SUBCOMMAND> -help' for more info.
`

	provisionUsage = `SUBCOMMANDS:
	hostconfig:
		Allows creating host configurations by spawning a TUI in
		which the user can input values into that are converted
		into a host_configuration.json file.

Use 'stmgr provision <SUBCOMMAND> -help' for more info.
`

	keygenUsage = `SUBCOMMANDS:
	certificate:
		Generate certificates for signing OS packages
		using ED25519 keys.

Use 'stmgr keygen <SUBCOMMAND> -help' for more info.
`
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// Display helptext if no arguments are given
	if len(args) < 3 {
		fmt.Print(usage)
		return nil
	}

	// Check which command is requested or display usage
	switch args[1] {
	case "ospkg":
		// Check for ospkg subcommands
		switch args[2] {
		case "create":
			// Create tool and flags
			createCmd := flag.NewFlagSet("createOSPKG", flag.ExitOnError)
			createOut := createCmd.String("out", "", "OS package output path. Two files will be created: the archive ZIP file and the descriptor JSON file. A directory or a filename can be passed. In case of a filename the file extensions will be set properly. Default name is system-transparency-os-package.")
			createLabel := createCmd.String("label", "", "Short description of the boot configuration. Defaults to 'System Transparency OS package <kernel>'.")
			createURL := createCmd.String("url", "", "URL of the OS package zip file in case of network boot mode.")
			createKernel := createCmd.String("kernel", "", "Operating system kernel.")
			createInitramfs := createCmd.String("initramfs", "", "Operating system initramfs.")
			createCmdLine := createCmd.String("cmdline", "", "Kernel command line.")

			if err := createCmd.Parse(args[3:]); err != nil {
				return err
			}
			return ospkg.Run(*createOut, *createLabel, *createURL, *createKernel, *createInitramfs, *createCmdLine)

		case "sign":
			// Sign tool and flags
			signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
			signKey := signCmd.String("key", "", "Private key for signing.")
			signCert := signCmd.String("cert", "", "Certificate corresponding to the private key.")
			signOSPKG := signCmd.String("ospkg", "", "OS package archive or descriptor file. Both need to be present.")

			if err := signCmd.Parse(args[3:]); err != nil {
				return err
			}
			return sign.Run(*signKey, *signCert, *signOSPKG)

		case "show":
			// Show tool and flags
			fmt.Println("Not implemented yet!")
			return nil

		default:
			// Display usage on unknown subcommand
			fmt.Print(ospkgUsage)
			return nil
		}

	case "provision":
		// Check for provision subcommands
		switch args[2] {
		case "hostconfig":
			// Host configuration tool and flags
			hostconfigCmd := flag.NewFlagSet("provision", flag.ExitOnError)
			hostconfigEfi := hostconfigCmd.Bool("efi", false, "Store host_configuration.json in the efivarfs.")
			hostconfigVersion := hostconfigCmd.Int("version", 1, "Hostconfig version.")
			hostconfigAddrMode := hostconfigCmd.String("addrMode", "", "Hostconfig network_mode.")
			hostconfigHostIP := hostconfigCmd.String("hostIP", "", "Hostconfig host_ip.")
			hostconfigGateway := hostconfigCmd.String("gateway", "", "Hostconfig gateway.")
			hostconfigDNS := hostconfigCmd.String("dns", "", "Hostconfig dns.")
			hostconfigInterface := hostconfigCmd.String("interface", "", "Hostconfig network_interface.")
			hostconfigURLs := hostconfigCmd.String("urls", "", "Hostconfig provisioning_urls.")
			hostconfigID := hostconfigCmd.String("id", "", "Hostconfig identity.")
			hostconfigAuth := hostconfigCmd.String("auth", "", "Hostconfig authentication.")

			if err := hostconfigCmd.Parse(args[3:]); err != nil {
				return err
			}
			return provision.Run(*hostconfigEfi, *hostconfigVersion, *hostconfigAddrMode, *hostconfigHostIP, *hostconfigGateway, *hostconfigDNS, *hostconfigInterface, *hostconfigURLs, *hostconfigID, *hostconfigAuth)

		default:
			// Display usage on unknown subcommand
			fmt.Print(provisionUsage)
			return nil
		}

	case "keygen":
		// Check for keygen subcommands
		switch args[2] {
		case "certificate":
			// Certificate tool and flags
			certificateCmd := flag.NewFlagSet("keygen", flag.ExitOnError)
			certificateRootCert := certificateCmd.String("rootCert", "", "Root certificate in PEM format to sign the new certificate. Ignored if -isCA is set.")
			certificateRootKey := certificateCmd.String("rootKey", "", "Root key in PEM format to sign the new certificate. Ignored if -isCA is set.")
			certificateIsCA := certificateCmd.Bool("isCA", false, "Generate self signed root certificate.")
			certificateValidFrom := certificateCmd.String("validFrom", "", "Date formatted as RFC822. Defaults to time of creation.")
			certificateValidUntil := certificateCmd.String("validUntil", "", "Date formatted as RFC822. Defaults to time of creation + 72h.")
			certificateCertOut := certificateCmd.String("certOut", "", "Output certificate file. Defaults to cert.pem or rootcert.pem is -isCA is set.")
			certificateKeyOut := certificateCmd.String("keyOut", "", "Output key file. Defaults to key.pem or rootkey.pem if -isCA is set.")

			if err := certificateCmd.Parse(args[3:]); err != nil {
				return err
			}
			return keygen.Run(*certificateIsCA, *certificateRootCert, *certificateRootKey, *certificateValidFrom, *certificateValidUntil, *certificateCertOut, *certificateKeyOut)

		default:
			// Display usage on unknown subcommand
			fmt.Print(keygenUsage)
			return nil
		}

	case "build":
		// Check for build subcommands
		switch args[2] {
		default:
			// Display usage on unknown subcommand
			fmt.Println("Not implemented yet!")
			return nil
		}

	default:
		// Display usage on unknown command
		fmt.Print(usage)
		return nil
	}
}
