package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type argFlags struct {
	Addr                   string
	User                   string
	Password               string
	Database               string
	OutputFolderPath       string
	Recover                bool
	Seed                   bool
	SshServerUser          string
	SshServerPassword      string
	SshServerHost          string
	SshServerPort          int
	SshServerDirectoryPath string
}

var sshConnection Connection

var currentFlags argFlags

func main() {

	err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if currentFlags.Seed == true {
		err := seed(db)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	if currentFlags.Recover == true {
		err := recover(db)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	err = backup(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func parseArgs() error {
	var opts struct {
		Addr                   string `short:"a" long:"addres" required:"true"`
		Password               string `short:"p" long:"password" required:"true"`
		User                   string `short:"u" long:"username" required:"true"`
		Database               string `short:"d" long:"database" required:"true"`
		OutputFolderPath       string `short:"o" long:"output_folder_path" description:"The location of the folder to be used for exit data"`
		Recover                bool   `short:"r" long:"recover" description:"Gets the data from the specified backup folder and inserts it in the database"`
		Seed                   bool   `short:"s" long:"seed" description:"Gets the data from <current_directory>/seed and inserts it in the database"`
		SshServerHost          string `short:"h" long:"sshhost"`
		SshServerUser          string `short:"l"  long:"sshuser"`
		SshServerPassword      string `short:"m"  long:"sshpwd"`
		SshServerPort          int    `short:"n"  long:"sshport" default:"-1"`
		SshServerDirectoryPath string `short:"k"  long:"sshdirectory"`
	}

	parser := flags.NewParser(&opts, flags.None)

	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		return err
	}

	if opts.Recover == true && opts.OutputFolderPath == "" {
		return fmt.Errorf("You cannot use recover without specifying a folder -o <folder_name> containing data.")
	}

	if opts.Recover == false && opts.Seed == false {
		switch {
		case opts.OutputFolderPath == "":
			return fmt.Errorf("You need to specify folder in which you want to backup your data.")
		case opts.SshServerDirectoryPath == "":
			return fmt.Errorf("You need to specify folder in which you want to backup your data on the remote server.")
		case opts.SshServerHost == "":
			return fmt.Errorf("You need to specify SSH server host.")
		case opts.SshServerPort == -1:
			return fmt.Errorf("You need to specify SSH server port.")
		case opts.SshServerUser == "":
			return fmt.Errorf("You need to specify SSH user name.")
		case opts.SshServerPassword == "":
			return fmt.Errorf("You need to specify SSH user password.")

		}

	}

	currentFlags = argFlags{
		Addr:             opts.Addr,
		Password:         opts.Password,
		User:             opts.User,
		Database:         opts.Database,
		OutputFolderPath: opts.OutputFolderPath,
		Recover:          opts.Recover,
		Seed:             opts.Seed,
	}

	sshConnection = Connection{
		host:                opts.SshServerHost,
		port:                opts.SshServerPort,
		user:                opts.SshServerUser,
		password:            opts.SshServerPassword,
		backupDirectoryPath: opts.SshServerDirectoryPath,
	}

	return nil
}
