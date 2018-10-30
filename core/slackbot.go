package core

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"strings"
)

const activeRecastToken = secondBotToken

type SlackBot struct {
	slackToken     string
	client         *slack.Client
	rtm            *slack.RTM
	creator        *MessageManager
	conversationID string
}

func NewSlackBot() (*SlackBot, error) {

	token := "xoxb-438453325860-438070557617-CviJFdimezMGe8FM04MwfO5a"
	client := slack.New(token)
	rtm := client.NewRTM()
	creator, err := NewMessageCreator(activeRecastToken)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create MessageManager")
	}

	return &SlackBot{slackToken: token, client: client, rtm: rtm, creator: creator, conversationID: "1"}, nil
}

func (bot *SlackBot) Run() {

	go bot.rtm.ManageConnection()
	for {
		select {
		case message := <-bot.rtm.IncomingEvents:
			fmt.Print("Event Received: ")

			switch event := message.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", event.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", event.Text)

				bot.Respond(event)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", event.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break
			}
		}
	}

}

func (bot *SlackBot) Respond(msg *slack.MessageEvent) {

	response := ""
	text := msg.Text

	if strings.ToLower(text) == "%new" {

		bot.conversationID = bot.getNewConversationID()
		response = "new conversation with ID:" + bot.conversationID

		bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(response, msg.Channel))

		newCreator, err := NewMessageCreator(activeRecastToken)
		if err != nil {
			fmt.Println(err)
		}

		bot.creator = newCreator
		return

	} else if strings.Contains(strings.ToLower(text), "%switch") {

		bot.conversationID = bot.getConversationID(text)
		response = "switch to conversation with ID:" + bot.conversationID

		bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(response, msg.Channel))
		return
	}

	response, err := bot.creator.Response(text, bot.conversationID)
	if err != nil {
		fmt.Println(err)
	}

	bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(response, msg.Channel))
}

func (bot *SlackBot) getNewConversationID() string {
	newID := bot.creator.NewConversationID()
	return newID
}

func (bot *SlackBot) getConversationID(text string) string {
	values := strings.Split(text, " ")
	convID := strings.TrimSpace(values[1])
	return convID
}
