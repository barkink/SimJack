# SimJack

🎲 **SimJack**, Go diliyle yazılmış, yüksek performanslı ve tamamen özelleştirilebilir bir blackjack simülasyon motorudur.

Kart sayma, dinamik bahis ayarlamaları, yan bahisler (Perfect Pair, 21+3) ve JSON formatında özel strateji dosyalarını destekler — strateji testleri ve istatistiksel analizler için idealdir.

> "Stratejinin simülasyonla buluştuğu yer."

---

## 🚀 Özellikler

- ♠️ Tam blackjack oyun motoru (7 koltuklu masa, krupiye yapay zekâsı, bölme/ikiye katlama desteği)
- 🧠 JSON üzerinden dinamik olarak yüklenen strateji dosyaları (derlemeye gerek yok!)
- 🧠 Özel strateji desteği (temel & kart sayma stratejileri)
- 📈 Gerçek sayıya (true count) göre sapma kuralları
- 🎯 Gerçek sayıya bağlı dinamik bahis sistemi
- 💼 Perfect Pair & 21+3 yan bahis desteği
- 📦 JSON ile yapılandırılabilir oyuncular, kurallar ve simülasyonlar
- 📊 Pandas ile analiz edilebilecek zengin CSV log çıktısı
- ⚡ Milyonlarca eli tamponlu loglama ile verimli şekilde işler
- 🧪 Zorunlu kartlar ve oyuncuya özel yapılandırmalar desteklenir
- ✅ Komut satırı kullanımına hazır ve API uyumlu tasarım

---

## 🚀 Hızlı Başlangıç

```bash
git clone https://github.com/barkink/SimJack.git
cd SimJack
go run main.go -config=test_config.json -log=simjack_log.csv -strategies=strategies
```

---

## ⚙️ Kullanım

### 🏗 Derleme

```bash
go build -o simjack main.go
```

### 🧪 Simülasyon Çalıştırma

Yapılandırma dosyasıyla:

```bash
./simjack -config=test_config.json -log=results.csv -strategies=strategies
```

Inline JSON ile:

```bash
./simjack -config_json='{"round_count":100000, "num_decks":6, ...}' -log=results.csv
```

### 🆘 Yardım

```bash
./simjack -help
```

---

## 📦 Proje Yapısı

```
simjack/
├── main.go              # CLI giriş noktası
├── config/              # Yapılandırma şeması
├── engine/              # Oyun motoru (Krupiye, Oyuncu, Kutu, El, Strateji, Logger...)
├── strategies/          # Strateji tanımları (örn. basic_chart.json)
├── test_config.json     # Çalıştırılabilir örnek yapılandırma dosyası
```

---

## 📈 Log Formatı

Her satır = bir el  
Simülasyon çalışırken log `logname_0.csv` olarak kaydedilir  
Simülasyon bitince `logname_1.csv` olarak yeniden adlandırılır

Sütunlar şunları içerir:

- Oyuncu ve kutu bilgileri
- Bahis miktarları ve ödemeler
- Yan bahisler ve sigorta
- Oyuncu kartları, sonuçlar, bölme/ikiye katlama bilgileri
- Krupiye kartları ve sonucu
- Deste bilgisi (kalan kartlar, kesme noktası, çekilen kartlar)

Pandas analizi için idealdir 🎯

---

## 📋 Strateji Formatı

Her strateji `.json` dosyası olarak tanımlanır:

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

Desteklenen anahtarlar:
- `hard_X_vs_Y`
- `soft_X_vs_Y`
- `pair_R_vs_Y`

---

## 🔬 Özel Yapılandırmalar

Yapılandırma dosyası şunlara izin verir:
- Oyuncuya özel bakiye, bahis birimi, bölme limiti
- Kutulara atama
- Yan bahis etkinleştirme
- Hata ayıklama için zorunlu kart tanımlama

Detaylar için `test_config.json` dosyasına bakınız.

---

## 🔧 Örnek Yapılandırma (test_config.json)

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

## 📚 Strateji Örneği (hi_lo.json)

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

## 📊 Log Çıktısı

Simülasyon çıktısı olarak CSV dosyası üretilir ve şu bilgileri içerir:
- 📌 her el için `strategy_key` ve karar izleri
- 🧠 (deviation), (fallback) gibi işaretleyiciler
- 💵 kullanılan bahis birimi ile planlanan birimin karşılaştırması
- 📈 detaylı finansal bilgiler (bakiye, ödeme, yan bahisler, sigorta)

