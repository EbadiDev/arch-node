package xray

import (
	"testing"
)

func TestValidateProtocolCompatibility(t *testing.T) {
	config := NewConfig("info")

	tests := []struct {
		name        string
		protocol    string
		streamSettings *StreamSettings
		expectError bool
		errorContains string
	}{
		{
			name:     "VMess with basic TCP should work",
			protocol: "vmess", 
			streamSettings: &StreamSettings{Network: "tcp"},
			expectError: false,
		},
		{
			name:     "VMess with WebSocket should work",
			protocol: "vmess",
			streamSettings: &StreamSettings{Network: "ws"},
			expectError: false,
		},
		{
			name:     "VMess with TLS should work",
			protocol: "vmess",
			streamSettings: &StreamSettings{Network: "tcp", Security: "tls"},
			expectError: false,
		},
		{
			name:     "VMess with REALITY should fail",
			protocol: "vmess",
			streamSettings: &StreamSettings{
				Network: "tcp", 
				Security: "reality",
				RealitySettings: &RealitySettings{},
			},
			expectError: true,
			errorContains: "VMess protocol does not support REALITY",
		},
		{
			name:     "VMess with XHTTP should fail",
			protocol: "vmess",
			streamSettings: &StreamSettings{
				Network: "xhttp",
				XhttpSettings: &XhttpSettings{},
			},
			expectError: true,
			errorContains: "VMess protocol does not support XHTTP",
		},
		{
			name:     "VLESS with REALITY should work",
			protocol: "vless",
			streamSettings: &StreamSettings{
				Network: "tcp",
				Security: "reality", 
				RealitySettings: &RealitySettings{},
			},
			expectError: false,
		},
		{
			name:     "VLESS with XHTTP should work",
			protocol: "vless",
			streamSettings: &StreamSettings{
				Network: "xhttp",
				XhttpSettings: &XhttpSettings{},
			},
			expectError: false,
		},
		{
			name:     "Trojan with non-TLS should fail",
			protocol: "trojan",
			streamSettings: &StreamSettings{
				Network: "tcp",
				Security: "reality",
			},
			expectError: true,
			errorContains: "Trojan protocol typically only supports TLS",
		},
		{
			name:     "Trojan with non-TCP should fail",
			protocol: "trojan",
			streamSettings: &StreamSettings{
				Network: "ws",
				Security: "tls",
			},
			expectError: true,
			errorContains: "Trojan protocol typically only supports TCP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.ValidateProtocolCompatibility(tt.protocol, tt.streamSettings)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.name)
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
			}
		})
	}
}

func TestSanitizeStreamSettingsForProtocol(t *testing.T) {
	config := NewConfig("info")

	t.Run("VMess with REALITY should fallback to TLS", func(t *testing.T) {
		original := &StreamSettings{
			Network: "tcp",
			Security: "reality",
			RealitySettings: &RealitySettings{Dest: "example.com"},
		}
		
		sanitized := config.SanitizeStreamSettingsForProtocol("vmess", original)
		
		if sanitized.Security != "tls" {
			t.Errorf("Expected security to be 'tls', got '%s'", sanitized.Security)
		}
		if sanitized.RealitySettings != nil {
			t.Errorf("Expected RealitySettings to be nil, got %v", sanitized.RealitySettings)
		}
	})

	t.Run("VMess with XHTTP should fallback to WebSocket", func(t *testing.T) {
		original := &StreamSettings{
			Network: "xhttp",
			XhttpSettings: &XhttpSettings{Host: "example.com"},
		}
		
		sanitized := config.SanitizeStreamSettingsForProtocol("vmess", original)
		
		if sanitized.Network != "ws" {
			t.Errorf("Expected network to be 'ws', got '%s'", sanitized.Network)
		}
		if sanitized.XhttpSettings != nil {
			t.Errorf("Expected XhttpSettings to be nil, got %v", sanitized.XhttpSettings)
		}
	})

	t.Run("VLESS should preserve all settings", func(t *testing.T) {
		original := &StreamSettings{
			Network: "xhttp",
			Security: "reality",
			XhttpSettings: &XhttpSettings{Host: "example.com"},
			RealitySettings: &RealitySettings{Dest: "example.com"},
		}
		
		sanitized := config.SanitizeStreamSettingsForProtocol("vless", original)
		
		if sanitized.Network != "xhttp" {
			t.Errorf("Expected network to be 'xhttp', got '%s'", sanitized.Network)
		}
		if sanitized.Security != "reality" {
			t.Errorf("Expected security to be 'reality', got '%s'", sanitized.Security)
		}
		if sanitized.XhttpSettings == nil {
			t.Error("Expected XhttpSettings to be preserved")
		}
		if sanitized.RealitySettings == nil {
			t.Error("Expected RealitySettings to be preserved")
		}
	})
}

func TestVmessVlessCompatibilityMethods(t *testing.T) {
	config := NewConfig("info")

	t.Run("VMess convenience methods should work", func(t *testing.T) {
		// Test VMess WebSocket using composable approach
		wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
		inbound := config.MakeVmessInbound("test", 8080, "uuid", "auto", wsSettings)
		if inbound.Protocol != "vmess" {
			t.Errorf("Expected protocol 'vmess', got '%s'", inbound.Protocol)
		}
		if inbound.StreamSettings.Network != "ws" {
			t.Errorf("Expected network 'ws', got '%s'", inbound.StreamSettings.Network)
		}

		// Test VMess gRPC using composable approach
		grpcSettings := config.MakeGrpcStreamSettings("service", "authority")
		grpcInbound := config.MakeVmessInbound("test", 8080, "uuid", "auto", grpcSettings)
		if grpcInbound.StreamSettings.Network != "grpc" {
			t.Errorf("Expected network 'grpc', got '%s'", grpcInbound.StreamSettings.Network)
		}
	})

	t.Run("VLESS convenience methods should work", func(t *testing.T) {
		// Test VLESS with REALITY using composable approach (should work)
		grpcSettings := config.MakeGrpcStreamSettings("service", "")
		grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443", []string{"example.com"}, "private", "public")
		inbound := config.MakeVlessInbound("test", 8080, "uuid", "tcp", grpcSettings)
		if inbound.Protocol != "vless" {
			t.Errorf("Expected protocol 'vless', got '%s'", inbound.Protocol)
		}
		if inbound.StreamSettings.Security != "reality" {
			t.Errorf("Expected security 'reality', got '%s'", inbound.StreamSettings.Security)
		}

		// Test VLESS with XHTTP using composable approach (should work)
		xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")
		xhttpInbound := config.MakeVlessInbound("test", 8080, "uuid", "tcp", xhttpSettings)
		if xhttpInbound.StreamSettings.Network != "xhttp" {
			t.Errorf("Expected network 'xhttp', got '%s'", xhttpInbound.StreamSettings.Network)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
