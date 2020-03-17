package responders

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"gitlab.com/digiresilience/link/quepasa/models"
	"log"
	"strings"
)

func init() {
	models.RegisterResponder(func(msg whatsapp.TextMessage, conn whatsapp.Conn) bool {
		if !strings.HasPrefix(msg.Text, "#") || !strings.HasSuffix(msg.Text, "bot") {
			return false
		}

		what := strings.TrimPrefix(strings.TrimSuffix(msg.Text, "bot"), "#")

		if strings.TrimSpace(what) == "" {
			return false
		}

		text := fmt.Sprintf("De %s is op!", what)

		if what == "kut" {
			text = "De vrouwen zijn op!"
		}

		if what == "lul" {
			text = "De mannen zijn op!"
		}

		if what == "wc" {
			text = "De wc papier is op!"
		}

		log.Printf("Responder is_op answering with: %s", text)

		_, err := conn.Send(whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: msg.Info.RemoteJid,
			},
			Text: text,
		})

		if err != nil {
			log.Printf("Error while sending message for is_op: %s", err)
		}

		return true
	}, 999)
}
