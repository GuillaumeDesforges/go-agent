package agent

import (
	"testing"
)

func TestAgentTell_Success(t *testing.T) {
	conversation := NewConversationFake()
	agent := NewReAct(conversation)

	err := agent.Tell(t.Context(), "This is a trivial test.")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
