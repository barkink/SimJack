# SimJack

![Go](https://img.shields.io/badge/Go-1.20+-blue)
![License](https://img.shields.io/github/license/barkink/SimJack)
![Platform](https://img.shields.io/badge/platform-cli-lightgrey)

🎲 **SimJack** is a high-performance, fully customizable blackjack simulation engine written in Go.

It supports card counting, dynamic bet ramping, sidebets (Perfect Pair, 21+3), and custom strategy files in JSON format — perfect for strategy testing or statistical analysis.

> "Where strategy meets simulation."

---

## 🚀 Features

- ♠️ Full blackjack game engine (7-box table, dealer AI, split/double support)
- 🧠 Strategy files loaded dynamically via JSON (no recompile!)
- 🧠 Custom strategy support (basic & card counting)
- 📈 Deviation rules based on true count
- 🎯 Bet ramping with true count multipliers
- 💼 Supports perfect pair & 21+3 sidebets
- 📦 JSON-configurable players, rules, and simulations
- 📊 Outputs a rich, Pandas-ready CSV log file
- ⚡ Handles millions of hands efficiently with buffered logging
- 🧪 Supports forced cards, custom config per player
- ✅ CLI-ready & API-compatible design

---

## 🚀 Getting Started

```bash
git clone https://github.com/barkink/SimJack.git
cd SimJack
go run main.go -config=test_config.json -log=simjack_log.csv -strategies=strategies/
```

---

## ⚙️ Usage

### 🏗 Build

```bash
go build -o simjack main.go
```

### 🧪 Run a Simulation

With config file:

```bash
./simjack -config=test_config.json -log=results.csv -strategies=strategies/
```

With inline JSON:

```bash
./simjack -config_json='{"round_count":100000, "num_decks":6, ...}' -log=results.csv
```

### 🆘 Help

```bash
./simjack -help
```
---
## 📦 Project Structure

```
simjack/
├── main.go              # CLI entry point
├── config/              # Config schema
├── engine/              # Game logic (Dealer, Player, Box, Hand, Strategy, Logger...)
├── strategies/          # Strategy definitions (e.g. basic_chart.json) Could be another directory, given as a parameter.
├── test_config.json     # Sample config to run out-of-the-box
```

---

## 📈 Log Format

Each row = one hand  
Log is saved as `logname_0.csv` while simulation is running  
Renamed to `logname_1.csv` when finished

Columns include:

- Player and box info
- Bet amounts and payouts
- Sidebets and insurance
- Hand cards, result, split/double flags
- Dealer cards and outcome
- Shoe info (deck state, cut position, round/total cards drawn)

Perfectly suited for Pandas analysis 🎯

---

## 📋 Strategy Format

Each strategy is a `.json` file defining:

```json
{
  "decide_insurance": true,
  "fallback": "stand",
  "actions": {
    "hard_13_vs_2": ["stand"],
    "pair_8_vs_10": ["split", "stand"],
    "soft_17_vs_3": ["hit"]
  }
}
```

Supports:
- `hard_X_vs_Y`
- `soft_X_vs_Y`
- `pair_R_vs_Y`

---

## 🔬 Custom Config

Config allows:
- Per-player balance, bet unit, split limit
- Box assignment
- Sidebet enablement
- Forced cards for debugging

See `test_config.json` for a working example.

---

---

## 🔧 Sample Config (test_config.json)

```json
{
  "num_decks": 6,
  "round_count": 1000,
  "strategy_directory": "strategies",
  "players": [
    {
      "player_id": 1,
      "strategy": "hi_lo",
      "bet_unit": 10,
      "box_indexes": [1, 2]
    }
  ]
}
```

---

## 📚 Strategy Sample (hi_lo.json)

```json
{
  "fallback": "stand",
  "counting_enabled": true,
  "actions": {
    "hard_16_vs_10": ["hit"]
  },
  "deviations": {
    "hard_16_vs_10": { "at_count": 4, "action": "stand" }
  },
  "bet_ramp": [
    { "min_count": -4, "bet_unit": 0.5 },
    { "min_count": 1,  "bet_unit": 2.0 },
    { "min_count": 3,  "bet_unit": 4.0 }
  ]
}
```

---

## 📊 Output Log

The simulator produces a CSV log file containing:
- 📌 strategy_key and full decision_trace per hand
- 🧠 indicators like (deviation), (fallback)
- 💵 bet_unit vs. bet_unit_used for ramp-up tracking
- 📈 full financials (balance, payout, sidebets, insurance)

Example log row:
```
..., hard_16_vs_10:stand (deviation), 10.00, 40.00
```

---

## 📄 License

MIT — see `LICENSE`

---

Happy simulating! 🧠♠️
