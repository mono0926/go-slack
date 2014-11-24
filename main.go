package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mono0926/go-backlog/api"
	qiita "github.com/mono0926/go-qiita/api"
	"github.com/mono0926/go-slack/web_hook"
	"html"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
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
		{
			Name:      "qiita",
			ShortName: "q",
			Usage:     `specify "teamKey", "apiKey", and "webHookUrl"`,
			Action: func(c *cli.Context) {
				log.Println("starting qiita...")
				args := c.Args()
				if len(args) < 3 {
					log.Println("not enough arguments.", args)
					return
				}
				teamKey := args[0]
				apiKey := args[1]
				webHookUrl := args[2]
				postNewQiita(teamKey, apiKey, webHookUrl)
			},
		},
	}
	app.Run(os.Args)

}

func postBacklogNotifications(teamKey string, apiKey string, webHookUrl string, toUsername string) {
	a := api.NewApi(teamKey, apiKey)
	client := a.NotificationClient()
	for _, n := range client.GetNotifications(false) {
		f := `%s さん「%s」\n「%s」\nhttps://%s.backlog.jp/view/%s`
		comment := n.Comment.Content
		comments := strings.Split(n.Comment.Content, "\n")
		if len(comments) > 0 {
			comment = comments[0]
		}
		m := fmt.Sprintf(f, n.Comment.CreatedUser.Name, comment, n.Issue.Summary, teamKey, n.Issue.IssueKey)
		m = html.EscapeString(m)
		log.Println(m)
		web_hook.Post(webHookUrl, toUsername, "backlog", m, ":backlog:")
	}
}

func postNewQiita(teamKey string, apiKey string, webHookUrl string) {
	a := qiita.NewTeamApi(teamKey, apiKey)
	client := a.ItemClient()
	items := client.GetItems()
	if len(items) == 0 {
		return
	}
	filename := "qiita.log"
	latest := read(filename)
	log.Println(latest)
	latestTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", latest)
	log.Println(latestTime)
	for _, n := range client.GetItems() {
		if n.CreatedAt.Unix() <= latestTime.Unix() {
			continue
		}
		f := "%sさんが「%s」を投稿しました。 %s"
		m := fmt.Sprintf(f, n.User.Name, n.Title, n.Url(teamKey+"."))
		m = html.EscapeString(m)
		log.Println(m)
		web_hook.Post(webHookUrl, "#qiita", "qiita", m, ":qiita:")
	}
	write(filename, fmt.Sprintf("%s", items[0].CreatedAt))
}

func read(filename string) string {
	bs, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Fatalln(e)
	}
	return string(bs)
}

func write(filename, content string) {
	fout, e := os.Create(filename)
	if e != nil {
		log.Fatalln(e)
	}
	defer fout.Close()
	fout.WriteString(content)
}
