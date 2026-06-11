package usecase_test

import (
	"context"
	"testing"

	"community-forum/backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatService_SendMessage_Success(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	msg, err := svc.SendMessage(context.Background(), 1, "Hello world")
	require.NoError(t, err)
	assert.Equal(t, uint(1), msg.ID)
	assert.Equal(t, "Hello world", msg.Content)
	assert.Equal(t, uint(1), msg.AuthorID)
}

func TestChatService_GetRecentMessages_ReturnsLatest(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	_, _ = svc.SendMessage(context.Background(), 1, "first")
	_, _ = svc.SendMessage(context.Background(), 2, "second")
	_, _ = svc.SendMessage(context.Background(), 1, "third")

	messages, err := svc.GetRecentMessages(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, messages, 2)
	assert.Equal(t, "second", messages[0].Content)
	assert.Equal(t, "third", messages[1].Content)
}

func TestChatService_GetMessagesBefore_ReturnsOlder(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	m1, _ := svc.SendMessage(context.Background(), 1, "first")
	m2, _ := svc.SendMessage(context.Background(), 2, "second")
	_, _ = svc.SendMessage(context.Background(), 1, "third")

	messages, err := svc.GetMessagesBefore(context.Background(), m2.ID, 1)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	assert.Equal(t, m1.Content, messages[0].Content)
}

func TestChatService_SendMessage_Sanitization(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	// Test script stripping
	msg, err := svc.SendMessage(context.Background(), 1, "<script>alert('xss')</script>Hello")
	require.NoError(t, err)
	assert.Equal(t, "Hello", msg.Content)

	// Test HTML tag stripping (StrictPolicy strips all tags)
	msg, err = svc.SendMessage(context.Background(), 1, "<b>Bold</b> and <i>Italic</i>")
	require.NoError(t, err)
	assert.Equal(t, "Bold and Italic", msg.Content)

	// Test invalid content after sanitization
	_, err = svc.SendMessage(context.Background(), 1, "<script></script>")
	assert.Error(t, err)
	assert.Equal(t, "message content invalid after sanitization", err.Error())
}
