package xray

import (
	"testing"
)

// TestComposableTransportExamples demonstrates the flexible composable approach
func TestComposableTransportExamples(t *testing.T) {
	config := NewConfig("info")

	t.Run("VMess with various transports", func(t *testing.T) {
		// VMess with basic TCP
		tcpInbound := config.MakeVmessInbound("vmess-tcp", 8080, "uuid", "auto", nil)
		if tcpInbound.Protocol != "vmess" {
			t.Errorf("Expected VMess protocol")
		}

		// VMess with WebSocket
		wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
		wsInbound := config.MakeVmessInbound("vmess-ws", 8081, "uuid", "auto", wsSettings)
		if wsInbound.StreamSettings.Network != "ws" {
			t.Errorf("Expected WebSocket transport")
		}

		// VMess with WebSocket + TLS
		wsSettings = config.MakeWebSocketStreamSettings("/secure", "secure.host.com")
		wsSettings = config.AddTlsToStreamSettings(wsSettings, "secure.host.com", false)
		wssTlsInbound := config.MakeVmessInbound("vmess-wss", 8082, "uuid", "auto", wsSettings)
		if wssTlsInbound.StreamSettings.Security != "tls" {
			t.Errorf("Expected TLS security")
		}

		// VMess with gRPC + TLS
		grpcSettings := config.MakeGrpcStreamSettings("gunService", "host.com")
		grpcSettings = config.AddTlsToStreamSettings(grpcSettings, "host.com", false)
		grpcTlsInbound := config.MakeVmessInbound("vmess-grpc", 8083, "uuid", "auto", grpcSettings)
		if grpcTlsInbound.StreamSettings.Network != "grpc" {
			t.Errorf("Expected gRPC transport")
		}

		// VMess with KCP
		kcpSettings := config.MakeKcpStreamSettings("utp", "password123")
		kcpInbound := config.MakeVmessInbound("vmess-kcp", 8084, "uuid", "auto", kcpSettings)
		if kcpInbound.StreamSettings.Network != "kcp" {
			t.Errorf("Expected KCP transport")
		}

		// VMess with HTTP Upgrade
		httpSettings := config.MakeHttpUpgradeStreamSettings("/upgrade", "host.com")
		httpInbound := config.MakeVmessInbound("vmess-http", 8085, "uuid", "auto", httpSettings)
		if httpInbound.StreamSettings.Network != "httpupgrade" {
			t.Errorf("Expected HTTP Upgrade transport")
		}
	})

	t.Run("VLESS with advanced features", func(t *testing.T) {
		// VLESS with basic TCP
		tcpInbound := config.MakeVlessInbound("vless-tcp", 9080, "uuid", "tcp", nil)
		if tcpInbound.Protocol != "vless" {
			t.Errorf("Expected VLESS protocol")
		}

		// VLESS with WebSocket + TLS
		wsSettings := config.MakeWebSocketStreamSettings("/vless", "vless.host.com")
		wsSettings = config.AddTlsToStreamSettings(wsSettings, "vless.host.com", false)
		vlessWsInbound := config.MakeVlessInbound("vless-ws", 9081, "uuid", "tcp", wsSettings)
		if vlessWsInbound.StreamSettings.Security != "tls" {
			t.Errorf("Expected TLS security")
		}

		// VLESS with gRPC + REALITY (VMess can't do this!)
		grpcSettings := config.MakeGrpcStreamSettings("vlessService", "")
		grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443",
			[]string{"example.com", "www.example.com"}, "privateKey", "publicKey")
		vlessRealityInbound := config.MakeVlessInbound("vless-reality", 9082, "uuid", "tcp", grpcSettings)
		if vlessRealityInbound.StreamSettings.Security != "reality" {
			t.Errorf("Expected REALITY security")
		}

		// VLESS with XHTTP (VMess can't do this!)
		xhttpSettings := config.MakeXhttpStreamSettings("cdn.host.com", "/api", "auto")
		vlessXhttpInbound := config.MakeVlessInbound("vless-xhttp", 9083, "uuid", "tcp", xhttpSettings)
		if vlessXhttpInbound.StreamSettings.Network != "xhttp" {
			t.Errorf("Expected XHTTP transport")
		}

		// VLESS with XHTTP + REALITY (most advanced combination!)
		xhttpSettings = config.MakeXhttpStreamSettings("advanced.host.com", "/advanced", "auto")
		xhttpSettings = config.AddRealityToStreamSettings(xhttpSettings, "cloudflare.com:443",
			[]string{"cloudflare.com", "www.cloudflare.com"}, "advancedPrivateKey", "advancedPublicKey")
		vlessAdvancedInbound := config.MakeVlessInbound("vless-advanced", 9084, "uuid", "tcp", xhttpSettings)
		if vlessAdvancedInbound.StreamSettings.Network != "xhttp" {
			t.Errorf("Expected XHTTP transport")
		}
		if vlessAdvancedInbound.StreamSettings.Security != "reality" {
			t.Errorf("Expected REALITY security")
		}
	})

	t.Run("Reusing transport configurations", func(t *testing.T) {
		// Create a transport configuration once
		wsSettings := config.MakeWebSocketStreamSettings("/shared", "shared.host.com")
		wsSettings = config.AddTlsToStreamSettings(wsSettings, "shared.host.com", false)

		// Reuse it for both VMess and VLESS
		vmessInbound := config.MakeVmessInbound("shared-vmess", 7080, "vmess-uuid", "auto", wsSettings)
		vlessInbound := config.MakeVlessInbound("shared-vless", 7081, "vless-uuid", "tcp", wsSettings)

		// Both should have the same transport configuration
		if vmessInbound.StreamSettings.Network != vlessInbound.StreamSettings.Network {
			t.Errorf("Transport networks don't match")
		}
		if vmessInbound.StreamSettings.Security != vlessInbound.StreamSettings.Security {
			t.Errorf("Security settings don't match")
		}
	})
}

