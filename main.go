package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mono0926/backlog-slack/web_hook"
	"github.com/mono0926/go-backlog/api"
	"html"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "backlog-slack"
	app.Usage = "TODO"
	app.Action = func(c *cli.Context) {

	}
	app.Commands = []cli.Command{
		{
			Name:      "notifications",
			ShortName: "n",
			Usage:     `specify "teamKey", "apiKey", "webHookUrl", and "toUsername"`,
			Action: func(c *cli.Context) {
				log.Println("starting notifications...")
				args := c.Args()
				if len(args) < 4 {
					log.Println("not enough arguments.", args)
					return
				}
				teamKey := args[0]
				apiKey := args[1]
				webHookUrl := args[2]
				toUsername := args[3]
				postBacklogNotifications(teamKey, apiKey, webHookUrl, toUsername)
			},
		},
	}
	app.Run(os.Args)

}

func postBacklogNotifications(teamKey string, apiKey string, webHookUrl string, toUsername string) {
	a := api.NewApi(teamKey, apiKey)
	client := a.NotificationClient()
	for _, n := range client.GetNotifications(true) {
		f := `%s さん「%s」\n「%s」\nhttps://%s.backlog.jp/view/%s`
		comment := n.Comment.Content
		comments := strings.Split(n.Comment.Content, "\n")
		if len(comments) > 0 {
			comment = comments[0]
		}
		m := fmt.Sprintf(f, n.Comment.CreatedUser.Name, comment, n.Issue.Summary, teamKey, n.Issue.IssueKey)
		m = html.EscapeString(m)
		log.Println(m)
		web_hook.Post(webHookUrl, toUsername, "backlog", m, ":+1:")
	}
}