Örnek Log formatı (Okunabilirlik için boşluk eklenmiştir):
```
round, shoe, deck_running_count, true_count, real_count_till_cut_card, box_id, player_id, hand_id, owner    , hand                                                                                 , result   , bet_from_config, bet_unit_used, hand_payout, main_payout, pp_bet, pp_win, pp_type     , p21_bet, p21_win, p21_type       , insurance_bet, insurance_payout, initial_balance, round_start_balance, player_balance, is_blackjack, is_doubled, is_split_child, split_count, dealer_upcard , dealer_final_hand                                                        , dealer_blackjack, dealer_bust, player_bust, player_draws                                                 , player_is_bankrupt, player_is_retired, num_decks, cut_card_position, cards_drawn_total, cards_drawn_round, cards_left_after_round, strategy_key
    1,    1, 4                 , 0         , 7                       ,      1,         1,       1, PlayerOne, Q of Clubs;7 of Diamonds                                                             , lose     ,           10.00,          5.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999993.00,      999986.00, False       , False     , False         ,           0, K of Spades   , K of Spades;Q of Spades                                                  , False           , False      , False      ,                                                              , False             , False            ,         8,               238,                19,                19,                    397, hard_17_vs_10:stand (fallback)
    1,    1, 4                 , 0         , 7                       ,      2,         1,       1, PlayerOne, 3 of Diamonds;5 of Diamonds;5 of Clubs;A of Spades;7 of Clubs                        , win      ,           10.00,          5.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999993.00,     1000006.00, False       , False     , False         ,           0, K of Spades   , K of Spades;Q of Spades                                                  , False           , False      , False      , 5 of Clubs;A of Spades;7 of Clubs                            , False             , False            ,         8,               238,                19,                19,                    397, hard_8_vs_10:hit (main)
    1,    1, 4                 , 0         , 7                       ,      3,         2,       1, PlayerTwo, Q of Diamonds;9 of Diamonds                                                          , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799990.00,      799985.00, False       , False     , False         ,           0, K of Spades   , K of Spades;Q of Spades                                                  , False           , False      , False      ,                                                              , False             , False            ,         8,               238,                19,                19,                    397, hard_19_vs_10:stand (fallback)
    1,    1, 4                 , 0         , 7                       ,      4,         2,       1, PlayerTwo, J of Diamonds;5 of Hearts;5 of Clubs                                                 , push     ,            5.00,          5.00,        5.00,        5.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799990.00,      799990.00, False       , False     , False         ,           0, K of Spades   , K of Spades;Q of Spades                                                  , False           , False      , False      , 5 of Clubs                                                   , False             , False            ,         8,               238,                19,                19,                    397, hard_15_vs_10:hit (main)
    1,    1, 4                 , 0         , 7                       ,      5,         2,       1, PlayerTwo, 4 of Hearts;2 of Hearts;3 of Clubs;5 of Hearts;6 of Hearts                           , push     ,            5.00,          5.00,        5.00,        5.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799990.00,      799995.00, False       , False     , False         ,           0, K of Spades   , K of Spades;Q of Spades                                                  , False           , False      , False      , 3 of Clubs;5 of Hearts;6 of Hearts                           , False             , False            ,         8,               238,                19,                19,                    397, hard_6_vs_10:hit (main)
    2,    1, 14                , 1         , 17                      ,      1,         1,       1, PlayerOne, 5 of Hearts;6 of Hearts;9 of Hearts                                                  , lose     ,           20.00,         40.00,        0.00,       11.00,   1.00,   0.00, none        ,    1.00,   11.00, Straight       ,          0.00,             0.00,      1000000.00,           999964.00,      999923.00, False       , True      , False         ,           0, 7 of Diamonds , 7 of Diamonds;4 of Spades;K of Clubs                                     , False           , False      , False      , 9 of Hearts                                                  , False             , False            ,         8,               238,                39,                20,                    377, hard_11_vs_7:double (main)
    2,    1, 14                , 1         , 17                      ,      2,         1,       1, PlayerOne, 6 of Clubs;2 of Clubs;6 of Diamonds;8 of Clubs                                       , lose     ,           10.00,         40.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999964.00,      999923.00, False       , False     , False         ,           0, 7 of Diamonds , 7 of Diamonds;4 of Spades;K of Clubs                                     , False           , False      , True       , 6 of Diamonds;8 of Clubs                                     , False             , False            ,         8,               238,                39,                20,                    377, hard_8_vs_7:hit (main)
    2,    1, 14                , 1         , 17                      ,      3,         2,       1, PlayerTwo, 4 of Clubs;7 of Spades;3 of Hearts                                                   , lose     ,           10.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      799970.00, False       , True      , False         ,           0, 7 of Diamonds , 7 of Diamonds;4 of Spades;K of Clubs                                     , False           , False      , False      , 3 of Hearts                                                  , False             , False            ,         8,               238,                39,                20,                    377, hard_11_vs_7:double (main)
    2,    1, 14                , 1         , 17                      ,      4,         2,       1, PlayerTwo, 6 of Diamonds;3 of Diamonds;7 of Clubs;3 of Spades                                   , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      799970.00, False       , False     , False         ,           0, 7 of Diamonds , 7 of Diamonds;4 of Spades;K of Clubs                                     , False           , False      , False      , 7 of Clubs;3 of Spades                                       , False             , False            ,         8,               238,                39,                20,                    377, hard_9_vs_7:hit (main)
    2,    1, 14                , 1         , 17                      ,      5,         2,       1, PlayerTwo, 2 of Clubs;9 of Diamonds;Q of Clubs                                                  , push     ,           10.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      799980.00, False       , True      , False         ,           0, 7 of Diamonds , 7 of Diamonds;4 of Spades;K of Clubs                                     , False           , False      , False      , Q of Clubs                                                   , False             , False            ,         8,               238,                39,                20,                    377, hard_11_vs_7:double (main)
    3,    1, 16                , 2         , 19                      ,      1,         1,       1, PlayerOne, 6 of Spades;3 of Clubs;K of Diamonds                                                 , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999861.00,      999819.00, False       , False     , False         ,           0, 8 of Spades   , 8 of Spades;9 of Diamonds                                                , False           , False      , False      , K of Diamonds                                                , False             , False            ,         8,               238,                56,                17,                    360, hard_9_vs_8:hit (main)
    3,    1, 16                , 2         , 19                      ,      2,         1,       1, PlayerOne, J of Diamonds;2 of Diamonds;8 of Clubs                                               , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999861.00,      999839.00, False       , False     , False         ,           0, 8 of Spades   , 8 of Spades;9 of Diamonds                                                , False           , False      , False      , 8 of Clubs                                                   , False             , False            ,         8,               238,                56,                17,                    360, hard_12_vs_8:hit (main)
    3,    1, 16                , 2         , 19                      ,      3,         2,       1, PlayerTwo, 8 of Hearts;10 of Spades                                                             , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799970.00,      799975.00, False       , False     , False         ,           0, 8 of Spades   , 8 of Spades;9 of Diamonds                                                , False           , False      , False      ,                                                              , False             , False            ,         8,               238,                56,                17,                    360, hard_18_vs_8:stand (fallback)
    3,    1, 16                , 2         , 19                      ,      4,         2,       1, PlayerTwo, 2 of Diamonds;2 of Diamonds;8 of Hearts;3 of Diamonds;4 of Spades                    , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799970.00,      799985.00, False       , False     , False         ,           0, 8 of Spades   , 8 of Spades;9 of Diamonds                                                , False           , False      , False      , 8 of Hearts;3 of Diamonds;4 of Spades                        , False             , False            ,         8,               238,                56,                17,                    360, pair_2_vs_8:hit (main)
    3,    1, 16                , 2         , 19                      ,      5,         2,       1, PlayerTwo, 10 of Clubs;J of Spades                                                              , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799970.00,      799995.00, False       , False     , False         ,           0, 8 of Spades   , 8 of Spades;9 of Diamonds                                                , False           , False      , False      ,                                                              , False             , False            ,         8,               238,                56,                17,                    360, hard_20_vs_8:stand (fallback)
    4,    1, 12                , 1         , 15                      ,      1,         1,       1, PlayerOne, 10 of Hearts;J of Clubs                                                              , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999777.00,      999735.00, False       , False     , False         ,           0, K of Spades   , K of Spades;3 of Spades;J of Spades                                      , False           , True       , False      ,                                                              , False             , False            ,         8,               238,                73,                17,                    343, hard_20_vs_10:stand (fallback)
    4,    1, 12                , 1         , 15                      ,      2,         1,       1, PlayerOne, 4 of Diamonds;Q of Diamonds;A of Diamonds;7 of Clubs                                 , lose     ,           10.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999777.00,      999735.00, False       , False     , False         ,           0, K of Spades   , K of Spades;3 of Spades;J of Spades                                      , False           , True       , True       , A of Diamonds;7 of Clubs                                     , False             , False            ,         8,               238,                73,                17,                    343, hard_14_vs_10:hit (main)
    4,    1, 12                , 1         , 15                      ,      3,         2,       1, PlayerTwo, A of Hearts;J of Hearts                                                              , blackjack,            5.00,          5.00,       12.50,       12.50,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      799992.50, True        , False     , False         ,           0, K of Spades   , K of Spades;3 of Spades;J of Spades                                      , False           , True       , False      ,                                                              , False             , False            ,         8,               238,                73,                17,                    343, soft_21_vs_10:stand (fallback)
    4,    1, 12                , 1         , 15                      ,      4,         2,       1, PlayerTwo, 4 of Hearts;6 of Spades;2 of Spades;8 of Clubs                                       , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      800002.50, False       , False     , False         ,           0, K of Spades   , K of Spades;3 of Spades;J of Spades                                      , False           , True       , False      , 2 of Spades;8 of Clubs                                       , False             , False            ,         8,               238,                73,                17,                    343, hard_10_vs_10:hit (main)
    4,    1, 12                , 1         , 15                      ,      5,         2,       1, PlayerTwo, 7 of Hearts;A of Spades                                                              , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      800012.50, False       , False     , False         ,           0, K of Spades   , K of Spades;3 of Spades;J of Spades                                      , False           , True       , False      ,                                                              , False             , False            ,         8,               238,                73,                17,                    343, hard_18_vs_10:stand (fallback)
    5,    1, 11                , 1         , 14                      ,      1,         1,       1, PlayerOne, 9 of Clubs;10 of Clubs                                                               , push     ,           10.00,         60.00,       10.00,       10.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999673.00,      999611.00, False       , False     , False         ,           0, 9 of Hearts   , 9 of Hearts;10 of Spades                                                 , False           , False      , False      ,                                                              , False             , False            ,         8,               238,                89,                16,                    327, hard_19_vs_9:stand (fallback)
    5,    1, 11                , 1         , 14                      ,      2,         1,       1, PlayerOne, 4 of Spades;7 of Spades;6 of Diamonds                                                , lose     ,           20.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999673.00,      999611.00, False       , True      , False         ,           0, 9 of Hearts   , 9 of Hearts;10 of Spades                                                 , False           , False      , False      , 6 of Diamonds                                                , False             , False            ,         8,               238,                89,                16,                    327, hard_11_vs_9:double (main)
    5,    1, 11                , 1         , 14                      ,      3,         2,       1, PlayerTwo, K of Hearts;8 of Spades                                                              , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           800002.50,      799997.50, False       , False     , False         ,           0, 9 of Hearts   , 9 of Hearts;10 of Spades                                                 , False           , False      , False      ,                                                              , False             , False            ,         8,               238,                89,                16,                    327, hard_18_vs_9:stand (fallback)
    5,    1, 11                , 1         , 14                      ,      4,         2,       1, PlayerTwo, A of Spades;4 of Diamonds;K of Spades;9 of Spades                                    , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           800002.50,      799997.50, False       , False     , False         ,           0, 9 of Hearts   , 9 of Hearts;10 of Spades                                                 , False           , False      , True       , K of Spades;9 of Spades                                      , False             , False            ,         8,               238,                89,                16,                    327, soft_15_vs_9:hit (main)
    5,    1, 11                , 1         , 14                      ,      5,         2,       1, PlayerTwo, 3 of Diamonds;4 of Diamonds;Q of Clubs                                               , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           800002.50,      799997.50, False       , False     , False         ,           0, 9 of Hearts   , 9 of Hearts;10 of Spades                                                 , False           , False      , False      , Q of Clubs                                                   , False             , False            ,         8,               238,                89,                16,                    327, hard_7_vs_9:hit (main)
    6,    1, 10                , 1         , 13                      ,      1,         1,       1, PlayerOne, 10 of Hearts;6 of Diamonds                                                           , lose     ,           10.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999549.00,      999487.00, False       , False     , False         ,           0, J of Hearts   , J of Hearts;7 of Hearts                                                  , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               104,                15,                    312, hard_16_vs_10:stand (deviation)
    6,    1, 10                , 1         , 13                      ,      2,         1,       1, PlayerOne, 10 of Spades;3 of Spades;9 of Diamonds                                               , lose     ,           10.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999549.00,      999487.00, False       , False     , False         ,           0, J of Hearts   , J of Hearts;7 of Hearts                                                  , False           , False      , True       , 9 of Diamonds                                                , False             , False            ,         8,               238,               104,                15,                    312, hard_13_vs_10:hit (main)
    6,    1, 10                , 1         , 13                      ,      3,         2,       1, PlayerTwo, 9 of Spades;2 of Clubs;9 of Clubs                                                    , win      ,           10.00,          5.00,       20.00,       20.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799987.50,      799997.50, False       , True      , False         ,           0, J of Hearts   , J of Hearts;7 of Hearts                                                  , False           , False      , False      , 9 of Clubs                                                   , False             , False            ,         8,               238,               104,                15,                    312, hard_11_vs_10:double (main)
    6,    1, 10                , 1         , 13                      ,      4,         2,       1, PlayerTwo, 9 of Spades;10 of Spades                                                             , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799987.50,      800007.50, False       , False     , False         ,           0, J of Hearts   , J of Hearts;7 of Hearts                                                  , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               104,                15,                    312, hard_19_vs_10:stand (fallback)
    6,    1, 10                , 1         , 13                      ,      5,         2,       1, PlayerTwo, 8 of Diamonds;5 of Spades;J of Spades                                                , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799987.50,      800007.50, False       , False     , False         ,           0, J of Hearts   , J of Hearts;7 of Hearts                                                  , False           , False      , True       , J of Spades                                                  , False             , False            ,         8,               238,               104,                15,                    312, hard_13_vs_10:hit (main)
    7,    1, 11                , 1         , 14                      ,      1,         1,       1, PlayerOne, 10 of Diamonds;K of Spades                                                           , push     ,           10.00,         60.00,       10.00,       10.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999425.00,      999373.00, False       , False     , False         ,           0, 10 of Hearts  , 10 of Hearts;K of Clubs                                                  , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               119,                15,                    297, hard_20_vs_10:stand (fallback)
    7,    1, 11                , 1         , 14                      ,      2,         1,       1, PlayerOne, 5 of Clubs;Q of Clubs                                                                , lose     ,           10.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999425.00,      999373.00, False       , False     , False         ,           0, 10 of Hearts  , 10 of Hearts;K of Clubs                                                  , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               119,                15,                    297, hard_15_vs_10:stand (deviation)
    7,    1, 11                , 1         , 14                      ,      3,         2,       1, PlayerTwo, 5 of Clubs;J of Clubs;4 of Diamonds                                                  , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799997.50,      799987.50, False       , False     , False         ,           0, 10 of Hearts  , 10 of Hearts;K of Clubs                                                  , False           , False      , False      , 4 of Diamonds                                                , False             , False            ,         8,               238,               119,                15,                    297, hard_15_vs_10:hit (main)
    7,    1, 11                , 1         , 14                      ,      4,         2,       1, PlayerTwo, 5 of Spades;6 of Diamonds;7 of Diamonds                                              , lose     ,           10.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799997.50,      799987.50, False       , True      , False         ,           0, 10 of Hearts  , 10 of Hearts;K of Clubs                                                  , False           , False      , False      , 7 of Diamonds                                                , False             , False            ,         8,               238,               119,                15,                    297, hard_11_vs_10:double (main)
    7,    1, 11                , 1         , 14                      ,      5,         2,       1, PlayerTwo, 7 of Hearts;6 of Diamonds;6 of Hearts                                                , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799997.50,      799987.50, False       , False     , False         ,           0, 10 of Hearts  , 10 of Hearts;K of Clubs                                                  , False           , False      , False      , 6 of Hearts                                                  , False             , False            ,         8,               238,               119,                15,                    297, hard_13_vs_10:hit (main)
    8,    1, 11                , 2         , 14                      ,      1,         1,       1, PlayerOne, 8 of Diamonds;10 of Clubs                                                            , lose     ,           10.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999311.00,      999249.00, False       , False     , False         ,           0, 7 of Spades   , 7 of Spades;4 of Spades;3 of Clubs;5 of Diamonds                         , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               137,                18,                    279, hard_18_vs_7:stand (fallback)
    8,    1, 11                , 2         , 14                      ,      2,         1,       1, PlayerOne, 4 of Hearts;A of Diamonds;8 of Hearts;A of Clubs;9 of Hearts                         , lose     ,           10.00,         60.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999311.00,      999249.00, False       , False     , False         ,           0, 7 of Spades   , 7 of Spades;4 of Spades;3 of Clubs;5 of Diamonds                         , False           , False      , True       , 8 of Hearts;A of Clubs;9 of Hearts                           , False             , False            ,         8,               238,               137,                18,                    279, hard_15_vs_7:hit (main)
    8,    1, 11                , 2         , 14                      ,      3,         2,       1, PlayerTwo, 7 of Diamonds;9 of Hearts                                                            , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799977.50,      799972.50, False       , False     , False         ,           0, 7 of Spades   , 7 of Spades;4 of Spades;3 of Clubs;5 of Diamonds                         , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               137,                18,                    279, hard_16_vs_7:stand (main)
    8,    1, 11                , 2         , 14                      ,      4,         2,       1, PlayerTwo, 8 of Clubs;K of Spades                                                               , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799977.50,      799972.50, False       , False     , False         ,           0, 7 of Spades   , 7 of Spades;4 of Spades;3 of Clubs;5 of Diamonds                         , False           , False      , False      ,                                                              , False             , False            ,         8,               238,               137,                18,                    279, hard_18_vs_7:stand (fallback)
    8,    1, 11                , 2         , 14                      ,      5,         2,       1, PlayerTwo, 8 of Spades;4 of Spades;Q of Clubs                                                   , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799977.50,      799972.50, False       , False     , False         ,           0, 7 of Spades   , 7 of Spades;4 of Spades;3 of Clubs;5 of Diamonds                         , False           , False      , True       , Q of Clubs                                                   , False             , False            ,         8,               238,               137,                18,                    279, hard_12_vs_7:hit (main)
    9,    1, 8                 , 1         , 11                      ,      1,         1,       1, PlayerOne, J of Spades;6 of Hearts                                                              , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999187.00,      999145.00, False       , False     , False         ,           0, 3 of Diamonds , 3 of Diamonds;5 of Clubs;4 of Diamonds;A of Spades;Q of Diamonds         , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               153,                16,                    263, hard_16_vs_3:stand (main)
    9,    1, 8                 , 1         , 11                      ,      2,         1,       1, PlayerOne, Q of Hearts;8 of Spades                                                              , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999187.00,      999165.00, False       , False     , False         ,           0, 3 of Diamonds , 3 of Diamonds;5 of Clubs;4 of Diamonds;A of Spades;Q of Diamonds         , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               153,                16,                    263, hard_18_vs_3:stand (fallback)
    9,    1, 8                 , 1         , 11                      ,      3,         2,       1, PlayerTwo, K of Diamonds;2 of Diamonds;J of Hearts                                              , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799962.50,      799957.50, False       , False     , False         ,           0, 3 of Diamonds , 3 of Diamonds;5 of Clubs;4 of Diamonds;A of Spades;Q of Diamonds         , False           , True       , True       , J of Hearts                                                  , False             , False            ,         8,               238,               153,                16,                    263, hard_12_vs_3:hit (main)
    9,    1, 8                 , 1         , 11                      ,      4,         2,       1, PlayerTwo, A of Hearts;Q of Clubs                                                               , blackjack,            5.00,          5.00,       12.50,       12.50,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799962.50,      799970.00, True        , False     , False         ,           0, 3 of Diamonds , 3 of Diamonds;5 of Clubs;4 of Diamonds;A of Spades;Q of Diamonds         , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               153,                16,                    263, soft_21_vs_3:stand (fallback)
    9,    1, 8                 , 1         , 11                      ,      5,         2,       1, PlayerTwo, 3 of Spades;K of Clubs                                                               , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799962.50,      799980.00, False       , False     , False         ,           0, 3 of Diamonds , 3 of Diamonds;5 of Clubs;4 of Diamonds;A of Spades;Q of Diamonds         , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               153,                16,                    263, hard_13_vs_3:stand (main)
   10,    1, 5                 , 1         , 8                       ,      1,         1,       1, PlayerOne, 10 of Spades;3 of Hearts                                                             , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999103.00,      999061.00, False       , False     , False         ,           0, 6 of Spades   , 6 of Spades;10 of Hearts;7 of Hearts                                     , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               166,                13,                    250, hard_13_vs_6:stand (main)
   10,    1, 5                 , 1         , 8                       ,      2,         1,       1, PlayerOne, 5 of Spades;10 of Diamonds                                                           , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999103.00,      999081.00, False       , False     , False         ,           0, 6 of Spades   , 6 of Spades;10 of Hearts;7 of Hearts                                     , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               166,                13,                    250, hard_15_vs_6:stand (main)
   10,    1, 5                 , 1         , 8                       ,      3,         2,       1, PlayerTwo, 10 of Clubs;J of Clubs                                                               , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799970.00,      799975.00, False       , False     , False         ,           0, 6 of Spades   , 6 of Spades;10 of Hearts;7 of Hearts                                     , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               166,                13,                    250, hard_20_vs_6:stand (fallback)
   10,    1, 5                 , 1         , 8                       ,      4,         2,       1, PlayerTwo, Q of Clubs;6 of Hearts                                                               , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799970.00,      799985.00, False       , False     , False         ,           0, 6 of Spades   , 6 of Spades;10 of Hearts;7 of Hearts                                     , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               166,                13,                    250, hard_16_vs_6:stand (main)
   10,    1, 5                 , 1         , 8                       ,      5,         2,       1, PlayerTwo, J of Spades;9 of Clubs                                                               , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799970.00,      799995.00, False       , False     , False         ,           0, 6 of Spades   , 6 of Spades;10 of Hearts;7 of Hearts                                     , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               166,                13,                    250, hard_19_vs_6:stand (fallback)
   11,    1, 7                 , 1         , 10                      ,      1,         1,       1, PlayerOne, A of Clubs;3 of Spades;8 of Diamonds;4 of Hearts;7 of Spades                         , lose     ,           10.00,         60.00,        0.00,       11.00,   1.00,   0.00, none        ,    1.00,   11.00, Straight       ,          0.00,             0.00,      1000000.00,           999019.00,      998968.00, False       , False     , False         ,           0, 2 of Hearts   , 2 of Hearts;4 of Clubs;2 of Diamonds;7 of Hearts;K of Hearts             , False           , True       , True       , 8 of Diamonds;4 of Hearts;7 of Spades                        , False             , False            ,         8,               238,               184,                18,                    232, soft_14_vs_2:hit (main)
   11,    1, 7                 , 1         , 10                      ,      2,         1,       1, PlayerOne, 4 of Hearts;A of Diamonds                                                            , win      ,           10.00,         60.00,       20.00,       20.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           999019.00,      998988.00, False       , False     , False         ,           0, 2 of Hearts   , 2 of Hearts;4 of Clubs;2 of Diamonds;7 of Hearts;K of Hearts             , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               184,                18,                    232, hard_15_vs_2:stand (main)
   11,    1, 7                 , 1         , 10                      ,      3,         2,       1, PlayerTwo, 9 of Hearts;J of Diamonds                                                            , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      799990.00, False       , False     , False         ,           0, 2 of Hearts   , 2 of Hearts;4 of Clubs;2 of Diamonds;7 of Hearts;K of Hearts             , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               184,                18,                    232, hard_19_vs_2:stand (fallback)
   11,    1, 7                 , 1         , 10                      ,      4,         2,       1, PlayerTwo, 5 of Hearts;10 of Diamonds                                                           , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      800000.00, False       , False     , False         ,           0, 2 of Hearts   , 2 of Hearts;4 of Clubs;2 of Diamonds;7 of Hearts;K of Hearts             , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               184,                18,                    232, hard_15_vs_2:stand (main)
   11,    1, 7                 , 1         , 10                      ,      5,         2,       1, PlayerTwo, K of Diamonds;6 of Spades                                                            , win      ,            5.00,          5.00,       10.00,       10.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           799985.00,      800010.00, False       , False     , False         ,           0, 2 of Hearts   , 2 of Hearts;4 of Clubs;2 of Diamonds;7 of Hearts;K of Hearts             , False           , True       , False      ,                                                              , False             , False            ,         8,               238,               184,                18,                    232, hard_16_vs_2:stand (main)
   12,    2, -2                , 0         , -3                      ,      1,         1,       1, PlayerOne, 9 of Hearts;K of Diamonds                                                            , lose     ,           10.00,          5.00,        0.00,        0.00,   1.00,   0.00, none        ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           998981.00,      998974.00, False       , False     , False         ,           0, K of Spades   , K of Spades;J of Diamonds                                                , False           , False      , False      ,                                                              , False             , False            ,         8,               224,                16,                16,                    400, hard_19_vs_10:stand (fallback)
   12,    2, -2                , 0         , -3                      ,      2,         1,       1, PlayerOne, 9 of Hearts;9 of Clubs                                                               , lose     ,           10.00,          5.00,        0.00,        7.00,   1.00,   7.00, Mixed Pair  ,    1.00,    0.00, none           ,          0.00,             0.00,      1000000.00,           998981.00,      998981.00, False       , False     , False         ,           0, K of Spades   , K of Spades;J of Diamonds                                                , False           , False      , False      ,                                                              , False             , False            ,         8,               224,                16,                16,                    400, pair_9_vs_10:stand (main)
   12,    2, -2                , 0         , -3                      ,      3,         2,       1, PlayerTwo, 7 of Clubs;Q of Diamonds                                                             , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           800000.00,      799995.00, False       , False     , False         ,           0, K of Spades   , K of Spades;J of Diamonds                                                , False           , False      , False      ,                                                              , False             , False            ,         8,               224,                16,                16,                    400, hard_17_vs_10:stand (fallback)
   12,    2, -2                , 0         , -3                      ,      4,         2,       1, PlayerTwo, 6 of Hearts;7 of Spades;A of Clubs;J of Clubs                                        , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           800000.00,      799995.00, False       , False     , False         ,           0, K of Spades   , K of Spades;J of Diamonds                                                , False           , False      , True       , A of Clubs;J of Clubs                                        , False             , False            ,         8,               224,                16,                16,                    400, hard_13_vs_10:hit (main)
   12,    2, -2                , 0         , -3                      ,      5,         2,       1, PlayerTwo, 9 of Clubs;4 of Hearts;2 of Clubs;4 of Diamonds                                      , lose     ,            5.00,          5.00,        0.00,        0.00,   0.00,   0.00, none        ,    0.00,    0.00, none           ,          0.00,             0.00,       800000.00,           800000.00,      799995.00, False       , False     , False         ,           0, K of Spades   , K of Spades;J of Diamonds                                                , False           , False      , False      , 2 of Clubs;4 of Diamonds                                     , False             , False            ,         8,               224,                16,                16,                    400, hard_13_vs_10:hit (main)

.
.
.
```

