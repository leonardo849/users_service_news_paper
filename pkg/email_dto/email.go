package email_dto

type SendEmailDTO struct {
	To []string `json:"to"`
	Subject string `json:"subject"`
	Text string `json:"text"`
}