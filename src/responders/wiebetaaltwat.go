package responders

import (
	"github.com/Rhymen/go-whatsapp"
	"gitlab.com/digiresilience/link/quepasa/models"
	"log"
	"strconv"
	"strings"
)

type WieBetaaltWatParticipant struct {
	ID string
	Amount
}

type WieBetaaltWatEntry struct {
	TotalAmount float64
	TotalParticipants int
	Participants []WieBetaaltWatParticipant
}

func init() {
	models.RegisterResponder(func(msg whatsapp.TextMessage, conn whatsapp.Conn) bool {
		if !strings.HasPrefix(msg.Text, "#wbw ") {
			return false
		}

		elements := strings.Split(strings.TrimPrefix(msg.Text, "#wbw "), " ")

		for _, element := range elements {
			if strings.TrimSpace(element) == "" {
				continue
			}

			price, err := strconv.ParseFloat(element, 64)
			if err == nil {
				totalAmount += price
				continue
			}

		}

		log.Printf("Total amount: %i", totalAmount)

		return true
	}, 50)
}
