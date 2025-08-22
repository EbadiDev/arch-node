package xray

import (
	"encoding/json"
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

type Log struct {
	LogLevel string `json:"loglevel" validate:"required"`
	Access   string `json:"access,omitempty"`
	Error    string `json:"error,omitempty"`
}

type Client struct {
	Password string `json:"password,omitempty" validate:"omitempty,min=1,max=64"`
	Method   string `json:"method,omitempty"`     // Required for Shadowsocks, optional for others
	Email    string `json:"email" validate:"required"`
	ID       string `json:"id,omitempty"`         // For VMess/VLESS UUID
	AlterId  int    `json:"alterId,omitempty"`    // For VMess
	Level    int    `json:"level,omitempty"`      // User level
	Security string `json:"security,omitempty"`   // For VMess/VLESS security
}

type InboundSettings struct {
	Address    string    `json:"address,omitempty"`
	Clients    []*Client `json:"clients,omitempty" validate:"omitempty,dive"`
	Network    string    `json:"network,omitempty"`
	Method     string    `json:"method,omitempty"`
	Password   string    `json:"password,omitempty"`
	Decryption string    `json:"decryption,omitempty"` // For VLESS
}

type Inbound struct {
	Listen         string           `json:"listen" validate:"required"`
	Port           int              `json:"port" validate:"required,min=1,max=65536"`
	Protocol       string           `json:"protocol" validate:"required"`
	Settings       *InboundSettings `json:"settings" validate:"required"`
	StreamSettings *StreamSettings  `json:"streamSettings,omitempty"`
	Tag            string           `json:"tag" validate:"required"`
}

type OutboundServer struct {
	Address  string `json:"address" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65536"`
	Method   string `json:"method,omitempty"`     // Only required for Shadowsocks
	Password string `json:"password,omitempty"`   // For Shadowsocks/Trojan
	Uot      bool   `json:"uot"`
	ID       string `json:"id,omitempty"`         // For VMess/VLESS UUID  
	AlterId  int    `json:"alterId,omitempty"`    // For VMess
	Level    int    `json:"level,omitempty"`      // User level
	Security string `json:"security,omitempty"`   // For VMess/VLESS/Trojan
}

type OutboundSettings struct {
	Servers []*OutboundServer `json:"servers,omitempty" validate:"omitempty,dive"`
	Vnext   []*VnextServer    `json:"vnext,omitempty" validate:"omitempty,dive"`
}

// VnextServer represents a VMess outbound server configuration
type VnextServer struct {
	Address string      `json:"address" validate:"required"`
	Port    int         `json:"port" validate:"required,min=1,max=65536"`
	Users   []*VmessUser `json:"users" validate:"required,dive"`
}

// VmessUser represents a VMess user configuration for outbound connections
type VmessUser struct {
	ID       string `json:"id" validate:"required"`
	AlterId  int    `json:"alterId,omitempty"`
	Level    int    `json:"level,omitempty"`
	Security string `json:"security,omitempty"`
}

type StreamSettings struct {
	Network            string                `json:"network" validate:"required"`
	Security           string                `json:"security,omitempty"`
	
	// Transport-specific settings
	TcpSettings        *TcpSettings         `json:"tcpSettings,omitempty"`
	WsSettings         *WebSocketSettings   `json:"wsSettings,omitempty"`
	HttpSettings       *HttpSettings        `json:"httpSettings,omitempty"`
	GrpcSettings       *GrpcSettings        `json:"grpcSettings,omitempty"`
	KcpSettings        *KcpSettings         `json:"kcpSettings,omitempty"`
	HttpUpgradeSettings *HttpUpgradeSettings `json:"httpupgradeSettings,omitempty"`
	XhttpSettings      *XhttpSettings       `json:"xhttpSettings,omitempty"`
	
	// Security settings
	TlsSettings        *TlsSettings         `json:"tlsSettings,omitempty"`
	RealitySettings    *RealitySettings     `json:"realitySettings,omitempty"`
	
	// Socket settings
	SocketSettings     *SocketSettings      `json:"sockopt,omitempty"`
}

// TCP with HTTP header masquerading
type TcpSettings struct {
	AcceptProxyProtocol bool              `json:"acceptProxyProtocol,omitempty"`
	Header             *TcpHeaderObject   `json:"header,omitempty"`
}

type TcpHeaderObject struct {
	Type     string                    `json:"type"`
	Request  *HttpRequestObject        `json:"request,omitempty"`
	Response *HttpResponseObject       `json:"response,omitempty"`
}

type HttpRequestObject struct {
	Version string              `json:"version,omitempty"`
	Method  string              `json:"method,omitempty"`
	Path    []string            `json:"path,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type HttpResponseObject struct {
	Version string              `json:"version,omitempty"`
	Status  string              `json:"status,omitempty"`
	Reason  string              `json:"reason,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// WebSocket
type WebSocketSettings struct {
	AcceptProxyProtocol bool   `json:"acceptProxyProtocol,omitempty"`
	Path               string `json:"path,omitempty"`
	Host               string `json:"host,omitempty"`
	HeartbeatPeriod    int    `json:"heartbeatPeriod,omitempty"`
	CustomHost         string `json:"custom_host,omitempty"`
}

// HTTP (Legacy HTTP transport)
type HttpSettings struct {
	Host []string `json:"host,omitempty"`
	Path string   `json:"path,omitempty"`
}

// gRPC
type GrpcSettings struct {
	AcceptProxyProtocol   bool   `json:"acceptProxyProtocol,omitempty"`
	ServiceName          string `json:"serviceName,omitempty"`
	Authority            string `json:"authority,omitempty"`
	MultiMode            bool   `json:"multiMode,omitempty"`
	UserAgent            string `json:"user_agent,omitempty"`
	IdleTimeout          int    `json:"idle_timeout,omitempty"`
	HealthCheckTimeout   int    `json:"health_check_timeout,omitempty"`
	PermitWithoutStream  bool   `json:"permit_without_stream,omitempty"`
	InitialWindowsSize   int    `json:"initial_windows_size,omitempty"`
}

// KCP
type KcpSettings struct {
	AcceptProxyProtocol bool              `json:"acceptProxyProtocol,omitempty"`
	Mtu                int               `json:"mtu,omitempty"`
	Tti                int               `json:"tti,omitempty"`
	UplinkCapacity     int               `json:"uplinkCapacity,omitempty"`
	DownlinkCapacity   int               `json:"downlinkCapacity,omitempty"`
	Congestion         bool              `json:"congestion,omitempty"`
	ReadBufferSize     int               `json:"readBufferSize,omitempty"`
	WriteBufferSize    int               `json:"writeBufferSize,omitempty"`
	Header             *KcpHeaderObject  `json:"header,omitempty"`
	Seed               string            `json:"seed,omitempty"`
}

type KcpHeaderObject struct {
	Type   string `json:"type"`
	Domain string `json:"domain,omitempty"`
}

// HTTP Upgrade
type HttpUpgradeSettings struct {
	AcceptProxyProtocol bool   `json:"acceptProxyProtocol,omitempty"`
	Host               string `json:"host,omitempty"`
	Path               string `json:"path,omitempty"`
	CustomHost         string `json:"custom_host,omitempty"`
}

// XHTTP
type XhttpSettings struct {
	AcceptProxyProtocol bool   `json:"acceptProxyProtocol,omitempty"`
	Host               string `json:"host,omitempty"`
	CustomHost         string `json:"custom_host,omitempty"`
	Path               string `json:"path,omitempty"`
	NoSSEHeader        bool   `json:"noSSEHeader,omitempty"`
	NoGRPCHeader       bool   `json:"noGRPCHeader,omitempty"`
	Mode               string `json:"mode,omitempty"`
}

// Socket Settings (common to all transports)
type SocketSettings struct {
	UseSocket           bool   `json:"useSocket,omitempty"`
	DomainStrategy      string `json:"DomainStrategy,omitempty"`
	TcpKeepAliveInterval int    `json:"tcpKeepAliveInterval,omitempty"`
	TcpUserTimeout      int    `json:"tcpUserTimeout,omitempty"`
	TcpMaxSeg           int    `json:"tcpMaxSeg,omitempty"`
	TcpWindowClamp      int    `json:"tcpWindowClamp,omitempty"`
	TcpKeepAliveIdle    int    `json:"tcpKeepAliveIdle,omitempty"`
	TcpMptcp            bool   `json:"tcpMptcp,omitempty"`
}

// TLS Settings
type TlsSettings struct {
	ServerName          string   `json:"serverName,omitempty"`
	RejectUnknownSni    bool     `json:"rejectUnknownSni,omitempty"`
	AllowInsecure       bool     `json:"allowInsecure,omitempty"`
	Fingerprint         string   `json:"fingerprint,omitempty"`
	Sni                 string   `json:"sni,omitempty"`
	CurvePreferences    string   `json:"curvepreferences,omitempty"`
	Alpn                []string `json:"alpn,omitempty"`
	ServerNameToVerify  string   `json:"serverNameToVerify,omitempty"`
}

// REALITY Settings
type RealitySettings struct {
	Show         bool     `json:"show,omitempty"`
	Dest         string   `json:"dest,omitempty"`
	PrivateKey   string   `json:"privatekey,omitempty"`
	MinClientVer string   `json:"minclientver,omitempty"`
	MaxClientVer string   `json:"maxclientver,omitempty"`
	MaxTimeDiff  int      `json:"maxtimediff,omitempty"`
	ProxyProtocol int     `json:"proxyprotocol,omitempty"`
	ShortIds     []string `json:"shortids,omitempty"`
	ServerNames  []string `json:"serverNames,omitempty"`
	Fingerprint  string   `json:"fingerprint,omitempty"`
	SpiderX      string   `json:"spiderx,omitempty"`
	PublicKey    string   `json:"publickey,omitempty"`
}

type Outbound struct {
	Protocol       string            `json:"protocol" validate:"required"`
	Tag            string            `json:"tag" validate:"required"`
	Settings       *OutboundSettings `json:"settings,omitempty"`
	StreamSettings *StreamSettings   `json:"streamSettings,omitempty"`
}

type DNS struct {
	Servers []string `json:"servers" validate:"required"`
}

type API struct {
	Tag      string   `json:"tag" validate:"required"`
	Services []string `json:"services" validate:"required"`
}

type PolicyLevels struct {
	StatsUserUplink   bool `json:"statsUserUplink"`
	StatsUserDownlink bool `json:"statsUserDownlink"`
}

type Policy struct {
	Levels map[string]map[string]bool `json:"levels"`
	System map[string]bool            `json:"system"`
}

type Rule struct {
	InboundTag  []string `json:"inboundTag" validate:"required"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	BalancerTag string   `json:"balancerTag,omitempty"`
	Domain      []string `json:"domain,omitempty"`
}

type RoutingSettings struct {
	Rules []*Rule `json:"rules" validate:"required,dive"`
}

type Balancer struct {
	Tag      string   `json:"tag" validate:"required"`
	Selector []string `json:"selector"`
}

type Routing struct {
	DomainStrategy string      `json:"domainStrategy" validate:"required"`
	DomainMatcher  string      `json:"domainMatcher" validate:"required"`
	Rules          []*Rule     `json:"rules,omitempty" validate:"omitempty,dive"`
	Balancers      []*Balancer `json:"balancers,omitempty" validate:"omitempty,dive"`
}

type Reverse struct {
	Bridges []*ReverseItem `json:"bridges,omitempty"  validate:"omitempty,dive"`
	Portals []*ReverseItem `json:"portals,omitempty"  validate:"omitempty,dive"`
}

type ReverseItem struct {
	Tag    string `json:"tag"  validate:"required"`
	Domain string `json:"domain"  validate:"required"`
}

type Metadata struct {
	UpdatedAt string `json:"updatedAt"`
	UpdatedBy string `json:"UpdatedBy"`
}

type Config struct {
	Log       *Log                   `json:"log" validate:"required"`
	Inbounds  []*Inbound             `json:"inbounds" validate:"required,dive"`
	Outbounds []*Outbound            `json:"outbounds" validate:"required,dive"`
	DNS       *DNS                   `json:"dns" validate:"required"`
	Stats     map[string]interface{} `json:"stats" validate:"required"`
	API       *API                   `json:"api" validate:"required"`
	Policy    *Policy                `json:"policy" validate:"required"`
	Routing   *Routing               `json:"routing" validate:"required"`
	Reverse   *Reverse               `json:"reverse,omitempty"`
	Metadata  *Metadata              `json:"_metadata,omitempty"`
}

func (c *Config) MakeShadowsocksInbound(tag, password, method, network string, port int, clients []*Client) *Inbound {
	return &Inbound{
		Tag:      tag,
		Protocol: "shadowsocks",
		Listen:   "0.0.0.0",
		Port:     port,
		Settings: &InboundSettings{
			Clients:  clients,
			Password: password,
			Method:   method,
			Network:  network,
		},
	}
}

func (c *Config) MakeShadowsocksOutbound(tag, host, password, method string, port int) *Outbound {
	return &Outbound{
		Tag:      tag,
		Protocol: "shadowsocks",
		Settings: &OutboundSettings{
			Servers: []*OutboundServer{
				{
					Address:  host,
					Port:     port,
					Method:   method,
					Password: password,
					Uot:      true,
				},
			},
		},
		StreamSettings: &StreamSettings{
			Network: "tcp",
		},
	}
}

// VLESS Protocol Support
func (c *Config) MakeVlessInbound(tag string, port int, uuid string, network string, streamSettings *StreamSettings) *Inbound {
	// Validate compatibility for VLESS protocol
	if streamSettings != nil {
		// VLESS supports all transport types including XHTTP
		if streamSettings.XhttpSettings != nil {
			// VLESS + XHTTP is valid
		}
		
		// VLESS supports both TLS and REALITY
		if streamSettings.Security == "reality" && streamSettings.RealitySettings != nil {
			// VLESS + REALITY is valid
		} else if streamSettings.Security == "tls" && streamSettings.TlsSettings != nil {
			// VLESS + TLS is valid
		}
	}

	settings := &InboundSettings{
		Clients: []*Client{
			{
				ID:    uuid,
				Email: "client@example.com",
				Level: 0,
			},
		},
		Decryption: "none",
	}
	
	inbound := &Inbound{
		Tag:      tag,
		Protocol: "vless",
		Listen:   "0.0.0.0",
		Port:     port,
		Settings: settings,
	}
	
	// Apply stream settings if provided
	if streamSettings != nil {
		inbound.StreamSettings = streamSettings
	} else if network != "" && network != "tcp" {
		// Fallback for backward compatibility
		inbound.StreamSettings = &StreamSettings{Network: network}
	}
	
	return inbound
}

func (c *Config) MakeVlessOutbound(tag, address string, port int, uuid, network string) *Outbound {
	return &Outbound{
		Tag:      tag,
		Protocol: "vless",
		Settings: &OutboundSettings{
			Servers: []*OutboundServer{
				{
					Address: address,
					Port:    port,
					ID:      uuid,
					Level:   0,
				},
			},
		},
		StreamSettings: &StreamSettings{
			Network: network,
		},
	}
}

// VMess Protocol Support  
func (c *Config) MakeVmessInbound(tag string, port int, uuid, encryption string, streamSettings *StreamSettings) *Inbound {
	// Validate compatibility for VMess protocol
	if streamSettings != nil {
		// VMess does NOT support REALITY (VLESS-only feature)
		if streamSettings.Security == "reality" || streamSettings.RealitySettings != nil {
			// This is invalid - VMess doesn't support REALITY
			// Could return error or fallback to TLS
			streamSettings.Security = "tls"
			streamSettings.RealitySettings = nil
		}
		
		// VMess does NOT support XHTTP (VLESS-only feature)  
		if streamSettings.XhttpSettings != nil {
			// This is invalid - VMess doesn't support XHTTP
			// Could return error or fallback to WebSocket
			streamSettings.XhttpSettings = nil
			if streamSettings.Network == "xhttp" {
				streamSettings.Network = "ws"
			}
		}
	}

	settings := &InboundSettings{
		Clients: []*Client{
			{
				ID:       uuid,
				Email:    "client@example.com",
				AlterId:  0,
				Level:    0,
				Security: encryption, // Use Security field for VMess encryption
			},
		},
	}
	
	// If no stream settings provided, default to TCP
	if streamSettings == nil {
		streamSettings = &StreamSettings{
			Network: "tcp",
		}
	}
	
	return &Inbound{
		Tag:            tag,
		Protocol:       "vmess",
		Listen:         "0.0.0.0",
		Port:           port,
		Settings:       settings,
		StreamSettings: streamSettings,
	}
}

func (c *Config) MakeVmessOutbound(tag, address string, port int, uuid, encryption string, streamSettings *StreamSettings) *Outbound {
	// If no stream settings provided, default to TCP
	if streamSettings == nil {
		streamSettings = &StreamSettings{
			Network: "tcp",
		}
	}

	return &Outbound{
		Tag:      tag,
		Protocol: "vmess",
		Settings: &OutboundSettings{
			Vnext: []*VnextServer{
				{
					Address: address,
					Port:    port,
					Users: []*VmessUser{
						{
							ID:       uuid,
							AlterId:  0,
							Level:    0,
							Security: encryption,
						},
					},
				},
			},
		},
		StreamSettings: streamSettings,
	}
}

// Trojan Protocol Support
func (c *Config) MakeTrojanInbound(tag string, port int, password, network string, security interface{}) *Inbound {
	settings := &InboundSettings{
		Clients: []*Client{
			{
				Password: password,
				Email:    "client@example.com",
			},
		},
	}
	
	return &Inbound{
		Tag:      tag,
		Protocol: "trojan",
		Listen:   "0.0.0.0",
		Port:     port,
		Settings: settings,
	}
}

func (c *Config) MakeTrojanOutbound(tag, address string, port int, password, network string) *Outbound {
	return &Outbound{
		Tag:      tag,
		Protocol: "trojan",
		Settings: &OutboundSettings{
			Servers: []*OutboundServer{
				{
					Address:  address,
					Port:     port,
					Password: password,
				},
			},
		},
		StreamSettings: &StreamSettings{
			Network: network,
		},
	}
}

func (c *Config) FindInbound(tag string) *Inbound {
	for _, inbound := range c.Inbounds {
		if inbound.Tag == tag {
			return inbound
		}
	}
	return nil
}

func (c *Config) FindOutbound(tag string) *Outbound {
	for _, outbound := range c.Outbounds {
		if outbound.Tag == tag {
			return outbound
		}
	}
	return nil
}

func (c *Config) FindBalancer(tag string) *Balancer {
	for _, balancer := range c.Routing.Balancers {
		if balancer.Tag == tag {
			return balancer
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if c.FindInbound("api") == nil {
		return errors.New("xray: config: api inbound not found")
	}
	
	// Protocol-aware validation
	if err := c.validateProtocolSpecific(); err != nil {
		return errors.WithStack(err)
	}
	
	return errors.WithStack(validator.New(validator.WithRequiredStructEnabled()).Struct(c))
}

// validateProtocolSpecific performs protocol-specific validation
func (c *Config) validateProtocolSpecific() error {
	// Validate inbounds
	for _, inbound := range c.Inbounds {
		if inbound.Settings != nil && inbound.Settings.Clients != nil {
			for _, client := range inbound.Settings.Clients {
				if err := c.validateClient(client, inbound.Protocol); err != nil {
					return errors.Wrapf(err, "invalid client for protocol %s in inbound %s", inbound.Protocol, inbound.Tag)
				}
			}
		}
	}
	
	// Validate outbounds
	for _, outbound := range c.Outbounds {
		if outbound.Settings != nil {
			// Validate servers (for Shadowsocks, VLESS, Trojan)
			if outbound.Settings.Servers != nil {
				for _, server := range outbound.Settings.Servers {
					if err := c.validateServer(server, outbound.Protocol); err != nil {
						return errors.Wrapf(err, "invalid server for protocol %s in outbound %s", outbound.Protocol, outbound.Tag)
					}
				}
			}
			
			// Validate vnext (for VMess)
			if outbound.Settings.Vnext != nil {
				for _, vnext := range outbound.Settings.Vnext {
					if err := c.validateVnext(vnext, outbound.Protocol); err != nil {
						return errors.Wrapf(err, "invalid vnext for protocol %s in outbound %s", outbound.Protocol, outbound.Tag)
					}
				}
			}
		}
	}
	
	return nil
}

// validateClient validates client configuration based on protocol
func (c *Config) validateClient(client *Client, protocol string) error {
	switch protocol {
	case "shadowsocks":
		if client.Method == "" {
			return errors.New("shadowsocks client requires method field")
		}
		if client.Password == "" {
			return errors.New("shadowsocks client requires password field")
		}
	case "vmess":
		if client.ID == "" {
			return errors.New("vmess client requires id (UUID) field")
		}
	case "vless":
		if client.ID == "" {
			return errors.New("vless client requires id (UUID) field")
		}
	case "trojan":
		if client.Password == "" {
			return errors.New("trojan client requires password field")
		}
	case "dokodemo-door", "freedom", "blackhole", "socks", "http":
		// System protocols - no client validation needed
		break
	default:
		// Allow unknown protocols to pass validation
		break
	}
	return nil
}

// validateServer validates server configuration based on protocol
func (c *Config) validateServer(server *OutboundServer, protocol string) error {
	switch protocol {
	case "shadowsocks":
		if server.Method == "" {
			return errors.New("shadowsocks server requires method field")
		}
		if server.Password == "" {
			return errors.New("shadowsocks server requires password field")
		}
	case "vmess":
		if server.ID == "" {
			return errors.New("vmess server requires id (UUID) field")
		}
	case "vless":
		if server.ID == "" {
			return errors.New("vless server requires id (UUID) field")
		}
	case "trojan":
		if server.Password == "" {
			return errors.New("trojan server requires password field")
		}
	case "freedom", "blackhole", "socks", "http":
		// System protocols - no server validation needed
		break
	default:
		// Allow unknown protocols to pass validation
		break
	}
	return nil
}

// validateVnext validates VMess vnext configuration
func (c *Config) validateVnext(vnext *VnextServer, protocol string) error {
	if protocol == "vmess" {
		if len(vnext.Users) == 0 {
			return errors.New("vmess vnext requires at least one user")
		}
		for _, user := range vnext.Users {
			if user.ID == "" {
				return errors.New("vmess user requires id (UUID) field")
			}
		}
	}
	return nil
}

func (c *Config) Equals(other *Config) bool {
	json1, err := json.Marshal(c)
	if err != nil {
		return false
	}

	json2, err := json.Marshal(other)
	if err != nil {
		return false
	}

	return string(json1) == string(json2)
}

func NewConfig(logLevel string) *Config {
	return &Config{
		Log: &Log{
			LogLevel: logLevel,
			Access:   "./storage/logs/xray-access.log",
			Error:    "./storage/logs/xray-error.log",
		},
		Inbounds: []*Inbound{
			{
				Tag:      "api",
				Protocol: "dokodemo-door",
				Listen:   "127.0.0.1",
				Port:     3411,
				Settings: &InboundSettings{
					Address: "127.0.0.1",
					Network: "tcp",
				},
			},
		},
		Outbounds: []*Outbound{
			{
				Tag:      "out",
				Protocol: "freedom",
			},
		},
		DNS: &DNS{
			Servers: []string{"8.8.8.8", "8.8.4.4", "localhost"},
		},
		Stats: map[string]interface{}{},
		API: &API{
			Tag:      "api",
			Services: []string{"StatsService"},
		},
		Policy: &Policy{
			Levels: map[string]map[string]bool{
				"0": {
					"statsUserUplink":   true,
					"statsUserDownlink": true,
				},
			},
			System: map[string]bool{
				"statsInboundUplink":    true,
				"statsInboundDownlink":  true,
				"statsOutboundUplink":   true,
				"statsOutboundDownlink": true,
			},
		},
		Routing: &Routing{
			DomainStrategy: "AsIs",
			DomainMatcher:  "hybrid",
			Rules: []*Rule{
				{
					InboundTag:  []string{"api"},
					OutboundTag: "api",
				},
			},
			Balancers: []*Balancer{},
		},
		Reverse: &Reverse{
			Bridges: []*ReverseItem{},
			Portals: []*ReverseItem{},
		},
	}
}

// Transport Helper Methods

// MakeTcpStreamSettings creates TCP stream settings with optional HTTP header masquerading
func (c *Config) MakeTcpStreamSettings(httpHeader bool) *StreamSettings {
	settings := &StreamSettings{
		Network: "tcp",
	}
	
	if httpHeader {
		settings.TcpSettings = &TcpSettings{
			Header: &TcpHeaderObject{
				Type: "http",
				Request: &HttpRequestObject{
					Version: "1.1",
					Method:  "GET",
					Path:    []string{"/"},
					Headers: map[string][]string{
						"Host":       {"www.example.com"},
						"User-Agent": {"Mozilla/5.0"},
					},
				},
				Response: &HttpResponseObject{
					Version: "1.1",
					Status:  "200",
					Reason:  "OK",
					Headers: map[string][]string{
						"Content-Type": {"text/html"},
					},
				},
			},
		}
	}
	
	return settings
}

// MakeWebSocketStreamSettings creates WebSocket stream settings
func (c *Config) MakeWebSocketStreamSettings(path, host string) *StreamSettings {
	return &StreamSettings{
		Network: "ws",
		WsSettings: &WebSocketSettings{
			Path: path,
			Host: host,
		},
	}
}

// MakeGrpcStreamSettings creates gRPC stream settings
func (c *Config) MakeGrpcStreamSettings(serviceName, authority string) *StreamSettings {
	return &StreamSettings{
		Network: "grpc",
		GrpcSettings: &GrpcSettings{
			ServiceName: serviceName,
			Authority:   authority,
		},
	}
}

// MakeKcpStreamSettings creates KCP stream settings
func (c *Config) MakeKcpStreamSettings(headerType, seed string) *StreamSettings {
	settings := &StreamSettings{
		Network: "kcp",
		KcpSettings: &KcpSettings{
			Mtu:              1350,
			Tti:              50,
			UplinkCapacity:   5,
			DownlinkCapacity: 20,
			Congestion:       false,
			ReadBufferSize:   2,
			WriteBufferSize:  2,
			Seed:             seed,
		},
	}
	
	if headerType != "" {
		settings.KcpSettings.Header = &KcpHeaderObject{
			Type: headerType,
		}
	}
	
	return settings
}

// MakeHttpUpgradeStreamSettings creates HTTP Upgrade stream settings
func (c *Config) MakeHttpUpgradeStreamSettings(host, path string) *StreamSettings {
	return &StreamSettings{
		Network: "httpupgrade",
		HttpUpgradeSettings: &HttpUpgradeSettings{
			Host: host,
			Path: path,
		},
	}
}

// MakeXhttpStreamSettings creates XHTTP stream settings
func (c *Config) MakeXhttpStreamSettings(host, path, mode string) *StreamSettings {
	return &StreamSettings{
		Network: "xhttp",
		XhttpSettings: &XhttpSettings{
			Host: host,
			Path: path,
			Mode: mode,
		},
	}
}

// Security Helper Methods

// AddTlsToStreamSettings adds TLS security to existing stream settings
func (c *Config) AddTlsToStreamSettings(streamSettings *StreamSettings, serverName string, allowInsecure bool) *StreamSettings {
	if streamSettings == nil {
		streamSettings = &StreamSettings{Network: "tcp"}
	}
	
	streamSettings.Security = "tls"
	streamSettings.TlsSettings = &TlsSettings{
		ServerName:    serverName,
		AllowInsecure: allowInsecure,
		Alpn:          []string{"h2", "http/1.1"},
	}
	
	return streamSettings
}

// AddRealityToStreamSettings adds REALITY security to existing stream settings
func (c *Config) AddRealityToStreamSettings(streamSettings *StreamSettings, dest string, serverNames []string, privateKey, publicKey string) *StreamSettings {
	if streamSettings == nil {
		streamSettings = &StreamSettings{Network: "tcp"}
	}
	
	streamSettings.Security = "reality"
	streamSettings.RealitySettings = &RealitySettings{
		Dest:        dest,
		ServerNames: serverNames,
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		ShortIds:    []string{"", "0123456789abcdef"},
	}
	
	return streamSettings
}

// Protocol Helper Methods for Transport Configuration
// These methods provide a composable way to build transport configurations

// Example usage:
// wsSettings := config.MakeWebSocketStreamSettings("/path", "host.com")
// wsSettings = config.AddTlsToStreamSettings(wsSettings, "server.com", false) 
// inbound := config.MakeVmessInbound("tag", 8080, "uuid", "auto", wsSettings)

// Note: VMess supports TCP, WebSocket, gRPC, KCP, HTTP Upgrade + TLS
//       VMess does NOT support REALITY or XHTTP (VLESS-only features)
//       VLESS supports all transports and security types

// Protocol Compatibility Validation

// ValidateProtocolCompatibility checks if the given protocol supports the transport and security configuration
func (c *Config) ValidateProtocolCompatibility(protocol string, streamSettings *StreamSettings) error {
	if streamSettings == nil {
		return nil // Basic TCP is supported by all protocols
	}
	
	switch protocol {
	case "vmess":
		// VMess restrictions
		if streamSettings.Security == "reality" || streamSettings.RealitySettings != nil {
			return errors.New("VMess protocol does not support REALITY security (use VLESS instead)")
		}
		if streamSettings.Network == "xhttp" || streamSettings.XhttpSettings != nil {
			return errors.New("VMess protocol does not support XHTTP transport (use VLESS instead)")
		}
		// VMess supports: TCP, WebSocket, gRPC, KCP, HTTP Upgrade + TLS
		
	case "vless":
		// VLESS supports all transports and security types
		// No restrictions needed
		
	case "trojan":
		// Trojan restrictions (usually only supports TCP + TLS)
		if streamSettings.Security != "tls" && streamSettings.Security != "" {
			return errors.New("Trojan protocol typically only supports TLS security")
		}
		if streamSettings.Network != "tcp" && streamSettings.Network != "" {
			return errors.New("Trojan protocol typically only supports TCP transport")
		}
		
	case "shadowsocks":
		// Shadowsocks restrictions (usually basic transports only)
		if streamSettings.Security == "reality" || streamSettings.RealitySettings != nil {
			return errors.New("Shadowsocks protocol does not support REALITY security")
		}
		if streamSettings.Network == "xhttp" || streamSettings.XhttpSettings != nil {
			return errors.New("Shadowsocks protocol does not support XHTTP transport")
		}
	}
	
	return nil
}

// SanitizeStreamSettingsForProtocol removes incompatible settings and returns a safe configuration
func (c *Config) SanitizeStreamSettingsForProtocol(protocol string, streamSettings *StreamSettings) *StreamSettings {
	if streamSettings == nil {
		return nil
	}
	
	// Create a copy to avoid modifying the original
	sanitized := &StreamSettings{
		Network:             streamSettings.Network,
		Security:            streamSettings.Security,
		TcpSettings:         streamSettings.TcpSettings,
		WsSettings:          streamSettings.WsSettings,
		HttpSettings:        streamSettings.HttpSettings,
		GrpcSettings:        streamSettings.GrpcSettings,
		KcpSettings:         streamSettings.KcpSettings,
		HttpUpgradeSettings: streamSettings.HttpUpgradeSettings,
		XhttpSettings:       streamSettings.XhttpSettings,
		TlsSettings:         streamSettings.TlsSettings,
		RealitySettings:     streamSettings.RealitySettings,
		SocketSettings:      streamSettings.SocketSettings,
	}
	
	switch protocol {
	case "vmess":
		// Remove REALITY (VMess doesn't support it)
		if sanitized.Security == "reality" {
			sanitized.Security = "tls" // fallback to TLS
		}
		sanitized.RealitySettings = nil
		
		// Remove XHTTP (VMess doesn't support it)
		if sanitized.Network == "xhttp" {
			sanitized.Network = "ws" // fallback to WebSocket
		}
		sanitized.XhttpSettings = nil
		
	case "trojan":
		// Ensure TLS for Trojan
		if sanitized.Security == "reality" {
			sanitized.Security = "tls"
		}
		sanitized.RealitySettings = nil
		
		// Ensure TCP for Trojan
		if sanitized.Network != "tcp" && sanitized.Network != "" {
			sanitized.Network = "tcp"
		}
		
	case "shadowsocks":
		// Remove advanced features not supported by Shadowsocks
		if sanitized.Security == "reality" {
			sanitized.Security = ""
		}
		sanitized.RealitySettings = nil
		
		if sanitized.Network == "xhttp" {
			sanitized.Network = "tcp"
		}
		sanitized.XhttpSettings = nil
	}
	
	return sanitized
}
