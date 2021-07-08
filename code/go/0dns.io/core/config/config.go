package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/0chain/gosdk/core/block"
	"github.com/spf13/viper"
)

//SetupDefaultConfig - setup the default config options that can be overridden via the config file
func SetupDefaultConfig() {
	viper.SetDefault("logging.level", "info")
}

/*SetupConfig - setup the configuration system */
func SetupConfig() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigName("0dns")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

const (
	DeploymentDevelopment = 0
	DeploymentTestNet     = 1
	DeploymentMainNet     = 2
	HTTPProtocol          = "http://"
	HTTPSProtocol         = "https://"
)

/*Config - all the config options passed from the command line*/
type Config struct {
	Port            int
	ChainID         string
	DeploymentMode  byte
	SignatureScheme string
	Miners          []string
	Sharders        []string

	MagicBlockWorkerTimerInSeconds int64

	UseHTTPS bool
	UsePath  bool

	CurrentMagicBlock *block.MagicBlock
}

/*Configuration of the system */
var Configuration Config

/*TestNet is the program running in TestNet mode? */
func TestNet() bool {
	return Configuration.DeploymentMode == DeploymentTestNet
}

/*Development - is the programming running in development mode? */
func Development() bool {
	return Configuration.DeploymentMode == DeploymentDevelopment
}

func (c *Config) UpdateMagicBlock(m *block.MagicBlock) {
	c.CurrentMagicBlock = m
}

func (c *Config) SetMinerSharderNodes() {
	if c.CurrentMagicBlock == nil {
		panic("No magic block set")
	}

	networkProtocol := HTTPProtocol
	if c.UseHTTPS {
		networkProtocol = HTTPSProtocol
	}

	var miners []string
	for _, miner := range c.CurrentMagicBlock.Miners.Nodes {
		host := miner.Host
		if strings.Contains(host, "localhost") {
			host = miner.N2NHost
		}

		if c.UsePath {
			miners = append(miners,
				networkProtocol+
					host+
					"/"+
					miner.Path)
		} else {
			miners = append(miners,
				networkProtocol+
					host+
					":"+
					strconv.Itoa(miner.Port))
		}
	}

	var sharders []string
	for _, sharder := range c.CurrentMagicBlock.Sharders.Nodes {
		host := sharder.Host
		if strings.Contains(host, "localhost") {
			host = sharder.N2NHost
		}
		if c.UsePath {
			sharders = append(sharders,
				networkProtocol+
					host+
					"/"+
					sharder.Path)
		} else {
			sharders = append(sharders,
				networkProtocol+
					host+
					":"+
					strconv.Itoa(sharder.Port))
		}
	}

	c.Miners = miners
	c.Sharders = sharders
}