---

## 📊 Log Sütunlarının Açıklamaları

| Sütun                      | Açıklama                                           |
|---------------------------|----------------------------------------------------|
| `round`                   | Tur numarası                                       |
| `shoe`                    | Deste numarası (yeniden karıştırma sayacı)         |
| `deck_running_count`      | Oyuncuya göre deste sayacı                         |
| `true_count`              | Gerçek sayaç (sayım ÷ kalan deste)                 |
| `real_count_till_cut_card`| Kesme kartına kadar olan gerçek sayım              |
| `box_id`                  | Masa kutusu ID’si                                  |
| `player_id`               | Oyuncu ID’si                                       |
| `hand_id`                 | Bu kutuya ait elin ID’si                           |
| `owner`                   | Oyuncu adı veya etiketi                            |
| `hand`                    | Elde bulunan kartlar (noktalı virgülle ayrılmış)  |
| `result`                  | Sonuç: win, lose, push, blackjack                  |
| `bet_from_config`         | Yapılandırmadan gelen bahis birimi                |
| `bet_unit_used`           | Kullanılan gerçek bahis birimi                    |
| `hand_payout`             | Toplam ödeme (yan bahis dahil)                    |
| `main_payout`             | Ana bahis üzerinden yapılan ödeme                 |
| `pp_bet`                  | Perfect Pair yan bahsi miktarı                    |
| `pp_win`                  | Perfect Pair kazancı                              |
| `pp_type`                 | Perfect Pair sonucu türü (örn. Mixed, Suited)     |
| `p21_bet`                 | 21+3 yan bahis miktarı                            |
| `p21_win`                 | 21+3 kazancı                                      |
| `p21_type`                | 21+3 sonucu türü                                  |
| `insurance_bet`           | Sigorta bahsi miktarı                             |
| `insurance_payout`        | Sigorta kazancı                                   |
| `initial_balance`         | Simülasyon başındaki başlangıç bakiyesi          |
| `round_start_balance`     | Tur başındaki oyuncu bakiyesi                     |
| `player_balance`          | Tur sonunda kalan oyuncu bakiyesi                 |
| `is_blackjack`            | El blackjack mi?                                  |
| `is_doubled`              | İkiye katlandı mı?                                |
| `is_split_child`          | Bu el bölünmeden mi geldi?                        |
| `split_count`             | Toplam bölünme sayısı                             |
| `dealer_upcard`           | Krupiyenin görünen kartı                          |
| `dealer_final_hand`       | Krupiyenin son eli                                |
| `dealer_blackjack`        | Krupiye blackjack yaptı mı?                       |
| `dealer_bust`             | Krupiye battı mı?                                 |
| `player_bust`             | Oyuncu battı mı?                                  |
| `player_draws`            | İlk iki karttan sonra çekilen kartlar            |
| `player_is_bankrupt`      | Oyuncu minimum bahis yapamıyor mu?               |
| `player_is_retired`       | Oyuncu hedef bakiyeye ulaştı mı?                 |
| `num_decks`               | Kullanılan deste sayısı                           |
| `cut_card_position`       | Kesme kartı konumu                                |
| `cards_drawn_total`       | Desteden çekilen toplam kart sayısı              |
| `cards_drawn_round`       | Bu turda çekilen kart sayısı                     |
| `cards_left_after_round`  | Tur sonrası destede kalan kart sayısı            |
| `strategy_key`            | Bu ele uygulanan strateji                         |

