package main

import (
    "fmt"
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
)

func init() {
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
    
    _, err := dg.ApplicationCommandBulkOverwrite(appID, guildID,
            []*discordgo.ApplicationCommand{
                {
                    Name: "hello-world",
                    Description: "Showcase of a basic slash command",
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
        case "hello-world":
            err := s.InteractionRespond(
                i.Interaction,
                &discordgo.InteractionResponse {
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData {
                        Content: "Hello world!",
                    },
                })
            if err != nil {
                fmt.Println(err)
                return
            }
        }
    })
    
    dg.Identify.Intents = discordgo.IntentsGuildMessages
    
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

