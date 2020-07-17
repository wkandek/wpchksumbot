package main


import (
    "crypto/sha1"
    "fmt"
    "io"
    "io/ioutil"
	"irc"
    "log"
	"net"
    "net/http"
    "os"
)


const server = "chat.freenode.net:6667"
const channel = "#demok8sws"
var urls = [...]string{"https://de.wikipedia.org", "https://el.wikipedia.org", "https://en.wikipedia.org"}


func send_irc(msgs []string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatalln(err)
	}

	config := irc.ClientConfig{
		Nick: "WPChkSumBot",
        Pass: "password",
        User: "chkbot",
        Name: "Checksum Bot",
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
			if m.Command == "001" {
				// 001 is a welcome event, so we join channels there
				c.Writef("JOIN %s", channel)
			} else if m.Command == "366" {
                for _, l := range msgs {
				    c.Writef("PRIVMSG %s %s", channel, l)
                }
				c.Write("QUIT ")
                os.Exit(0)
			}
		}),
	}

	// Create the client
	client := irc.NewClient(conn, config)
	err = client.Run()
	if err != nil {
		log.Fatalln(err)
	}
}


func main() {
    msgs := []string{}

    for _, url := range urls {
        response, err := http.Get(url)
        if err != nil {
            log.Fatal(err)
        }
        defer response.Body.Close()
        bodybytes, err := ioutil.ReadAll(response.Body)

        h :=  sha1.New()
        h.Write(bodybytes)
        bs := h.Sum(nil)

        _, err = io.Copy(os.Stdout, response.Body)
        if err != nil {
            log.Fatal(err)
        }
        msgs = append(msgs, fmt.Sprintf("%s:%x\n", url, bs))
    }
    send_irc(msgs)
}
