package responders

import (
	"encoding/json"
	"github.com/Rhymen/go-whatsapp"
	"gitlab.com/digiresilience/link/quepasa/models"
	"log"
	"net/http"
	"strings"
)

func init() {
	models.RegisterResponder(func(msg whatsapp.TextMessage, conn whatsapp.Conn) bool {
		if !strings.HasPrefix(msg.Text, "#moppenbot") && !strings.HasPrefix(msg.Text, "#mopbot") && !strings.HasPrefix(msg.Text, "#grapjesbot") && !strings.HasPrefix(msg.Text, "#grapbot") {
			return false
		}

		text := "Helaas, vandaag geen mop!"

		response, err := http.Get("http://api.apekool.nl/services/jokes/getjoke.php")

		if err == nil {
			type MopStruct struct {
				Status string
				Joke string
			}

			var data MopStruct
			json.NewDecoder(response.Body).Decode(&data)

			if data.Status == "ok" {
				text = data.Joke
			}
		}

		_, err = conn.Send(whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: msg.Info.RemoteJid,
			},
			Text: text,
		})

		if err != nil {
			log.Printf("Error while sending message for is_op: %s", err)
		}

		return true
	}, 998)
}
