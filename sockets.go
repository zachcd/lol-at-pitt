package main

import (
	"log"
	"sync"

	"fmt"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/lab-d8/lol-at-pitt/draft"
	"github.com/lab-d8/oauth2"
	"github.com/martini-contrib/render"
)

// Message is a single message sent between clients
type Message struct {
	Type string `json:"type"`
	From string `json:"from"`
	Text string `json:"text"`
}

// DraftRoom used for socket communication DraftRoom
type DraftRoom struct {
	sync.Mutex
	clients []*Client
}

// Client is a single connection to be used for websockets
type Client struct {
	ID         string
	in         <-chan *Message
	out        chan<- *Message
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
}

// Add a client to a room
func (r *DraftRoom) appendClient(client *Client) {
	r.Lock()
	r.clients = append(r.clients, client)
	for _, c := range r.clients {
		if c != client {
			c.out <- &Message{"status", client.ID, "Joined this chat"}
		}
	}
	r.Unlock()
}

// Remove a client from a room
func (r *DraftRoom) removeClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	for index, c := range r.clients {
		if c == client {
			r.clients = append(r.clients[:index], r.clients[(index+1):]...)
		} else {
			c.out <- &Message{"status", client.ID, "Left this chat"}
		}
	}
}

// Message all the other clients in the same room
func (r *DraftRoom) messageOtherClients(client *Client, msg *Message) {
	r.Lock()
	msg.From = client.ID

	for _, c := range r.clients {
		if c != client {
			c.out <- msg
		}
	}
	defer r.Unlock()
}

func (r *DraftRoom) broadcast(msg *Message) {
	r.Lock()
	for _, c := range r.clients {
		c.out <- msg
	}
	defer r.Unlock()
}

func (r *DraftRoom) messageWithID(id string, msg *Message) {
	for _, c := range r.clients {
		if c.ID == id {
			r.message(c, msg)
		}
	}
}

func (r *DraftRoom) message(client *Client, msg *Message) {
	client.out <- msg
}
func newDraftRoom() *DraftRoom {
	return &DraftRoom{sync.Mutex{}, make([]*Client, 0)}
}

var room *DraftRoom

// SocketRouter is used to setup main martini function
func SocketRouter(m *martini.ClassicMartini) {
	room = newDraftRoom()
	Init()
	m.Get("/draft", LoginRequired, func(r render.Render, token oauth2.Tokens) {
		id, _ := GetId(token.Access())
		r.HTML(200, "draft", id)
	})

	m.Get("/admin/start", func() {
		draft.Paused = false
		allowTicks = false
		fmt.Println("Hello!")
		Handle(Message{Type: "event", Text: "The round has started, LET THE BIDDING BEGIN"})
	})

	m.Get("/admin/reset", func() {
		draft.GetCurrentPlayer().HighestBid = 0
		draft.GetCurrentPlayer().Team = ""
		Handle(Message{Type: "event", Text: "Admin reset current round, starting when they press the button.."})
		allowTicks = false
		draft.Paused = true
	})

	m.Get("/admin/skip", func() {
		Handle(Message{Type: "event", Text: "Admin skipped current player, waiting on him to start.."})
		Handle(Message{Type: "update"})

		draft.Next()
	})

	m.Get("/admin/previous", func() {
		Handle(Message{Type: "event", Text: "Admin undid previous round, waiting on him to start.."})
		Handle(Message{Type: "update"})
		draft.Previous()
	})
	// This is the sockets connection for the room, it is a json mapping to sockets.
	m.Get("/draft/:clientname", sockets.JSON(Message{}), func(params martini.Params, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
		client := &Client{params["clientname"], receiver, sender, done, err, disconnect}
		room.appendClient(client)
		Handle(Message{Type: "login"})
		// A single select can be used to do all the messaging
		for {
			select {
			case <-client.err:
				// Its gone jim
			case msg := <-client.in:
				// Most of code will be handled here
				log.Println(*msg)
				Handle(*msg)
			case <-client.done:
				room.removeClient(client)
				return 200, "OK"
			}
		}
	})

}
