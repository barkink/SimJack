package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"simjack/config"
	"simjack/engine"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to simulation config JSON file (default: config.json)")
	configJSON := flag.String("config_json", "", "Inline JSON for simulation config")
	logPath := flag.String("log", "output.csv", "Path to log output CSV file (default: output.csv)")
	strategyDir := flag.String("strategies", "strategies", "Directory containing strategy JSON files (default: ./strategies)")
	help := flag.Bool("help", false, "Show usage")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	var cfg config.SimulationConfig
	var err error

	if *configJSON != "" {
		err = json.Unmarshal([]byte(*configJSON), &cfg)
		if err != nil {
			fmt.Printf("Failed to parse inline config: %v\n\n", err)
			printUsage()
			os.Exit(1)
		}
	} else if *configPath != "" {
		file, err := os.Open(*configPath)
		if err != nil {
			fmt.Printf("Failed to open config file: %v\n\n", err)
			printUsage()
			os.Exit(1)
		}
		defer file.Close()
		err = json.NewDecoder(file).Decode(&cfg)
		if err != nil {
			fmt.Printf("Failed to parse config file: %v\n\n", err)
			printUsage()
			os.Exit(1)
		}
	} else {
		fmt.Println("Either -config or -config_json must be provided.")
		printUsage()
		os.Exit(1)
	}

	cfg.StrategyDirectory = *strategyDir

	if err := engine.SetStrategyDirectory(cfg.StrategyDirectory); err != nil {
		fmt.Printf("Invalid strategy directory: %v\n\n", err)
		printUsage()
		os.Exit(1)
	}

	logger, err := engine.NewLogger(*logPath)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n\n", err)
		printUsage()
		os.Exit(1)
	}
	defer logger.Close()

	eng := engine.NewEngine(cfg, logger)
	eng.Run()

	fmt.Println("Simulation completed. Log written to", *logPath)
}

func printUsage() {
	fmt.Println("SimJack - Blackjack Simulation")
	fmt.Println("Usage:")
	fmt.Println("  -config string        Path to simulation config JSON file")
	fmt.Println("  -config_json string   Inline JSON config instead of file")
	fmt.Println("  -log string           Path to log output CSV file (default: output.csv)")
	fmt.Println("  -strategies string    Directory containing strategy JSON files (default: strategies/)")
	fmt.Println("  -help                 Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  simjack -config=test_config.json -log=simjack_log.csv -strategies=.")
	fmt.Println("  simjack -config_json='{\"round_count\":1000,...}' -log=simjack_log.csv")
}
