package main

import (
	"bytes"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	sqlClient *SqlClient
	graphqlClient *GraphqlClient
)

func init() {
	// Load config
	snapConfig := Config{}
	err := snapConfig.InitConfig()
	if err != nil{
		log.Fatal(err)
	}
	// Init sql client
	sqlClient = &SqlClient{}
	err = sqlClient.Init()
	if err != nil{
		log.Fatal(err)
	}

	// Init graphql client
	graphqlClient = &GraphqlClient{}
	graphqlClient.Init(viper.GetString("graphql.url"))
}

func main() {
	app := iris.Default()

	app.UseRouter(Cors)

	ac := makeAccessLog()
	defer ac.Close()

	// Register the middleware (UseRouter to catch http errors too).
	app.UseRouter(ac.Handler)

	booksAPI := app.Party("/snapshot")
	{
		booksAPI.Use(iris.Compression)
		booksAPI.Post("/webhooks", SnapShotHooks)
		booksAPI.Post("/subscribe", subscribe)
	}

	app.Listen(":8009")
}

// CORS middleware
func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Method() == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}

func makeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	ac := accesslog.File("./access.log")
	// Remove this line to disable logging to console:
	ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})
	return ac
}

func subscribe(ctx iris.Context) {
	var s SubscribeModel
	err := ctx.ReadJSON(&s)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Key("message", "Subscribe failure").DetailErr(err))
		return
	}
	// Verify email format
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //match email address
	reg := regexp.MustCompile(pattern)
	if ok := reg.MatchString(s.Email);!ok{
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Key("message", "Subscribe failure, please check the eamil address format."))
		return
	}

	// Add a new subscribe
	err = sqlClient.NewSubscribe(s)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Key("message", "Subscribe failure").DetailErr(err))
		return
	}
	ctx.JSON(iris.Map{
		"status":  iris.StatusOK,
		"message": "Subscribe success",
	})
}

func SendMail(mailTo []string, subject string, body string) error {

	mailConn := viper.GetStringMapString("mailConn")

	port, _ := strconv.Atoi(mailConn["port"]) // convert string to int

	m := gomail.NewMessage()

	m.SetHeader("From",  m.FormatAddress(mailConn["user"], mailConn["name"]))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}

type Event struct {
	ID string `json:"id"`
	Event string `json:"event"`
	Space string `json:"space"`
	Expire int `json:"expire"`
}

func SnapShotHooks(ctx iris.Context) {
	var e Event
	err := ctx.ReadJSON(&e)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().DetailErr(err))
		log.Println(err)
		return
	}
	// Judge whether it's a proposal event
	if strings.HasPrefix(e.Event, "proposal/"){
		propRes 	:= graphqlClient.GetSingleProposal(strings.ReplaceAll(e.ID, "proposal/", ""))
		prop 		:= propRes.Proposal
		render 		:= PropRender{}
		timeFormat  := "2006-01-02"
		render.Start = time.Unix(prop.Start, 0).Format(timeFormat)
		render.End   = time.Unix(prop.End, 0).Format(timeFormat)
		render.SpaceName = prop.Space.Name
		render.Title = prop.Title
		if len(prop.Body) > 100{
			render.Body  = prop.Body[0:150] + "..."
		} else {
			render.Body  = prop.Body
		}
		render.State = prop.State
		render.Type  = prop.Type
		render.Link  = prop.Link

		switch e.Event {
			case "proposal/created":
				return
			case "proposal/start":
				render.Subject = `<h2 style="line-height: 36px; font-size: 1.5em; font-weight: bold;">THE PROPOSAL YOU</h2><h2 style="line-height: 36px; font-size: 1.5em; font-weight: bold;">SUBSCRIBED BEGAN</h2>`
				break
			case "proposal/end":
				render.Subject = `<h2 style="line-height: 36px; font-size: 1.5em; font-weight: bold;">THE PROPOSAL YOU</h2><h2 style="line-height: 36px; font-size: 1.5em; font-weight: bold;">SUBSCRIBED ENDED</h2>`
				break
			case "proposal/deleted":
				render.Subject = `<h2 style="line-height: 36px; font-size: 1.5em; font-weight: bold;">THE PROPOSAL YOU</h2><h2 style="line-height: 36px; font-size: 1.5em; font-weight: bold;">SUBSCRIBED WAS DELETED</h2>`
				break
			default:
				return
		}

		// Mail recipients
		mailTo, err := sqlClient.GetSubscribeEmail(strings.ReplaceAll(e.ID, "proposal/", ""))
		if err != nil{
			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().DetailErr(err))
			return
		}
		// Load mail template
		body, err := RenderMailTmpl(render)
		if err != nil{
			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().DetailErr(err))
			return
		}
		subject := fmt.Sprintf("[Snapshot/%s]%s(#%s)", render.SpaceName, render.Title, render.State)
		SendMail(mailTo, subject, body)
	}
	ctx.JSON(iris.Map{
		"status":  iris.StatusOK,
		"message": "success",
	})
}

type PropRender struct{
	Subject string
	Title string
	Start string
	End string
	SpaceName string
	Body string
	State string
	Type string
	Link string
}

func RenderMailTmpl(prop PropRender) (string, error) {
	tmpl := *template.New("email.tmpl")
	tpl, err := tmpl.ParseFiles("./email.tmpl")
	if err != nil{
		return "1", err
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, prop)
	if err != nil {
		return "2", err
	}
	return buf.String(), nil
}
