package xray

import (
	"testing"
)

func TestProtocolValidation(t *testing.T) {
	config := NewConfig("info")

	// Test VMess configuration without required UUID
	t.Run("VMess Missing UUID", func(t *testing.T) {
		vmessInbound := &Inbound{
			Tag:      "vmess-invalid",
			Protocol: "vmess",
			Listen:   "0.0.0.0",
			Port:     10001,
			Settings: &InboundSettings{
				Clients: []*Client{
					{
						Email: "test@example.com",
						// Missing ID field
					},
				},
			},
		}
		config.Inbounds = append(config.Inbounds, vmessInbound)
		
		err := config.validateProtocolSpecific()
		if err == nil {
			t.Error("Expected validation error for VMess without UUID")
		}
		
		// Reset for next test
		config.Inbounds = config.Inbounds[:len(config.Inbounds)-1]
	})

	// Test Shadowsocks configuration without method
	t.Run("Shadowsocks Missing Method", func(t *testing.T) {
		ssOutbound := &Outbound{
			Tag:      "ss-invalid",
			Protocol: "shadowsocks",
			Settings: &OutboundSettings{
				Servers: []*OutboundServer{
					{
						Address:  "example.com",
						Port:     443,
						Password: "password123",
						// Missing Method field
					},
				},
			},
		}
		config.Outbounds = append(config.Outbounds, ssOutbound)
		
		err := config.validateProtocolSpecific()
		if err == nil {
			t.Error("Expected validation error for Shadowsocks without method")
		}
		
		// Reset for next test
		config.Outbounds = config.Outbounds[:len(config.Outbounds)-1]
	})

	// Test valid VMess configuration
	t.Run("Valid VMess", func(t *testing.T) {
		vmessInbound := config.MakeVmessInbound("vmess-valid", 10002, "550e8400-e29b-41d4-a716-446655440000", "auto", "tcp")
		config.Inbounds = append(config.Inbounds, vmessInbound)
		
		err := config.validateProtocolSpecific()
		if err != nil {
			t.Errorf("Unexpected validation error for valid VMess: %v", err)
		}
		
		// Reset for next test
		config.Inbounds = config.Inbounds[:len(config.Inbounds)-1]
	})

	// Test valid Shadowsocks configuration
	t.Run("Valid Shadowsocks", func(t *testing.T) {
		ssOutbound := config.MakeShadowsocksOutbound("ss-valid", "example.com", "password123", "aes-256-gcm", 443)
		config.Outbounds = append(config.Outbounds, ssOutbound)
		
		err := config.validateProtocolSpecific()
		if err != nil {
			t.Errorf("Unexpected validation error for valid Shadowsocks: %v", err)
		}
		
		// Reset for next test
		config.Outbounds = config.Outbounds[:len(config.Outbounds)-1]
	})

	// Test unknown protocol (should pass)
	t.Run("Unknown Protocol", func(t *testing.T) {
		unknownInbound := &Inbound{
			Tag:      "unknown",
			Protocol: "unknown-protocol",
			Listen:   "0.0.0.0",
			Port:     10003,
			Settings: &InboundSettings{
				Clients: []*Client{
					{
						Email: "test@example.com",
					},
				},
			},
		}
		config.Inbounds = append(config.Inbounds, unknownInbound)
		
		err := config.validateProtocolSpecific()
		if err != nil {
			t.Errorf("Unexpected validation error for unknown protocol: %v", err)
		}
		
		// Reset
		config.Inbounds = config.Inbounds[:len(config.Inbounds)-1]
	})
}
