// Statup
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statup
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package notifiers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/statping/statping/types/failures"
	"github.com/statping/statping/types/notifications"
	"github.com/statping/statping/types/notifier"
	"github.com/statping/statping/types/services"
	"github.com/statping/statping/utils"
	"net/url"
	"strings"
	"time"
)

var _ notifier.Notifier = (*telegram)(nil)

type telegram struct {
	*notifications.Notification
}

func (t *telegram) Select() *notifications.Notification {
	return t.Notification
}

var Telegram = &telegram{&notifications.Notification{
	Method:      "telegram",
	Title:       "Telegram",
	Description: "Receive notifications on your Telegram channel when a service has an issue. You must get a Telegram API token from the /botfather. Review the <a target=\"_blank\" href=\"http://techthoughts.info/how-to-create-a-telegram-bot-and-send-messages-via-api\">Telegram API Tutorial</a> to learn how to generate a new API Token.",
	Author:      "Hunter Long",
	AuthorUrl:   "https://github.com/hunterlong",
	Icon:        "fab fa-telegram-plane",
	Delay:       time.Duration(5 * time.Second),
	Limits:      60,
	Form: []notifications.NotificationForm{{
		Type:        "text",
		Title:       "Telegram API Token",
		Placeholder: "383810182:EEx829dtCeufeQYXG7CUdiQopqdmmxBPO7-s",
		SmallText:   "Enter the API Token given to you from the /botfather chat.",
		DbField:     "api_secret",
		Required:    true,
	}, {
		Type:        "text",
		Title:       "Channel or User ID",
		Placeholder: "789325392",
		SmallText:   "Insert your Telegram Channel ID or User ID here.",
		DbField:     "Var1",
		Required:    true,
	}}},
}

// Send will send a HTTP Post to the Telegram API. It accepts type: string
func (u *telegram) sendMessage(message string) error {
	apiEndpoint := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", u.ApiSecret)

	v := url.Values{}
	v.Set("chat_id", u.Var1)
	v.Set("text", message)
	rb := *strings.NewReader(v.Encode())

	contents, _, err := utils.HttpRequest(apiEndpoint, "GET", "application/x-www-form-urlencoded", nil, &rb, time.Duration(10*time.Second), true)

	success, _ := telegramSuccess(contents)
	if !success {
		errorOut := telegramError(contents)
		out := fmt.Sprintf("Error code %v - %v", errorOut.ErrorCode, errorOut.Description)
		return errors.New(out)
	}
	return err
}

// OnFailure will trigger failing service
func (u *telegram) OnFailure(s *services.Service, f *failures.Failure) error {
	msg := fmt.Sprintf("Your service '%v' is currently offline!", s.Name)
	return u.sendMessage(msg)
}

// OnSuccess will trigger successful service
func (u *telegram) OnSuccess(s *services.Service) error {
	msg := fmt.Sprintf("Your service '%v' is currently online!", s.Name)
	return u.sendMessage(msg)
}

// OnTest will test the Twilio SMS messaging
func (u *telegram) OnTest() error {
	msg := fmt.Sprintf("Testing the Twilio SMS Notifier on your Statping server")
	return u.sendMessage(msg)
}

func telegramSuccess(res []byte) (bool, telegramResponse) {
	var obj telegramResponse
	json.Unmarshal(res, &obj)
	if obj.Ok {
		return true, obj
	}
	return false, obj
}

func telegramError(res []byte) telegramErrorObj {
	var obj telegramErrorObj
	json.Unmarshal(res, &obj)
	return obj
}

type telegramErrorObj struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type telegramResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}
