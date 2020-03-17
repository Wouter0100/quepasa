package models

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	wa "github.com/Rhymen/go-whatsapp"
	qrcode "github.com/skip2/go-qrcode"
)

type WhatsAppServer struct {
	connections map[string]*wa.Conn
	handlers    map[string]*messageHandler
}

var server *WhatsAppServer

type messageHandler struct {
	botID       string
	userIDs     map[string]bool
	messages    map[string]Message
	synchronous bool
}

//
// Start
//
func StartServer() error {
	log.Println("Starting WhatsApp server")

	connections := make(map[string]*wa.Conn)
	handlers := make(map[string]*messageHandler)
	server = &WhatsAppServer{connections, handlers}

	return startHandlers()
}

func startHandlers() error {
	bots, err := FindAllBots(GetDB())
	if err != nil {
		return err
	}

	for _, bot := range bots {
		log.Printf("Adding message handlers for %s", bot.Number)

		err = startHandler(bot.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func  startHandler(botID string) error {
	con, err := createConnection()
	if err != nil {
		return err
	}

	server.connections[botID] = con

	userIDs := make(map[string]bool)
	messages := make(map[string]Message)
	startupHandler := &messageHandler{botID, userIDs, messages, true}
	con.AddHandler(startupHandler)

	session, err := readSession(botID)
	if err != nil {
		return err
	}

	session, err = con.RestoreWithSession(session)
	if err != nil {
		return err
	}

	<-time.After(3 * time.Second)

	if err := writeSession(botID, session); err != nil {
		return err
	}

	con.RemoveHandlers()

	log.Println("Fetching initial messages")

	initialMessages, err := fetchMessages(con, botID, startupHandler.userIDs)
	if err != nil {
		return err
	}

	log.Println("Setting up long-running message handler")

	asyncMessageHandler := &messageHandler{botID, startupHandler.userIDs, initialMessages, false}
	server.handlers[botID] = asyncMessageHandler
	con.AddHandler(asyncMessageHandler)

	return nil
}

func getConnection(botID string) (*wa.Conn, error) {
	connection, ok := server.connections[botID]
	if !ok {
		return nil, fmt.Errorf("Connection not found for botID %s", botID)
	}

	return connection, nil
}

func createConnection() (*wa.Conn, error) {
	con, err := wa.NewConn(30 * time.Second)
	if err != nil {
		return con, err
	}

	con.SetClientName("QuePasa for Link", "QuePasa")

	return con, err
}

func writeSession(botID string, session wa.Session) error {
	_, err := GetOrCreateStore(GetDB(), botID)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err = encoder.Encode(session); err != nil {
		return err
	}

	_, err = UpdateStore(GetDB(), botID, buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func readSession(botID string) (wa.Session, error) {
	var session wa.Session
	store, err := GetStore(GetDB(), botID)
	if err != nil {
		return session, err
	}

	r := bytes.NewReader(store.Data)
	decoder := gob.NewDecoder(r)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func SignIn(botID string, out chan<- []byte) error {
	con, err := createConnection()
	if err != nil {
		return err
	}

	qr := make(chan string)
	go func() {
		var png []byte
		png, err := qrcode.Encode(<-qr, qrcode.Medium, 256)
		if err != nil {
			log.Println(err)
		}
		encodedPNG := base64.StdEncoding.EncodeToString(png)
		out <- []byte(encodedPNG)
	}()

	session, err := con.Login(qr)
	if err != nil {
		return err
	}

	return writeSession(botID, session)
}

func SendMessage(botID string, recipient string, message string) (string, error) {
	log.Print(botID, recipient, message)
	var messageID string
	con, err := getConnection(botID)
	if err != nil {
		return messageID, err
	}

	textMessage := wa.TextMessage{
		Info: wa.MessageInfo{
			RemoteJid: recipient,
		},
		Text: message,
	}

	messageID, err = con.Send(textMessage)
	if err != nil {
		return messageID, err
	}

	return messageID, nil
}

//
// ReceiveMessages
//

func ReceiveMessages(botID string, timestamp string) ([]Message, error) {
	var messages []Message
	searchTimestamp, err := strconv.ParseUint(timestamp, 10, 64)
	if err != nil {
		searchTimestamp = 1000000
	}

	handler, ok := server.handlers[botID]
	if !ok {
		return messages, nil
	}

	for _, msg := range handler.messages {
		if msg.Timestamp >= searchTimestamp {
			messages = append(messages, msg)
		}
	}

	sort.Sort(ByTimestamp(messages))

	return messages, nil
}

func loadMessages(con *wa.Conn, botID string, userID string, count int) (map[string]Message, error) {
	userIDs := make(map[string]bool)
	messages := make(map[string]Message)
	handler := &messageHandler{botID, userIDs, messages, true}
	con.LoadFullChatHistory(userID, count, time.Millisecond*300, handler)
	con.RemoveHandlers()
	return messages, nil
}

func fetchMessages(con *wa.Conn, botID string, userIDs map[string]bool) (map[string]Message, error) {
	messages := make(map[string]Message)

	for userID, _ := range userIDs {
		if string(userID[0]) == "+" {
			continue
		}
		userMessages, err := loadMessages(con, botID, userID, 50)
		if err != nil {
			return messages, err
		}

		for messageID, message := range userMessages {
			messages[messageID] = message
		}
	}

	return messages, nil
}

// Message handler

func (h *messageHandler) HandleTextMessage(msg wa.TextMessage) {
	con, err := getConnection(h.botID)
	if err != nil {
		return
	}

	currentUserID := CleanPhoneNumber(con.Info.Wid) + "@s.whatsapp.net"
	message := Message{}
	message.ID = msg.Info.Id
	message.Timestamp = msg.Info.Timestamp
	message.Body = msg.Text
	contact, ok := con.Store.Contacts[msg.Info.RemoteJid]
	if ok {
		message.Name = contact.Name
	}
	if msg.Info.FromMe {
		message.Source = currentUserID
		message.Recipient = msg.Info.RemoteJid
	} else {
		message.Source = msg.Info.RemoteJid
		message.Recipient = currentUserID
	}

	h.userIDs[msg.Info.RemoteJid] = true
	h.messages[message.ID] = message
}

func (h *messageHandler) HandleError(err error) {
	if e, ok := err.(*wa.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		<-time.After(10 * time.Second)
		log.Println("Reconnecting...")

		con, err := getConnection(h.botID)
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}

		err = con.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else if strings.Contains(err.Error(), "tag 174") {
		log.Printf("Binary decode error, underlying error: %v", err)
	} else {
		log.Printf("Message handler error: %v\n", err)
	}
}

func (h *messageHandler) ShouldCallSynchronously() bool {
	return h.synchronous
}
