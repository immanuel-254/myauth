package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/resend/resend-go/v2"
)

func GetData(data map[string]string, w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
		return
	}
}

func SendData(data map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func SendEmail(email, subject, link string, template func(route string) string, w http.ResponseWriter, r *http.Request) {
	// send email
	client := resend.NewClient(os.Getenv("RESENDAPIKEY"))

	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESENDEMAIL"),
		To:      []string{email},
		Html:    template(link),
		Subject: subject,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}