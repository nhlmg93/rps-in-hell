package main

// TODO:
//  Drowdown with RPS options
//  store RPS game away with callengers option
//  display accept button on command being sent
//  on accept display same challenge options dropdown
//    if no one accepts after an hour remove stored challenge and message
//  calculate winner, remove from stored games, display winner in message
import (
	"flag"
	"fmt"
	"go/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)


var(
    guildID string
    appID string
    token string
    dg *discordgo.Session
    cmdUpd = flag.Bool("cmdupd", false, "deploy a cmd")
)


func init() {
    flag.Parse()

    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

    token = os.Getenv("TOKEN")
    guildID = os.Getenv("GUILD_ID")
    appID = os.Getenv("APP_ID")

    // Create a new Discord session using the provided bot token.
    dg, err = discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }
}

func main() {
    if *cmdUpd {
        cmds, _ := dg.ApplicationCommands(appID,"")
        for _, cmd := range cmds {
            dg.ApplicationCommandDelete(appID, "", cmd.ID)
        } 
   
        cmds, _ = dg.ApplicationCommands(appID,"")
        fmt.Printf("\ncmd: %v\n", cmds)
    }
    _, err := dg.ApplicationCommandBulkOverwrite(appID, guildID,
            []*discordgo.ApplicationCommand{
                {
                    Name: "challenge",
                    Description: "Start a hellish game of RPS",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Name: "object",
                            Description: "Play your move with a fiery passion",
                            Type: discordgo.ApplicationCommandOptionString,
                            Required: true,
                            Choices: []*discordgo.ApplicationCommandOptionChoice{
                                {
                                    Name: "Rock",
                                    Value: 1,
                                },
                                {
                                    Name: "Paper",
                                    Value: 2,
                                },
                                {
                                    Name: "Scissors",
                                    Value: 3,
                                },
                            },
                        },
                    },
                },
            },
        )
    if err != nil {
        fmt.Println(err) 
        return
    }
    dg.AddHandler(func (s *discordgo.Session, i *discordgo.InteractionCreate) {
        data := i.ApplicationCommandData()
        switch data.Name {
        case "challenge":
            var err error
            //err := s.InteractionRespond(
            //    i.Interaction,
            //    &discordgo.InteractionResponse {
            //        Type: discordgo.InteractionResponseChannelMessageWithSource,
            //        Data: &discordgo.InteractionResponseData {
            //            Content: "Hello world!",
            //        },
            //    })
            if err != nil {
                fmt.Println(err)
                return
            }
        }
    })
    
    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close() 
}

