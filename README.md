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
./simjack -config=test_config.json -log=results.csv -strategies=strategies
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
round,shoe,deck_running_count,true_count,real_count_till_cut_card,box_id,player_id,hand_id,owner,hand,result,bet_from_config,bet_unit_used,hand_payout,main_payout,pp_bet,pp_win,pp_type,p21_bet,p21_win,p21_type,insurance_bet,insurance_payout,initial_balance,round_start_balance,player_balance,is_blackjack,is_doubled,is_split_child,split_count,dealer_upcard,dealer_final_hand,dealer_blackjack,dealer_bust,player_bust,player_draws,player_is_bankrupt,player_is_retired,num_decks,cut_card_position,cards_drawn_total,cards_drawn_round,cards_left_after_round,strategy_key
1,1,4,0,7,1,1,1,PlayerOne,Q of Clubs;7 of Diamonds,lose,10.00,5.00,0.00,0.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999993.00,999986.00,False,False,False,0,K of Spades,K of Spades;Q of Spades,False,False,False,,False,False,8,238,19,19,397,hard_17_vs_10:stand (fallback)
1,1,4,0,7,2,1,1,PlayerOne,3 of Diamonds;5 of Diamonds;5 of Clubs;A of Spades;7 of Clubs,win,10.00,5.00,20.00,20.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999993.00,1000006.00,False,False,False,0,K of Spades,K of Spades;Q of Spades,False,False,False,5 of Clubs;A of Spades;7 of Clubs,False,False,8,238,19,19,397,hard_8_vs_10:hit (main)
1,1,4,0,7,3,2,1,PlayerTwo,Q of Diamonds;9 of Diamonds,lose,5.00,5.00,0.00,0.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799990.00,799985.00,False,False,False,0,K of Spades,K of Spades;Q of Spades,False,False,False,,False,False,8,238,19,19,397,hard_19_vs_10:stand (fallback)
1,1,4,0,7,4,2,1,PlayerTwo,J of Diamonds;5 of Hearts;5 of Clubs,push,5.00,5.00,5.00,5.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799990.00,799990.00,False,False,False,0,K of Spades,K of Spades;Q of Spades,False,False,False,5 of Clubs,False,False,8,238,19,19,397,hard_15_vs_10:hit (main)
1,1,4,0,7,5,2,1,PlayerTwo,4 of Hearts;2 of Hearts;3 of Clubs;5 of Hearts;6 of Hearts,push,5.00,5.00,5.00,5.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799990.00,799995.00,False,False,False,0,K of Spades,K of Spades;Q of Spades,False,False,False,3 of Clubs;5 of Hearts;6 of Hearts,False,False,8,238,19,19,397,hard_6_vs_10:hit (main)
2,1,14,1,17,1,1,1,PlayerOne,5 of Hearts;6 of Hearts;9 of Hearts,lose,20.00,40.00,0.00,11.00,1.00,0.00,none,1.00,11.00,Straight,0.00,0.00,1000000.00,999964.00,999923.00,False,True,False,0,7 of Diamonds,7 of Diamonds;4 of Spades;K of Clubs,False,False,False,9 of Hearts,False,False,8,238,39,20,377,hard_11_vs_7:double (main)
2,1,14,1,17,2,1,1,PlayerOne,6 of Clubs;2 of Clubs;6 of Diamonds;8 of Clubs,lose,10.00,40.00,0.00,0.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999964.00,999923.00,False,False,False,0,7 of Diamonds,7 of Diamonds;4 of Spades;K of Clubs,False,False,True,6 of Diamonds;8 of Clubs,False,False,8,238,39,20,377,hard_8_vs_7:hit (main)
2,1,14,1,17,3,2,1,PlayerTwo,4 of Clubs;7 of Spades;3 of Hearts,lose,10.00,5.00,0.00,0.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799985.00,799970.00,False,True,False,0,7 of Diamonds,7 of Diamonds;4 of Spades;K of Clubs,False,False,False,3 of Hearts,False,False,8,238,39,20,377,hard_11_vs_7:double (main)
2,1,14,1,17,4,2,1,PlayerTwo,6 of Diamonds;3 of Diamonds;7 of Clubs;3 of Spades,lose,5.00,5.00,0.00,0.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799985.00,799970.00,False,False,False,0,7 of Diamonds,7 of Diamonds;4 of Spades;K of Clubs,False,False,False,7 of Clubs;3 of Spades,False,False,8,238,39,20,377,hard_9_vs_7:hit (main)
2,1,14,1,17,5,2,1,PlayerTwo,2 of Clubs;9 of Diamonds;Q of Clubs,push,10.00,5.00,10.00,10.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799985.00,799980.00,False,True,False,0,7 of Diamonds,7 of Diamonds;4 of Spades;K of Clubs,False,False,False,Q of Clubs,False,False,8,238,39,20,377,hard_11_vs_7:double (main)
3,1,16,2,19,1,1,1,PlayerOne,6 of Spades;3 of Clubs;K of Diamonds,win,10.00,60.00,20.00,20.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999861.00,999819.00,False,False,False,0,8 of Spades,8 of Spades;9 of Diamonds,False,False,False,K of Diamonds,False,False,8,238,56,17,360,hard_9_vs_8:hit (main)
3,1,16,2,19,2,1,1,PlayerOne,J of Diamonds;2 of Diamonds;8 of Clubs,win,10.00,60.00,20.00,20.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999861.00,999839.00,False,False,False,0,8 of Spades,8 of Spades;9 of Diamonds,False,False,False,8 of Clubs,False,False,8,238,56,17,360,hard_12_vs_8:hit (main)
3,1,16,2,19,3,2,1,PlayerTwo,8 of Hearts;10 of Spades,win,5.00,5.00,10.00,10.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799970.00,799975.00,False,False,False,0,8 of Spades,8 of Spades;9 of Diamonds,False,False,False,,False,False,8,238,56,17,360,hard_18_vs_8:stand (fallback)
3,1,16,2,19,4,2,1,PlayerTwo,2 of Diamonds;2 of Diamonds;8 of Hearts;3 of Diamonds;4 of Spades,win,5.00,5.00,10.00,10.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799970.00,799985.00,False,False,False,0,8 of Spades,8 of Spades;9 of Diamonds,False,False,False,8 of Hearts;3 of Diamonds;4 of Spades,False,False,8,238,56,17,360,pair_2_vs_8:hit (main)
3,1,16,2,19,5,2,1,PlayerTwo,10 of Clubs;J of Spades,win,5.00,5.00,10.00,10.00,0.00,0.00,none,0.00,0.00,none,0.00,0.00,800000.00,799970.00,799995.00,False,False,False,0,8 of Spades,8 of Spades;9 of Diamonds,False,False,False,,False,False,8,238,56,17,360,hard_20_vs_8:stand (fallback)
4,1,12,1,15,1,1,1,PlayerOne,10 of Hearts;J of Clubs,win,10.00,60.00,20.00,20.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999777.00,999735.00,False,False,False,0,K of Spades,K of Spades;3 of Spades;J of Spades,False,True,False,,False,False,8,238,73,17,343,hard_20_vs_10:stand (fallback)
4,1,12,1,15,2,1,1,PlayerOne,4 of Diamonds;Q of Diamonds;A of Diamonds;7 of Clubs,lose,10.00,60.00,0.00,0.00,1.00,0.00,none,1.00,0.00,none,0.00,0.00,1000000.00,999777.00,999735.00,False,False,False,0,K of Spades,K of Spades;3 of Spades;J of Spades,False,True,True,A of Diamonds;7 of Clubs,False,False,8,238,73,17,343,hard_14_vs_10:hit (main)
.
.
.
```

---

⚡️ How Long Would 10 Million Hands Take in Real Life?

SimJack can simulate 10 million hands in a single run.But how does that compare to real-life gameplay?

Scenario

Real Life (Casino)

SimJack (Simulation)

Average time per round

~5 seconds

< 0.002 seconds

1 player x 10,000,000 rounds

~578 days (non-stop)

~2 minutes

7-player full table

~4s x 7 = 28s/round

still ~2 minutes

🌟 SimJack compresses 1.5 years of gameplay into just 2 minutes.

🧠 Why It Matters

With SimJack, you can:

Make decisions based on statistical confidence, not just gut feeling

Avoid wasting days (or bankroll) testing bad strategies

Gather more data in 1 hour than 1 year of table experience

---

## 📄 License

MIT — see `LICENSE`

---

Happy simulating! 🧠♠️