// TestInvalidCombinationsHandling tests that invalid combinations are handled gracefully
func TestInvalidCombinationsHandling(t *testing.T) {
	config := NewConfig("info")

	t.Run("VMess with REALITY should be sanitized", func(t *testing.T) {
		// Try to create VMess with REALITY (should be auto-corrected)
		grpcSettings := config.MakeGrpcStreamSettings("service", "")
		grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443",
			[]string{"example.com"}, "privateKey", "publicKey")

		// This should be sanitized by the VMess method
		vmessInbound := config.MakeVmessInbound("test", 8080, "uuid", "auto", grpcSettings)

		// VMess should have sanitized the REALITY settings
		if vmessInbound.StreamSettings.Security == "reality" {
			t.Errorf("VMess should not have REALITY security")
		}
	})

	t.Run("VMess with XHTTP should be sanitized", func(t *testing.T) {
		// Try to create VMess with XHTTP (should be auto-corrected)
		xhttpSettings := config.MakeXhttpStreamSettings("host.com", "/path", "auto")

		// This should be sanitized by the VMess method
		vmessInbound := config.MakeVmessInbound("test", 8080, "uuid", "auto", xhttpSettings)

		// VMess should have sanitized the XHTTP settings
		if vmessInbound.StreamSettings.Network == "xhttp" {
			t.Errorf("VMess should not have XHTTP transport")
		}
	})

	t.Run("Validation functions work correctly", func(t *testing.T) {
		// Valid combination
		wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
		err := config.ValidateProtocolCompatibility("vmess", wsSettings)
		if err != nil {
			t.Errorf("VMess + WebSocket should be valid")
		}

		// Invalid combination
		grpcSettings := config.MakeGrpcStreamSettings("service", "")
		grpcSettings = config.AddRealityToStreamSettings(grpcSettings, "example.com:443",
			[]string{"example.com"}, "privateKey", "publicKey")
		err = config.ValidateProtocolCompatibility("vmess", grpcSettings)
		if err == nil {
			t.Errorf("VMess + REALITY should be invalid")
		}
	})
}
