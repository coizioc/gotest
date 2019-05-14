package main

import (
	"deathmatch"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

var cmdHandler *CmdHandler

// Creds struct to read in the creds yaml file
type Creds struct {
	AppToken          string
	AppID             string
	BotToken          string
	BotName           string
	BotDiscriminator  string
	BotPermissionsInt int
	BotPrefix         string
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	configFilePath := "config.yaml"
	creds := Creds{}
	err := fillCreds(configFilePath, &creds)
	if err != nil {
		log.Fatal(err)
	}

	dg, err := discordgo.New("Bot " + creds.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	initCmds(dg, creds.BotPrefix)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func fillCreds(configFilePath string, creds *Creds) (err error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &creds)
	if err != nil {
		return err
	}

	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmdHandler.Handle(m)
}

func initCmds(s *discordgo.Session, prefix string) {
	cmdHandler = NewCmdHandler(s, prefix)
	cmdHandler.NewCmd("snap", cmdSnap)

	cmdHandler.NewCmd("dm", cmdDeathmatch)
	cmdHandler.NewCmd("deathmatch", cmdDeathmatch)
}

func cmdSnap(ctx *Ctx, args []string) error {
	charTotal := 0
	for _, str := range args[1:] {
		for _, c := range str {
			charTotal += int(c)
		}
	}

	outStr := ""
	if charTotal%2 == 0 {
		outStr = fmt.Sprintf("%s, you have been slain by Thanos, for the good of the Universe.", strings.Join(args[1:], " "))
	} else {
		outStr = fmt.Sprintf("%s, you were spared by Thanos.", strings.Join(args[1:], " "))
	}

	ctx.Send(outStr)
	return nil
}

func cmdDeathmatch(ctx *Ctx, args []string) error {
	members := ctx.Guild.Members

	var m1, m2 *discordgo.Member
	if len(args) == 1 {
		m1, m2 = members[rand.Intn(len(members))], members[rand.Intn(len(members))]
	} else if len(args) == 2 {
		for _, member := range members {
			if strings.Contains(member.User.Username, args[1]) {
				m1 = member
			}
		}
		m2 = members[rand.Intn(len(members))]
	} else {
		for _, member := range members {
			if strings.Contains(member.User.Username, args[1]) {
				m1 = member
			}
			if strings.Contains(member.User.Username, args[2]) {
				m2 = member
			}
		}
	}

	if m1 == nil {
		msg := fmt.Sprintf("Name %s not found.", args[1])
		ctx.Send(msg)
		return nil
	}
	if m2 == nil {
		msg := fmt.Sprintf("Name %s not found.", args[2])
		ctx.Send(msg)
		return nil
	}

	dmMessages := deathmatch.Deathmatch(m1, m2)
	msg, _ := ctx.Send(dmMessages[0])
	for _, dmMsg := range dmMessages[1:] {
		ctx.Edit(msg.ID, dmMsg)
		time.Sleep(2 * time.Second)
	}
	return nil
}
