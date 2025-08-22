package xray

import (
	"encoding/json"
	"testing"
)

func TestProtocolCompatibility(t *testing.T) {
	config := NewConfig("info")
	
	// Test that all protocol methods return valid JSON structures
	protocols := []struct {
		name string
		testFunc func() interface{}
	}{
		{
			"VLESS Inbound",
			func() interface{} {
				return config.MakeVlessInbound("test", 10001, "550e8400-e29b-41d4-a716-446655440000", "tcp", nil)
			},
		},
		{
			"VLESS Outbound",
			func() interface{} {
				return config.MakeVlessOutbound("test", "example.com", 443, "550e8400-e29b-41d4-a716-446655440000", "tcp")
			},
		},
		{
			"VMess Inbound",
			func() interface{} {
				return config.MakeVmessInbound("test", 10002, "550e8400-e29b-41d4-a716-446655440000", "auto", "tcp")
			},
		},
		{
			"VMess Outbound",
			func() interface{} {
				return config.MakeVmessOutbound("test", "example.com", 443, "550e8400-e29b-41d4-a716-446655440000", "auto", "tcp")
			},
		},
		{
			"Trojan Inbound",
			func() interface{} {
				return config.MakeTrojanInbound("test", 10003, "password123", "tcp", nil)
			},
		},
		{
			"Trojan Outbound",
			func() interface{} {
				return config.MakeTrojanOutbound("test", "example.com", 443, "password123", "tcp")
			},
		},
		{
			"Shadowsocks Inbound",
			func() interface{} {
				return config.MakeShadowsocksInbound("test", "password123", "aes-256-gcm", "tcp", 10004, []*Client{})
			},
		},
		{
			"Shadowsocks Outbound",
			func() interface{} {
				return config.MakeShadowsocksOutbound("test", "example.com", "password123", "aes-256-gcm", 443)
			},
		},
	}
	
	for _, protocol := range protocols {
		t.Run(protocol.name, func(t *testing.T) {
			result := protocol.testFunc()
			
			// Test JSON serialization
			jsonData, err := json.Marshal(result)
			if err != nil {
				t.Errorf("Failed to serialize %s: %v", protocol.name, err)
				return
			}
			
			// Test JSON deserialization
			var unmarshaled interface{}
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to deserialize %s: %v", protocol.name, err)
				return
			}
			
			t.Logf("%s JSON test passed", protocol.name)
		})
	}
}

func TestConfigWithAllProtocols(t *testing.T) {
	config := NewConfig("info")
	
	// Add all protocol types
	vlessInbound := config.MakeVlessInbound("vless-in", 10001, "550e8400-e29b-41d4-a716-446655440000", "tcp", nil)
	vmessInbound := config.MakeVmessInbound("vmess-in", 10002, "550e8400-e29b-41d4-a716-446655440001", "auto", "tcp")
	trojanInbound := config.MakeTrojanInbound("trojan-in", 10003, "password123", "tcp", nil)
	ssInbound := config.MakeShadowsocksInbound("ss-in", "password456", "aes-256-gcm", "tcp", 10004, []*Client{})
	
	config.Inbounds = append(config.Inbounds, vlessInbound, vmessInbound, trojanInbound, ssInbound)
	
	vlessOutbound := config.MakeVlessOutbound("vless-out", "example.com", 443, "550e8400-e29b-41d4-a716-446655440002", "tcp")
	vmessOutbound := config.MakeVmessOutbound("vmess-out", "example.com", 443, "550e8400-e29b-41d4-a716-446655440003", "auto", "tcp")
	trojanOutbound := config.MakeTrojanOutbound("trojan-out", "example.com", 443, "password789", "tcp")
	ssOutbound := config.MakeShadowsocksOutbound("ss-out", "example.com", "password101112", "aes-256-gcm", 443)
	
	config.Outbounds = append(config.Outbounds, vlessOutbound, vmessOutbound, trojanOutbound, ssOutbound)
	
	// Test full config validation
	err := config.Validate()
	if err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
	
	// Test JSON serialization of full config
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Errorf("Failed to serialize full config: %v", err)
	}
	
	// Test finding inbounds/outbounds
	if foundInbound := config.FindInbound("vless-in"); foundInbound == nil {
		t.Error("Failed to find VLESS inbound")
	}
	
	if foundOutbound := config.FindOutbound("trojan-out"); foundOutbound == nil {
		t.Error("Failed to find Trojan outbound")
	}
	
	t.Logf("Full config test passed, JSON size: %d bytes", len(jsonData))
}

func TestProtocolSpecificFeatures(t *testing.T) {
	config := NewConfig("info")
	
	// Test VLESS specific features
	t.Run("VLESS Features", func(t *testing.T) {
		inbound := config.MakeVlessInbound("vless-test", 10001, "550e8400-e29b-41d4-a716-446655440000", "tcp", nil)
		
		if inbound.Settings.Decryption != "none" {
			t.Error("VLESS should have decryption set to 'none'")
		}
		
		if len(inbound.Settings.Clients) != 1 {
			t.Error("VLESS should have exactly one client")
		}
		
		client := inbound.Settings.Clients[0]
		if client.ID == "" {
			t.Error("VLESS client should have UUID in ID field")
		}
		
		if client.Password != "" {
			t.Error("VLESS client should have empty password field")
		}
	})
	
	// Test VMess specific features
	t.Run("VMess Features", func(t *testing.T) {
		inbound := config.MakeVmessInbound("vmess-test", 10002, "550e8400-e29b-41d4-a716-446655440000", "auto", "tcp")
		
		client := inbound.Settings.Clients[0]
		if client.ID == "" {
			t.Error("VMess client should have UUID in ID field")
		}
		
		if client.AlterId != 0 {
			t.Error("VMess client should have AlterId set to 0")
		}
		
		if client.Security != "auto" {
			t.Error("VMess client should have correct encryption security")
		}
	})
	
	// Test Trojan specific features
	t.Run("Trojan Features", func(t *testing.T) {
		inbound := config.MakeTrojanInbound("trojan-test", 10003, "mypassword", "tcp", nil)
		
		client := inbound.Settings.Clients[0]
		if client.Password != "mypassword" {
			t.Error("Trojan client should have password set correctly")
		}
		
		if client.Method != "" {
			t.Error("Trojan client should not have method field set")
		}
		
		if client.ID != "" {
			t.Error("Trojan client should not have ID field set")
		}
	})
}
