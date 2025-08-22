package xray

import (
	"encoding/json"
	"testing"
)

func TestMakeVlessInbound(t *testing.T) {
	config := NewConfig("info")
	
	inbound := config.MakeVlessInbound("vless-test", 10001, "550e8400-e29b-41d4-a716-446655440000", "tcp", nil)
	
	if inbound.Tag != "vless-test" {
		t.Errorf("Expected tag 'vless-test', got %s", inbound.Tag)
	}
	
	if inbound.Protocol != "vless" {
		t.Errorf("Expected protocol 'vless', got %s", inbound.Protocol)
	}
	
	if inbound.Port != 10001 {
		t.Errorf("Expected port 10001, got %d", inbound.Port)
	}
	
	if len(inbound.Settings.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(inbound.Settings.Clients))
	}
	
	client := inbound.Settings.Clients[0]
	if client.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected UUID, got %s", client.ID)
	}
	
	if inbound.Settings.Decryption != "none" {
		t.Errorf("Expected decryption 'none', got %s", inbound.Settings.Decryption)
	}
	
	// Test JSON serialization
	_, err := json.Marshal(inbound)
	if err != nil {
		t.Errorf("Failed to serialize VLESS inbound: %v", err)
	}
}

func TestMakeVmessInbound(t *testing.T) {
	config := NewConfig("info")
	
	inbound := config.MakeVmessInbound("vmess-test", 10002, "550e8400-e29b-41d4-a716-446655440000", "auto", "tcp")
	
	if inbound.Tag != "vmess-test" {
		t.Errorf("Expected tag 'vmess-test', got %s", inbound.Tag)
	}
	
	if inbound.Protocol != "vmess" {
		t.Errorf("Expected protocol 'vmess', got %s", inbound.Protocol)
	}
	
	if inbound.Port != 10002 {
		t.Errorf("Expected port 10002, got %d", inbound.Port)
	}
	
	if len(inbound.Settings.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(inbound.Settings.Clients))
	}
	
	client := inbound.Settings.Clients[0]
	if client.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected UUID, got %s", client.ID)
	}
	
	if client.Method != "auto" {
		t.Errorf("Expected method 'auto', got %s", client.Method)
	}
	
	// Test JSON serialization
	_, err := json.Marshal(inbound)
	if err != nil {
		t.Errorf("Failed to serialize VMess inbound: %v", err)
	}
}

func TestMakeTrojanInbound(t *testing.T) {
	config := NewConfig("info")
	
	inbound := config.MakeTrojanInbound("trojan-test", 10003, "mypassword123", "tcp", nil)
	
	if inbound.Tag != "trojan-test" {
		t.Errorf("Expected tag 'trojan-test', got %s", inbound.Tag)
	}
	
	if inbound.Protocol != "trojan" {
		t.Errorf("Expected protocol 'trojan', got %s", inbound.Protocol)
	}
	
	if inbound.Port != 10003 {
		t.Errorf("Expected port 10003, got %d", inbound.Port)
	}
	
	if len(inbound.Settings.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(inbound.Settings.Clients))
	}
	
	client := inbound.Settings.Clients[0]
	if client.Password != "mypassword123" {
		t.Errorf("Expected password 'mypassword123', got %s", client.Password)
	}
	
	// Test JSON serialization
	_, err := json.Marshal(inbound)
	if err != nil {
		t.Errorf("Failed to serialize Trojan inbound: %v", err)
	}
}

func TestMakeVlessOutbound(t *testing.T) {
	config := NewConfig("info")
	
	outbound := config.MakeVlessOutbound("vless-out", "example.com", 443, "550e8400-e29b-41d4-a716-446655440000", "tcp")
	
	if outbound.Tag != "vless-out" {
		t.Errorf("Expected tag 'vless-out', got %s", outbound.Tag)
	}
	
	if outbound.Protocol != "vless" {
		t.Errorf("Expected protocol 'vless', got %s", outbound.Protocol)
	}
	
	if len(outbound.Settings.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(outbound.Settings.Servers))
	}
	
	server := outbound.Settings.Servers[0]
	if server.Address != "example.com" {
		t.Errorf("Expected address 'example.com', got %s", server.Address)
	}
	
	if server.Port != 443 {
		t.Errorf("Expected port 443, got %d", server.Port)
	}
	
	if server.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected UUID, got %s", server.ID)
	}
}

func TestMakeVmessOutbound(t *testing.T) {
	config := NewConfig("info")
	
	outbound := config.MakeVmessOutbound("vmess-out", "example.com", 443, "550e8400-e29b-41d4-a716-446655440000", "auto", "tcp")
	
	if outbound.Tag != "vmess-out" {
		t.Errorf("Expected tag 'vmess-out', got %s", outbound.Tag)
	}
	
	if outbound.Protocol != "vmess" {
		t.Errorf("Expected protocol 'vmess', got %s", outbound.Protocol)
	}
	
	if len(outbound.Settings.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(outbound.Settings.Servers))
	}
	
	server := outbound.Settings.Servers[0]
	if server.Address != "example.com" {
		t.Errorf("Expected address 'example.com', got %s", server.Address)
	}
	
	if server.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected UUID, got %s", server.ID)
	}
	
	if server.Method != "auto" {
		t.Errorf("Expected method 'auto', got %s", server.Method)
	}
}

func TestMakeTrojanOutbound(t *testing.T) {
	config := NewConfig("info")
	
	outbound := config.MakeTrojanOutbound("trojan-out", "example.com", 443, "mypassword123", "tcp")
	
	if outbound.Tag != "trojan-out" {
		t.Errorf("Expected tag 'trojan-out', got %s", outbound.Tag)
	}
	
	if outbound.Protocol != "trojan" {
		t.Errorf("Expected protocol 'trojan', got %s", outbound.Protocol)
	}
	
	if len(outbound.Settings.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(outbound.Settings.Servers))
	}
	
	server := outbound.Settings.Servers[0]
	if server.Address != "example.com" {
		t.Errorf("Expected address 'example.com', got %s", server.Address)
	}
	
	if server.Password != "mypassword123" {
		t.Errorf("Expected password 'mypassword123', got %s", server.Password)
	}
}
