package main

import (
	steno "github.com/cloudfoundry/gosteno"
	"dea_logging_agent"
	"flag"
	"os"
	"strings"
)

var instancesJsonFilePath = flag.String("instancesFile", "/var/vcap/data/dea_next/db/instances.json", "The DEA instances JSON file")
var logFilePath = flag.String("logFile", "", "The agent log file, defaults to STDOUT")
var loggregatorAddress = flag.String("server", "localhost:3456", "The loggregator TCP host:port for log forwarding")
var logLevel = flag.Bool("v", false, "Verbose logging")

func main() {
	flag.Parse()

	level := steno.LOG_INFO

	if *logLevel {
		level = steno.LOG_DEBUG
	}

	loggingConfig := &steno.Config{
		Sinks: make([]steno.Sink, 1),
		Level:     level,
		Codec:     steno.NewJsonCodec(),
		EnableLOC: true}
	if strings.TrimSpace(*logFilePath) == "" {
		loggingConfig.Sinks[0] = steno.NewIOSink(os.Stdout)
	} else {
		loggingConfig.Sinks[0] = steno.NewFileSink(*logFilePath)
	}
	steno.Init(loggingConfig)
	logger := steno.NewLogger("dea_logging_agent")

	config := &dea_logging_agent.Config{
		InstancesJsonFilePath: *instancesJsonFilePath,
		LoggregatorAddress: *loggregatorAddress}

	dea_logging_agent.Start(config, logger)
}
