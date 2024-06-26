package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	t.Run("Default config", func(t *testing.T) {
		os.Clearenv()
		cfg, err := parseConfig()
		if err != nil {
			t.Error()
			return
		}
		if cfg.Port != 8080 {
			t.Error("Default Port not 8080")
		}
		if len(cfg.HoneyPots) != 1 || cfg.HoneyPots[0] != "_t_email" {
			t.Error("Default HoneyPots are wrong")
		}
		if cfg.SMTPPort != 587 {
			t.Error("SMTP Port not 587")
		}
		if len(cfg.BlacklistArray) != 2 || cfg.BlacklistArray[0] != "gambling" || cfg.BlacklistArray[1] != "casino" {
			t.Error("Default Blacklist is wrong")
		}
	})
	t.Run("Correct config parsing", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("PORT", "1111")
		_ = os.Setenv("HONEYPOTS", "pot,abc")
		_ = os.Setenv("EMAIL_TO", "mail@example.com")
		_ = os.Setenv("ALLOWED_TO", "mail@example.com,test@example.com")
		_ = os.Setenv("EMAIL_FROM", "forms@example.com")
		_ = os.Setenv("SMTP_USER", "test@example.com")
		_ = os.Setenv("SMTP_PASS", "secret")
		_ = os.Setenv("SMTP_HOST", "smtp.example.com")
		_ = os.Setenv("SMTP_PORT", "100")
		_ = os.Setenv("GOOGLE_API_KEY", "abc")
		_ = os.Setenv("BLACKLIST", "test,abc")
		cfg, err := parseConfig()
		if err != nil {
			t.Error()
			return
		}
		if !reflect.DeepEqual(cfg.Port, 1111) {
			t.Error("Port is wrong")
		}
		if !reflect.DeepEqual(cfg.HoneyPots, []string{"pot", "abc"}) {
			t.Error("HoneyPots are wrong")
		}
		if !reflect.DeepEqual(cfg.DefaultRecipient, "mail@example.com") {
			t.Error("DefaultRecipient is wrong")
		}
		if !reflect.DeepEqual(cfg.AllowedRecipient, []string{"mail@example.com", "test@example.com"}) {
			t.Error("AllowedRecipient are wrong")
		}
		if !reflect.DeepEqual(cfg.Sender, "forms@example.com") {
			t.Error("Sender is wrong")
		}
		if !reflect.DeepEqual(cfg.SMTPUser, "test@example.com") {
			t.Error("SMTP user is wrong")
		}
		if !reflect.DeepEqual(cfg.SMTPPassword, "secret") {
			t.Error("SMTP password is wrong")
		}
		if !reflect.DeepEqual(cfg.SMTPHost, "smtp.example.com") {
			t.Error("SMTP host is wrong")
		}
		if !reflect.DeepEqual(cfg.SMTPPort, 100) {
			t.Error("SMTP port is wrong")
		}
		if !reflect.DeepEqual(cfg.GoogleAPIKey, "abc") {
			t.Error("Google API Key is wrong")
		}
		if !reflect.DeepEqual(cfg.BlacklistArray, []string{"test", "abc"}) {
			t.Error("Blacklist is wrong")
		}
	})
	t.Run("Error when wrong config", func(t *testing.T) {
		os.Clearenv()
		_ = os.Setenv("PORT", "ABC")
		_, err := parseConfig()
		if err == nil {
			t.Error()
		}
	})
}

func Test_checkRequiredConfig(t *testing.T) {
	validConfig := &config{
		Port:             8080,
		HoneyPots:        []string{"_t_email"},
		DefaultRecipient: "mail@example.com",
		AllowedRecipient: []string{"mail@example.com"},
		Sender:           "forms@example.com",
		SMTPUser:         "test@example.com",
		SMTPPassword:     "secret",
		SMTPHost:         "smtp.example.com",
		SMTPPort:         587,
	}
	t.Run("Valid config", func(t *testing.T) {
		if true != checkRequiredConfig(validConfig) {
			t.Error()
		}
	})
	t.Run("Default recipient missing", func(t *testing.T) {
		newConfig := *validConfig
		newConfig.DefaultRecipient = ""
		if false != checkRequiredConfig(&newConfig) {
			t.Error()
		}
	})
	t.Run("Allowed recipients missing", func(t *testing.T) {
		newConfig := *validConfig
		newConfig.AllowedRecipient = nil
		if false != checkRequiredConfig(&newConfig) {
			t.Error()
		}
	})
	t.Run("Sender missing", func(t *testing.T) {
		newConfig := *validConfig
		newConfig.Sender = ""
		if false != checkRequiredConfig(&newConfig) {
			t.Error()
		}
	})
	t.Run("SMTP user missing", func(t *testing.T) {
		newConfig := *validConfig
		newConfig.SMTPUser = ""
		if false != checkRequiredConfig(&newConfig) {
			t.Error()
		}
	})
	t.Run("SMTP password missing", func(t *testing.T) {
		newConfig := *validConfig
		newConfig.SMTPPassword = ""
		if false != checkRequiredConfig(&newConfig) {
			t.Error()
		}
	})
	t.Run("SMTP host missing", func(t *testing.T) {
		newConfig := *validConfig
		newConfig.SMTPHost = ""
		if false != checkRequiredConfig(&newConfig) {
			t.Error()
		}
	})
}
