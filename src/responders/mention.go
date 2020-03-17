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
		mention := fmt.Sprintf("@%s", strings.Split(conn.Info.Wid, "@")[0])

		log.Printf("line: %s with mention %s", msg.Text, mention)

		if !strings.Contains(msg.Text, mention) {
			return false
		}

		_, err := conn.Send(whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: msg.Info.RemoteJid,
			},
			Text: "Leuk dat je me riep in een bericht, maar ik weet niet wat je bedoelt.",
		})

		if err != nil {
			log.Printf("Error while sending message for is_op: %s", err)
		}

		return true
	}, 999)
}
