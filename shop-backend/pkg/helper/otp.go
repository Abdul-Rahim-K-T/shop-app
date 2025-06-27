package helper

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"shop-backend/config"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000)) // 6 dgt
	log.Printf("Generated OTP: %s\n", otp)
	return otp
}

func SendOtpEmail(cfg *config.Config, email, otp string) error {
	// Simulate email send

	from := cfg.EmailFrom
	password := cfg.EmailPassword
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	fmt.Printf("EMAIL_FROM: %s\n", from)
	fmt.Printf("EMAIL_PASSWORD: %s\n", password)
	fmt.Printf("SMTP_HOST: %s\n", smtpHost)
	fmt.Printf("SMTP_PORT: %s\n", smtpPort)

	to := []string{email}

	subject := "Your OTP code for KMS"
	body := fmt.Sprintf("Your OTP code is: %s", otp)

	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send OTP email to %s: %v\n", email, err)
	} else {
		log.Printf("Sent OTP %s to email: %s\n", otp, email)
	}
	return err
}

func SendOtpSMS(cfg *config.Config, phone string, otp string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(cfg.TwilioPhoneNumber)
	params.SetBody(fmt.Sprintf("Your OTP code is: %s", otp))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send OTP SMS to %s: %v\n", phone, err)
		return err
	}

	log.Printf("Sent OTP %s to phone %s. SID: %s\n", otp, phone, *resp.Sid)
	return nil
}

/*
func SendOtpSMS(cfg *config.Config, phone string, otp string) error {
	endpoint := "https://www.fast2sms.com/dev/bulkV2"
	params := url.Values{}
	FastSMSApiKey := strings.TrimSpace(cfg.Fast2SMSAPIKey)
	fmt.Printf("Fast2smsAPIKey for debugging:'%s'\n", FastSMSApiKey)
	params.Add("authorization", cfg.Fast2SMSAPIKey)
	params.Add("route", "q")
	params.Add("message", fmt.Sprintf("Your OTP is: %s", otp))
	params.Add("language", "english")
	params.Add("flash", "0")
	params.Add("numbers", phone)

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		log.Println("Failed to build Fast2SMS request:", err)
		return err
	}

	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("authorization", FastSMSApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send SMS:", err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Sent OTP %s to phone: %s, status: %s, response: %s\n", otp, phone, resp.Status, string(bodyBytes))

	return nil
}
*/

/*
func SendOtpSMS(cfg *config.Config, phone, otp string) error {
	// Simulate SMS send

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	message := fmt.Sprintf("Your KMS OTP is: %s", otp)

	fmt.Printf("TWILIO_ACCOUNT_SID: %s", cfg.TwilioAccountSID)
	fmt.Printf("TWILIO_AUTH_TOKEN: %s", cfg.TwilioAuthToken)

	params := &openapi.CreateMessagParams{}
	params.SetTo(phone)
	// params.SetChannel("sms")
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(os.Getenv("TWILIO_VERIFY_SERVICE_SID"), params)
	if err != nil {
		log.Printf("Failed to send OTP SMS to %s: %v\n", phone, err)
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	if resp.Sid != nil {
		log.Printf("Sent OTP SID: %s to phone: %s\n", *resp.Sid, phone)
	} else {
		log.Printf("Sent OTP to phone: %s (SID not returned)\n", phone)
	}

	log.Printf("Sent OTP %s to phone: %s\n", otp, phone)
	return nil
}
*/
