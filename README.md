# SimJack

🎲 **SimJack** is a high-performance, fully customizable blackjack simulation engine written in Go.

> "Where strategy meets simulation."

---

## 🚀 Features

- ♠️ Full blackjack game engine (7-box table, dealer AI, split/double support)
- 🧠 Strategy files loaded dynamically via JSON (no recompile!)
- 💼 Supports perfect pair & 21+3 sidebets
- 📊 Outputs a rich, Pandas-ready CSV log file
- ⚡ Handles millions of hands efficiently with buffered logging
- 🧪 Supports forced cards, custom config per player
- ✅ CLI-ready & API-compatible design

---

## 📦 Project Structure

```
simjack/
├── main.go              # CLI entry point
├── config/              # Config schema
├── engine/              # Game logic (Dealer, Player, Box, Hand, Strategy, Logger...)
├── strategies/          # Strategy definitions (e.g. basic_chart.json)
├── test_config.json     # Sample config to run out-of-the-box
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

## 📄 License

MIT — see `LICENSE`

---

Happy simulating! 🧠♠️
