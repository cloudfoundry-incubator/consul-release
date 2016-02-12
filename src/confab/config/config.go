package config

import "encoding/json"

type Config struct {
	Node   ConfigNode
	Confab ConfigConfab
	Consul ConfigConsul
	Path   ConfigPath
}

type ConfigConfab struct {
	TimeoutInSeconds int `json:"timeout_in_seconds"`
}

type ConfigConsul struct {
	Agent       ConfigConsulAgent
	RequireSSL  bool     `json:"require_ssl"`
	EncryptKeys []string `json:"encrypt_keys"`
}

type ConfigPath struct {
	AgentPath       string `json:"agent_path"`
	ConsulConfigDir string `json:"consul_config_dir"`
	PIDFile         string `json:"pid_file"`
}

type ConfigNode struct {
	Name       string
	Index      int
	ExternalIP string `json:"external_ip"`
}

type ConfigConsulAgent struct {
	Servers         ConfigConsulAgentServers
	Services        map[string]ServiceDefinition
	Mode            string
	Datacenter      string `json:"datacenter"`
	Domain          string `json:"domain"`
	Ports           ConfigConsulAgentPorts
	LogLevel        string `json:"log_level"`
	ProtocolVersion int    `json:"protocol_version"`
}

type ConfigConsulAgentServers struct {
	LAN []string
}

type ConfigConsulAgentPorts struct {
	DNS int `json:"dns"`
}

func Default() Config {
	return Config{
		Path: ConfigPath{
			AgentPath:       "/var/vcap/packages/consul/bin/consul",
			ConsulConfigDir: "/var/vcap/jobs/consul_agent/config",
			PIDFile:         "/var/vcap/sys/run/consul_agent/consul_agent.pid",
		},
		Consul: ConfigConsul{
			RequireSSL: true,
			Agent: ConfigConsulAgent{
				Servers: ConfigConsulAgentServers{
					LAN: []string{},
				},
			},
		},
		Confab: ConfigConfab{
			TimeoutInSeconds: 55,
		},
	}
}

func ConfigFromJSON(configData []byte) (Config, error) {
	config := Default()
	if err := json.Unmarshal(configData, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
