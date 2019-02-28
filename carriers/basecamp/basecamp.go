package basecamp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/imdario/mergo"
	"github.com/think-it-labs/notifyme/carriers"
	"github.com/think-it-labs/notifyme/notification"
)

type Basecamp struct {
	AccessToken  string
	AccountID    string
	ClientID     string
	ClientSecret string
	RefreshToken string
	Project      string
	Board        string
	// boards       []string
}

// '{"subject":"Kickoff","content":"<div><strong>Welcome to Basecamp, everyone.</strong></div>","status":"active"}'
type baseCampMessage struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

func init() {
	carriers.RegisterCarrier("basecamp", new)
}

func new(conf map[string]interface{}) (carriers.Carrier, error) {
	// Default config
	var basecampCarrierConfig = Basecamp{}

	if err := mergo.Map(&basecampCarrierConfig, conf, mergo.WithOverride); err != nil {
		return nil, err
	}

	if basecampCarrierConfig.AccessToken == "" {
		return nil, fmt.Errorf("Basecamp: missing access token")
	}

	// Build channel list
	// for _, board := range strings.Split(basecampCarrierConfig.Boards, ",") {
	// 	basecampCarrierConfig.boards = append(basecampCarrierConfig.boards, strings.TrimSpace(board))
	// }

	return &basecampCarrierConfig, nil
}

func (c *Basecamp) Send(notif *notification.Notification) error {
	endpoint := fmt.Sprintf("https://3.basecampapi.com/%s/buckets/%s/message_boards/%s/messages.json", c.AccountID, c.Project, c.Board)

	msg := baseCampMessage{
		Subject: notif.Cmd,
		Content: fmt.Sprintf("%s", notif.Logs),
		Status:  "active",
	}
	js, _ := json.Marshal(msg)
	bytes.NewReader(js)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(js))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Basecamp: Error sending notification: %v", err)
	}
	return nil
}

type refreshTokenResponse struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
}

func (c *Basecamp) newToken() (string, error) {
	endpoint := fmt.Sprintf("https://launchpad.37signals.com/authorization/token?type=refresh&refresh_token=%s&client_id=%s&redirect_uri=%s&client_secret=%s",
		c.RefreshToken,
		c.ClientID,
		"https://example.com",
		c.ClientSecret,
	)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var jsonResp refreshTokenResponse
	json.NewDecoder(resp.Body).Decode(&jsonResp)
	return jsonResp.AccessToken, nil
}
