package models

type SMTPService struct {
	SmtpEmail    string
	SmtpPassword string
	SmtpService  string
	SmtpAddress  string
}

type CfgUsecaseData struct {
	PasswordLength   uint8
	AesEncryptionKey string
	SmtpService      *SMTPService
}
