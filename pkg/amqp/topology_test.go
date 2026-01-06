package amqp

import (
	"testing"
)

func TestGetRetryQueue(t *testing.T) {
	tests := []struct {
		attempt  int
		expected string
	}{
		{0, Retry5sQueue},
		{1, Retry30sQueue},
		{2, Retry5mQueue},
		{3, ""},
		{4, ""},
	}

	for _, tt := range tests {
		result := GetRetryQueue(tt.attempt)
		if result != tt.expected {
			t.Errorf("GetRetryQueue(%d) = %s; want %s", tt.attempt, result, tt.expected)
		}
	}
}

func TestGetRetryExchangeAndKey(t *testing.T) {
	tests := []struct {
		attempt         int
		expectedExch    string
		expectedKey     string
	}{
		{0, RetryExchange, Retry5sQueue},
		{1, RetryExchange, Retry30sQueue},
		{2, RetryExchange, Retry5mQueue},
		{3, DLQExchange, DLQRoutingKey},
	}

	for _, tt := range tests {
		exch, key := GetRetryExchangeAndKey(tt.attempt)
		if exch != tt.expectedExch || key != tt.expectedKey {
			t.Errorf("GetRetryExchangeAndKey(%d) = (%s, %s); want (%s, %s)",
				tt.attempt, exch, key, tt.expectedExch, tt.expectedKey)
		}
	}
}

func TestDefaultTopologyConfig(t *testing.T) {
	config := DefaultTopologyConfig()

	if config.Retry5sTTL != 5000 {
		t.Errorf("Retry5sTTL = %d; want 5000", config.Retry5sTTL)
	}

	if config.Retry30sTTL != 30000 {
		t.Errorf("Retry30sTTL = %d; want 30000", config.Retry30sTTL)
	}

	if config.Retry5mTTL != 300000 {
		t.Errorf("Retry5mTTL = %d; want 300000", config.Retry5mTTL)
	}
}


