package service

import (
	"app/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
)

var TestCasesValidateNewMessageLength = []struct {
	name           string
	newMessage     model.NewMessage
	expectedResult bool
}{
	{
		name: "success",
		newMessage: model.NewMessage{
			ChatId:  uuid.MustParse("27397d31-9ec3-4f40-9725-2597e2b4a9ef"),
			UserId:  uuid.MustParse("27397d31-9ec3-4f40-9725-2597e2b4a9ef"),
			Content: "This message have less than 500 characters",
		},
		expectedResult: false,
	},
	{
		name: "failure",
		newMessage: model.NewMessage{
			ChatId: uuid.MustParse("27397d31-9ec3-4f40-9725-2597e2b4a9ef"),
			UserId: uuid.MustParse("27397d31-9ec3-4f40-9725-2597e2b4a9ef"),
			Content: "This message have over 500 characters, Lorem ipsum " +
				"dolor sit amet, consectetur adipiscing elit. Nulla facilisi. " +
				"Integer euismod, nisl ac elementum feugiat, ligula libero pharetra risus, " +
				"ac vehicula felis libero nec nunc. Suspendisse potenti. Quisque sit amet " +
				"justo auctor, tempus elit non, ultricies nisi. Nam ac elit id justo aliquam " +
				"sagittis. Donec dapibus, augue et gravida malesuada, elit nunc iaculis felis, " +
				"vel tempus urna elit ac odio. Fusce nec velit eu risus scelerisque cursus. " +
				"Aenean vel dui ut velit rhoncus tincidunt. Curabitur vel urna id erat vehicula gravida. " +
				"Sed gravida, risus a interdum fermentum, odio lectus sodales nunc, at egestas libero velit non turpis",
		},
		expectedResult: true,
	},
}

func TestValidateNewMessageLength(t *testing.T) {

	for _, tc := range TestCasesValidateNewMessageLength {
		t.Run(tc.name, func(t *testing.T) {
			messageService := &messageService{}

			result := messageService.validateNewMessageLength(tc.newMessage.Content)

			if result != tc.expectedResult {
				t.Errorf("scenario: %s, expected: %v, got: %v", tc.name, tc.expectedResult, result)
			}
		})
	}
}

func TestConvertNewMessageToMessage(t *testing.T) {
	newMessage := model.NewMessage{
		ChatId:  uuid.MustParse("06eed8ac-9339-4b76-b8df-b487e0276e5f"),
		UserId:  uuid.MustParse("06eed8ac-9339-4b76-b8df-b487e0276e5f"),
		Content: "Here goes message content",
	}

	messageService := &messageService{}

	messageID := uuid.MustParse("06eed8ac-9339-4b76-b8df-b487e0276e5f")
	createdAt := time.Now()

	message := messageService.convertNewMessageToMessage(newMessage, messageID, createdAt)

	var _, ok = interface{}(*message).(model.Message)
	if !ok {
		t.Errorf("Expected message to be of type model.Message, got %T", message)
	}

	if message.Id != messageID {
		t.Errorf("Expected ID to be %v, got %v", messageID, message.Id)
	}
	if message.ChatId != newMessage.ChatId {
		t.Errorf("Expected ChatId to be %v, got %v", newMessage.ChatId, message.ChatId)
	}
	if message.UserId != newMessage.UserId {
		t.Errorf("Expected UserId to be %v, got %v", newMessage.UserId, message.UserId)
	}
	if message.Content != newMessage.Content {
		t.Errorf("Expected Content to be %q, got %q", newMessage.Content, message.Content)
	}
	if message.Created.Unix() != createdAt.Unix() {
		t.Errorf("Expected CreatedAt to be %v, got %v", createdAt, message.Created)
	}
}
