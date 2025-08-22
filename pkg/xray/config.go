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
	Password string `json:"password" validate:"omitempty,min=1,max=64"`
	Method   string `json:"method" validate:"required"`
	Email    string `json:"email" validate:"required"`
	ID       string `json:"id,omitempty"`       // For VMess/VLESS UUID
	AlterId  int    `json:"alterId,omitempty"`  // For VMess
	Level    int    `json:"level,omitempty"`    // User level
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
	Listen   string           `json:"listen" validate:"required"`
	Port     int              `json:"port" validate:"required,min=1,max=65536"`
	Protocol string           `json:"protocol" validate:"required"`
	Settings *InboundSettings `json:"settings" validate:"required"`
	Tag      string           `json:"tag" validate:"required"`
}

type OutboundServer struct {
	Address  string `json:"address" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65536"`
	Method   string `json:"method" validate:"required"`
	Password string `json:"password" validate:"omitempty"`
	Uot      bool   `json:"uot"`
	ID       string `json:"id,omitempty"`      // For VMess/VLESS UUID  
	AlterId  int    `json:"alterId,omitempty"` // For VMess
	Level    int    `json:"level,omitempty"`   // User level
}

type OutboundSettings struct {
	Servers []*OutboundServer `json:"servers,omitempty" validate:"omitempty,dive"`
}

type StreamSettings struct {
	Network string `json:"network" validate:"required"`
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
func (c *Config) MakeVlessInbound(tag string, port int, uuid string, network string, security interface{}) *Inbound {
	settings := &InboundSettings{
		Clients: []*Client{
			{
				ID:       uuid,
				Password: "",
				Method:   "none",
				Email:    "client@example.com",
				Level:    0,
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
	
	return inbound
}

func (c *Config) MakeVlessOutbound(tag, address string, port int, uuid, network string) *Outbound {
	return &Outbound{
		Tag:      tag,
		Protocol: "vless",
		Settings: &OutboundSettings{
			Servers: []*OutboundServer{
				{
					Address:  address,
					Port:     port,
					Method:   "none",
					Password: "",
					ID:       uuid,
					Level:    0,
				},
			},
		},
		StreamSettings: &StreamSettings{
			Network: network,
		},
	}
}

// VMess Protocol Support  
func (c *Config) MakeVmessInbound(tag string, port int, uuid, encryption, network string) *Inbound {
	settings := &InboundSettings{
		Clients: []*Client{
			{
				ID:       uuid,
				Password: "",
				Method:   encryption,
				Email:    "client@example.com",
				AlterId:  0,
				Level:    0,
			},
		},
	}
	
	return &Inbound{
		Tag:      tag,
		Protocol: "vmess",
		Listen:   "0.0.0.0",
		Port:     port,
		Settings: settings,
	}
}

func (c *Config) MakeVmessOutbound(tag, address string, port int, uuid, encryption, network string) *Outbound {
	return &Outbound{
		Tag:      tag,
		Protocol: "vmess",
		Settings: &OutboundSettings{
			Servers: []*OutboundServer{
				{
					Address:  address,
					Port:     port,
					Method:   encryption,
					Password: "",
					ID:       uuid,
					AlterId:  0,
					Level:    0,
				},
			},
		},
		StreamSettings: &StreamSettings{
			Network: network,
		},
	}
}

// Trojan Protocol Support
func (c *Config) MakeTrojanInbound(tag string, port int, password, network string, security interface{}) *Inbound {
	settings := &InboundSettings{
		Clients: []*Client{
			{
				Password: password,
				Method:   "none",
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
					Method:   "none",
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
	return errors.WithStack(validator.New(validator.WithRequiredStructEnabled()).Struct(c))
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