---

## ⚡️ Gerçek Hayatta 10 Milyon El Ne Kadar Sürerdi?

SimJack, tek çalıştırmada 10 milyon eli simüle edebilir. Peki bu gerçek hayatta ne kadar sürerdi?

| Senaryo                   | Gerçek Hayat (Kumarhane)   | SimJack (Simülasyon) |
|---------------------------|-----------------------------|-----------------------|
| Tur başına ortalama süre  | ~5 saniye                   | < 0.002 saniye        |
| 1 oyuncu x 10M el         | ~578 gün (durmaksızın)      | ~2 dakika             |
| 7 oyunculu tam masa       | ~28 saniye/tur              | yine ~2 dakika        |

> 🌟 **SimJack 1.5 yıllık oyunu 2 dakikaya sıkıştırır.**

---

## 🧠 Neden Önemli?

SimJack ile:

* 🎯 Kararlarınızı içgüdüye değil, **istatistiksel güvene** dayandırırsınız
* ⏳ Kötü stratejileri test etmek için günlerinizi (ve paranızı) harcamazsınız
* 📊 Gerçek masada 1 yılda toplanabilecek veriyi **1 saatte** toplayabilirsiniz

---

## 📄 Lisans

MIT — `LICENSE` dosyasına bakınız.

---

İyi simülasyonlar! 🧠♠️