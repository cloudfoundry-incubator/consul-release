package chaperon_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/consul-release/src/confab/chaperon"
	"github.com/cloudfoundry-incubator/consul-release/src/confab/config"
	"github.com/cloudfoundry-incubator/consul-release/src/confab/fakes"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("ConfigWriter", func() {
	var (
		configDir string
		cfg       config.Config
		writer    chaperon.ConfigWriter
		logger    *fakes.Logger
	)

	Describe("Write", func() {

		BeforeEach(func() {
			logger = &fakes.Logger{}

			var err error
			configDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			cfg = config.Default()
			cfg.Node = config.ConfigNode{Name: "node", Index: 0}
			cfg.Path.ConsulConfigDir = configDir

			writer = chaperon.NewConfigWriter(configDir, logger)
		})

		It("writes a config file to the consul_config dir", func() {
			err := writer.Write(cfg)
			Expect(err).NotTo(HaveOccurred())

			buf, err := ioutil.ReadFile(filepath.Join(configDir, "config.json"))
			Expect(err).NotTo(HaveOccurred())

			conf := map[string]interface{}{
				"server":     false,
				"domain":     "",
				"datacenter": "",
				"data_dir":   "/var/vcap/store/consul_agent",
				"log_level":  "",
				"node_name":  "node-0",
				"ports": map[string]interface{}{
					"dns": 53,
				},
				"rejoin_after_leave":     true,
				"retry_join":             []string{},
				"retry_join_wan":         []string{},
				"bind_addr":              "",
				"disable_remote_exec":    true,
				"disable_update_check":   true,
				"protocol":               0,
				"verify_outgoing":        true,
				"verify_incoming":        true,
				"verify_server_hostname": true,
				"ca_file":                filepath.Join(configDir, "certs", "ca.crt"),
				"key_file":               filepath.Join(configDir, "certs", "agent.key"),
				"cert_file":              filepath.Join(configDir, "certs", "agent.crt"),
			}
			body, err := json.Marshal(conf)
			Expect(err).To(BeNil())
			Expect(buf).To(MatchJSON(body))

			Expect(logger.Messages).To(ContainSequence([]fakes.LoggerMessage{
				{
					Action: "config-writer.write.generate-configuration",
				},
				{
					Action: "config-writer.write.write-file",
					Data: []lager.Data{{
						"config": config.GenerateConfiguration(cfg, configDir),
					}},
				},
				{
					Action: "config-writer.write.success",
				},
			}))
		})

		Context("failure cases", func() {
			It("returns an error when the config file can't be written to", func() {
				configFile := filepath.Join(configDir, "config.json")
				Expect(os.Mkdir(configFile, os.ModeDir)).To(Succeed())

				err := writer.Write(cfg)
				Expect(err).To(MatchError(ContainSubstring("is a directory")))

				Expect(logger.Messages).To(ContainSequence([]fakes.LoggerMessage{
					{
						Action: "config-writer.write.generate-configuration",
					},
					{
						Action: "config-writer.write.write-file",
						Data: []lager.Data{{
							"config": config.GenerateConfiguration(cfg, configDir),
						}},
					},
					{
						Action: "config-writer.write.write-file.failed",
						Error:  fmt.Errorf("open %s: is a directory", filepath.Join(configDir, "config.json")),
					},
				}))
			})
		})
	})
})
