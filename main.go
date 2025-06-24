package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"simjack/config"
	"simjack/engine"
)

type CombinedInput struct {
	Config     config.SimulationConfig                 `json:"config"`
	Strategies map[string]engine.CountingStrategyFile `json:"strategies"`
}

func flagIsPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	configPath := flag.String("config", "config.json", "Path to simulation config JSON file (default: config.json)")
	configJSON := flag.String("config_json", "", "Inline JSON for simulation config")
	logPath := flag.String("log", "output.csv", "Path to log output CSV file (default: output.csv)")
	strategyDir := flag.String("strategies", "", "Directory containing strategy JSON files")
	useStdinStrategies := flag.Bool("use-stdin-strategies", false, "Load all strategies from stdin as JSON map")
	useStdinCombined := flag.Bool("use-stdin-combined", false, "Load config + strategies from single JSON on stdin")
	showProgress := flag.Bool("progress", false, "Show progress bar during simulation")
	debug := flag.Bool("debug", false, "Enable debug mode for round-level output")
	help := flag.Bool("help", false, "Show usage")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	if *useStdinCombined && *useStdinStrategies {
		fmt.Println("Cannot use both --use-stdin-combined and --use-stdin-strategies.")
		os.Exit(1)
	}

	var cfg config.SimulationConfig
	var strategyBundle map[string]engine.CountingStrategyFile
	var err error

	// 👇 Öncelik: --use-stdin-combined
	if *useStdinCombined {
		var combined CombinedInput
		err = json.NewDecoder(os.Stdin).Decode(&combined)
		if err != nil {
			fmt.Println("Failed to parse combined stdin input:", err)
			os.Exit(1)
		}
		cfg = combined.Config
		strategyBundle = combined.Strategies

	} else {
		// 👇 İkinci öncelik: -config_json
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

		// 👇 Sadece normal modda strategy yükleme yapılır
		if *useStdinStrategies {
			err = json.NewDecoder(os.Stdin).Decode(&strategyBundle)
			if err != nil {
				fmt.Println("Failed to parse strategies from stdin:", err)
				os.Exit(1)
			}
		} else {
			// Dosyadan yükleme
			if flagIsPassed("strategies") {
				cfg.StrategyDirectory = *strategyDir
			}
			if err := engine.SetStrategyDirectory(cfg.StrategyDirectory); err != nil {
				fmt.Printf("Strategy directory error: %v\n\n", err)
				printUsage()
				os.Exit(1)
			}
		}
	}

	// Simülasyonu başlat
	logger, err := engine.NewLogger(*logPath, cfg.GzipEnabled)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n\n", err)
		printUsage()
		os.Exit(1)
	}
	defer logger.Close()

	eng := engine.NewEngine(cfg, logger, *showProgress, *debug, strategyBundle)
	eng.Run()

	if *debug {
		fmt.Println("Simulation completed. Log written to", logger.FinalPath)
	}
}

func printUsage() {
	fmt.Println("SimJack - Blackjack Simulation")
	fmt.Println("Usage:")
	flag.PrintDefaults()
}
