package state

import (
	"strconv"
	"strings"
	"sync"

	"0dns.io/core/config"
	"0dns.io/core/logging"
	"github.com/0chain/gosdk/core/block"
	"go.uber.org/zap"
)

const (
	HTTPProtocol  = "http://"
	HTTPSProtocol = "https://"
)

// State defines the latest state of network based on most recent magic block.
type State struct {
	sync.RWMutex
	Miners            []string // External URLs (localhost) for serving to clients
	Sharders          []string // External URLs (localhost) for serving to clients
	InternalSharders  []string // Internal URLs (N2NHost) for fetching magic blocks
	CurrentMagicBlock *block.MagicBlock
}

// state is a global status of 0dns.
var state = &State{}

// Get return a copy of state.
func Get() State {
	state.RLock()
	defer state.RUnlock()

	return State{
		Miners:            state.Miners,
		Sharders:          state.Sharders,
		InternalSharders:  state.InternalSharders,
		CurrentMagicBlock: state.CurrentMagicBlock,
	}
}

func SetFromCurrentMagicBlock(c config.Config, b *block.MagicBlock) bool {
	if b == nil {
		panic("Unexpected missing magic block")
	}

	// Only update if new magic block number is higher than current
	state.RLock()
	currentMB := state.CurrentMagicBlock
	state.RUnlock()
	if currentMB != nil && b.MagicBlockNumber <= currentMB.MagicBlockNumber {
		logging.Logger.Debug("ignoring older magic block",
			zap.Int64("current", currentMB.MagicBlockNumber),
			zap.Int64("received", b.MagicBlockNumber))
		return false
	}

	networkProtocol := HTTPProtocol
	if c.UseHTTPS {
		networkProtocol = HTTPSProtocol
	}

	var miners []string
	for _, miner := range b.Miners.Nodes {
		host := miner.Host
		// When use_localhost is enabled, always use 127.0.0.1 for external URLs (local dev with port mappings)
		// Otherwise, replace localhost with N2NHost for production
		if c.UseLocalhost {
			host = "127.0.0.1"
		} else if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
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
	var internalSharders []string
	for _, sharder := range b.Sharders.Nodes {
		// External URL: when use_localhost is enabled, always use 127.0.0.1 (local dev with port mappings)
		host := sharder.Host
		if c.UseLocalhost {
			host = "127.0.0.1"
		} else if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
			host = sharder.N2NHost
		}

		// Internal URL: always use N2NHost for fetching from container
		internalHost := sharder.N2NHost

		if c.UsePath {
			sharders = append(sharders,
				networkProtocol+
					host+
					"/"+
					sharder.Path)
			internalSharders = append(internalSharders,
				networkProtocol+
					internalHost+
					"/"+
					sharder.Path)
		} else {
			sharders = append(sharders,
				networkProtocol+
					host+
					":"+
					strconv.Itoa(sharder.Port))
			internalSharders = append(internalSharders,
				networkProtocol+
					internalHost+
					":"+
					strconv.Itoa(sharder.Port))
		}
	}

	logging.Logger.Info("miners: " + strings.Join(miners, ", "))
	logging.Logger.Info("sharders: " + strings.Join(sharders, ", "))
	logging.Logger.Info("internal sharders: " + strings.Join(internalSharders, ", "))

	state.Lock()
	defer state.Unlock()

	state.CurrentMagicBlock = b
	state.Miners = miners
	state.Sharders = sharders
	state.InternalSharders = internalSharders
	return true
}
