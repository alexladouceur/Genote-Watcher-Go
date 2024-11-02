package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func NotifyUser(url, courseCode string) {
	data := []byte(fmt.Sprintf(`
	{
    "content": "<@134425266083594241>",
    "embeds": [
      {
        "description": "Nouvelle note en %s est disponible sur Genote!",
        "color": 100425
      }
    ],
    "attachments": []
	}
	`, courseCode))

	contentType := "application/json"
	r, err := http.Post(url, contentType, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode == 204 {
		log.Println("Notification sent successfully")
	} else {
		log.Println("Failed to send notification")
	}
}
