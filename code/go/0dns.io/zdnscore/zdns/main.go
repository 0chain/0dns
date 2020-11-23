package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"0dns.io/core/common"
	"0dns.io/core/config"
	"0dns.io/core/datastore"
	"0dns.io/core/logging"
	. "0dns.io/core/logging"

	"0dns.io/zdnscore/models"
	"0dns.io/zdnscore/worker"

	"github.com/0chain/gosdk/core/block"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func initializeConfig() {
	config.Configuration.ChainID = viper.GetString("server_chain.id")
	config.Configuration.SignatureScheme = viper.GetString("server_chain.signature_scheme")
	config.Configuration.Port = viper.GetInt("port")

	config.Configuration.MongoURL = viper.GetString("mongo.url")
	config.Configuration.DBName = viper.GetString("mongo.db_name")
	config.Configuration.MongoPoolSize = viper.GetInt64("mongo.pool_size")

	config.Configuration.MagicBlockWorkerTimerInSeconds = viper.GetInt64("worker.magic_block_worker")

	config.Configuration.UseHTTPS = viper.GetBool("use_https")
}

func initializeMagicBlock(magicBlockFile string) {
	magicBlock, err := os.Open(magicBlockFile)
	if err != nil {
		Logger.Error("Failed to read magic block with error", zap.Error(err))
		panic("unable to read magic block file")
	}

	magicBlockBytes, err := ioutil.ReadAll(magicBlock)
	if err != nil {
		Logger.Error("Failed to read magic block with error", zap.Error(err))
		panic("unable to read magic block file")
	}

	var m block.MagicBlock
	err = json.Unmarshal(magicBlockBytes, &m)
	if err != nil {
		Logger.Error("Failed to unmarshal magic block bytes", zap.Error(err))
		panic("Unable to unmarshal magic block bytes")
	}

	if !models.CheckMagicBlockPresentInDB(context.Background(), m.MagicBlockNumber) {
		err = models.InsertMagicBlock(context.Background(), &m)
		if err != nil {
			Logger.Error("Failed to insert magic block to the DB", zap.Error(err))
			panic("Unable to insert magic blockto the DB")
		}
	}

	// fetch old blocks
	if m.MagicBlockNumber != 1 {
		go worker.FetchOldMagicBlocks(context.Background(), m.MagicBlockNumber-1)
	}

	config.Configuration.CurrentMagicBlock = &m
	config.Configuration.SetMinerSharderNodes()
}

func initHandlers(r *mux.Router) {
	r.HandleFunc("/", common.UserRateLimit(HomePageHandler))
	r.HandleFunc("/network", common.UserRateLimit(NetworkDetailsHandler))
	r.HandleFunc("/magic_block", common.UserRateLimit(LatestMagicBlockHandler))
}

var startTime time.Time

func main() {
	deploymentMode := flag.Int("deployment_mode", 2, "deployment_mode")
	magicBlockFile := flag.String("magic_block", "", "magic_block")

	flag.Parse()

	config.Configuration.DeploymentMode = byte(*deploymentMode)
	config.SetupDefaultConfig()
	config.SetupConfig()

	if config.Development() {
		logging.InitLogging("development")
	} else {
		logging.InitLogging("production")
	}
	initializeConfig()

	common.SetupRootContext(context.Background())

	checkForDBConnection(context.Background())

	initializeMagicBlock(*magicBlockFile)

	address := fmt.Sprintf(":%v", config.Configuration.Port)

	var server *http.Server
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})
	rHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)
	if config.Development() {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        rHandler,
		}
	} else {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        rHandler,
		}
	}
	common.HandleShutdown(server)

	common.ConfigRateLimits()
	initHandlers(r)
	go worker.SetupWorkers(context.Background())

	startTime = time.Now().UTC()
	Logger.Info("Ready to listen to the requests on ", zap.Any("port", config.Configuration.Port))
	log.Fatal(server.ListenAndServe())
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div>Running since %v ...\n", startTime)
	fmt.Fprintf(w, "<div>Working on the chain: %v</div>\n", config.Configuration.ChainID)
	fmt.Fprintf(w, "<div>I am 0dns with <ul><li>miners:%v</li><li>sharders:%v</li></ul></div>\n",
		config.Configuration.Miners, config.Configuration.Sharders)
}

func NetworkDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Miners   []string `json:"miners"`
		Sharders []string `json:"sharders"`
	}

	response.Miners = config.Configuration.Miners
	response.Sharders = config.Configuration.Sharders

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LatestMagicBlockHandler(w http.ResponseWriter, r *http.Request) {
	magicBlock := config.Configuration.CurrentMagicBlock

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(magicBlock)
}

func checkForDBConnection(ctx context.Context) {
	retries := 0
	var err error
	for retries < 600 {
		Logger.Info("Trying to connect to mongoDB ...")
		err = datastore.GetStore().Open(ctx)
		if err != nil {
			time.Sleep(1 * time.Second)
			retries++
			continue
		}
		Logger.Info("DB Connection done.")
		break
	}

	if err != nil {
		Logger.Error("Error in opening the database. Shutting the server down")
		panic(err)
	}
}
