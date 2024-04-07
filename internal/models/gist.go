package models

type GistText struct {
	From_message    string `json: "from_message"`
	Subject_message string `json:"subject_message"`
	Email_message   string `json:"email_message"`
}
