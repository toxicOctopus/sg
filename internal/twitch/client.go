package twitch

import (
	"bufio"
	"github.com/pkg/errors"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

const (
	ircAddress   = "irc.chat.twitch.tv:6667"
	readInterval = time.Millisecond * 10
	ircNL        = "\r\n"

	maxAuthVerificationTries = 10
	authVerificationInterval = time.Millisecond * 100
)

// Regex for parsing PRIVMSG strings.
//
// First matched group is the user's name and the second matched group is the content of the
// user's message.
var MsgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)

type Client struct {
	connection net.Conn
}

//GetClient establishes connection including authorization
func GetClient(nick, pass string) (Client, error) {
	client := Client{}

	conn, err := net.Dial("tcp", ircAddress)
	if err != nil {
		return client, errors.Wrap(err, "Cannot connect to "+ircAddress)
	}
	client.connection = conn

	_, err = client.connection.Write([]byte("PASS " + pass + ircNL))
	if err != nil {
		return client, errors.Wrap(err, "auth pass failed")
	}
	_, err = client.connection.Write([]byte("NICK " + nick + ircNL))
	if err != nil {
		return client, errors.Wrap(err, "auth nick failed")
	}

	tp := textproto.NewReader(bufio.NewReader(client.connection))
	authSuccessful := false
	for i := 0; i < maxAuthVerificationTries; i++ {
		line, err := tp.ReadLine()
		if "PING :tmi.twitch.tv" == line {
			// respond to PING message with a PONG message, to maintain the connection
			_, err = client.connection.Write([]byte("PONG :tmi.twitch.tv" + ircNL))
		} else if nil == err {
			matches := MsgRegex.FindStringSubmatch(line)
			if nil == matches {
				if strings.Contains(line, "Welcome, GLHF!") {
					authSuccessful = true
					break
				}
			}
		}
		time.Sleep(authVerificationInterval)
	}
	if !authSuccessful {
		return client, errors.New("auth failed")
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}

func (c *Client) Listen(channel string, messageCallback func(from, message string), errorCallback func(err error), warnCallback func(err error)) {
	_, err := c.connection.Write([]byte("JOIN #" + channel + ircNL))
	if err != nil {
		warnCallback(errors.Wrap(err, "failed to join channel"))
	}

	// reads from connection
	tp := textproto.NewReader(bufio.NewReader(c.connection))

	// listens for chat messages
	for {
		line, err := tp.ReadLine()

		if nil != err {
			err = c.connection.Close()
			if err != nil {
				warnCallback(errors.Wrap(err, "failed to close connection to twitch"))
			}
			errorCallback(errors.Wrap(err, "connection closed unexpectedly"))
			return
		}

		if "PING :tmi.twitch.tv" == line {
			// respond to PING message with a PONG message, to maintain the connection
			_, err = c.connection.Write([]byte("PONG :tmi.twitch.tv" + ircNL))
			if err != nil {
				warnCallback(errors.Wrap(err, "failed to pong"))
			}
			continue
		} else {
			// handle a PRIVMSG message
			matches := MsgRegex.FindStringSubmatch(line)
			if nil != matches {
				userName := matches[1]
				msgType := matches[2]

				switch msgType {
				case "PRIVMSG":
					msg := matches[3]
					messageCallback(userName, msg)
				default:
					// do nothing
				}
			}
		}
		time.Sleep(readInterval)
	}
}