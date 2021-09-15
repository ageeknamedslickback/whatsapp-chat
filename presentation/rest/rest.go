package rest

import (
	"encoding/xml"
	"log"
	"net/http"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"github.com/ageeknamedslickback/whatsapp-chat/usecases"
)

// Handlers defines our REST layer interface
type Handlers interface {
	IncomingMessage(w http.ResponseWriter, r *http.Request)
}

// Rest ..
type Rest struct {
	service usecases.MessagesService
}

// NewRestHandlers ..
func NewRestHandlers(s usecases.MessagesService) *Rest {
	return &Rest{service: s}
}

// IncomingMessage is Twilio's incoming message webhook
func (rst *Rest) IncomingMessage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("failed to parse request data: %v", err)
		return
	}
	if r.Form == nil || len(r.Form) == 0 {
		return
	}

	_, err := rst.service.InboundMessages(domain.Message{
		AccountSid:       r.Form.Get("AccountSid"),
		From:             r.Form.Get("From"),
		To:               r.Form.Get("To"),
		Body:             r.Form.Get("Body"),
		NumMedia:         r.Form.Get("NumMedia"),
		NumSegments:      r.Form.Get("NumSegments"),
		APIVersion:       r.Form.Get("ApiVersion"),
		ProfileName:      r.Form.Get("ProfileName"),
		SmsMessageSid:    r.Form.Get("SmsMessageSid"),
		SmsSid:           r.Form.Get("SmsSid"),
		SmsStatus:        r.Form.Get("SmsStatus"),
		WaID:             r.Form.Get("WaID"),
		MediaContentType: r.Form.Get("MediaContentType0"),
		MediaURL:         r.Form.Get("MediaUrl0"),
	},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response domain.Response
	x, err := xml.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/xml")
	_, err = w.Write(x)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
