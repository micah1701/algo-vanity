# ⚡ Algorand Vanity Address Generator (Go)

Generate custom Algorand wallet addresses that start with a desired prefix — **fast**.

This project started as a simple curiosity about how vanity address generators work and evolved into a highly optimized, multi-core brute-force engine written in Go.

---

## 🚀 What This Does

This tool generates random Algorand keypairs until it finds an address that starts with a user-defined prefix.

Example:

ABCD... HELLO... DOGE...

When a match is found, it outputs:

- ✅ Algorand Address
- ✅ 25-word mnemonic (seed phrase)
- ✅ Total attempts
- ✅ Execution time

---

## ⚙️ How It Works

Algorand addresses are derived from:

base32( publicKey (32 bytes) + checksum (4 bytes) )

Each character in the address represents **5 bits**, meaning:

| Prefix Length | Difficulty          |
| ------------- | ------------------- |
| 4 chars       | ~1 million attempts |
| 5 chars       | ~33 million         |
| 6 chars       | ~1 billion          |
| 7 chars       | ~34 billion         |

This tool brute-forces key generation until a match is found.

---

## 🔥 Performance Evolution

This project went through several major optimizations:

### ❌ Naive Approach (Browser / Basic JS)

- Single-threaded
- Full base32 encoding per attempt
- Extremely slow (hours+ for 4 chars)

---

### ✅ Go + Goroutines

- Multi-core parallelism
- Massive speed improvement

---

### ⚡ Removed `Address.String()` Calls

- Avoided repeated base32 encoding
- Reduced allocations
- ~2–4x speed boost

---

### 🚀 Final Optimization: **NO Base32 Encoding**

Instead of encoding addresses:

✅ Convert target prefix → bit pattern  
✅ Generate `(publicKey + checksum)`  
✅ Compare bits directly

This eliminates the biggest bottleneck entirely.

---

## 🧠 Why This Is Fast

- ✅ Zero string allocations in hot loop
- ✅ Bit-level comparison (early exit)
- ✅ Full CPU utilization
- ✅ Minimal memory overhead

---

## 📊 Real Performance Example

On a typical machine (e.g., my AMD Ryzen 7 5825U Mini PC):

~230,000 attempts/sec

### Observed results:

- ✅ 5-character prefix in **15 seconds** (lucky run)
- ✅ 5-character prefix in **~2 minutes** (expected)

---

## 🛠️ Installation

### 1. Clone the repo

```bash
git clone https://github.com/micah1701/algo-vanity.git
cd algo-vanity
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run

```bash
go run main.go
```

### Or Compile to executable for even fasterness

```
go build -o algo-vanity.exe

./algo-vanity.exe HELLO
```

## ▶️ Usage

### Default (prefix = ABCD)

```
go run main.go
```

### Custom Prefix

```
go run main.go HELLO
```

## 🧾 Example Output

```
Searching for prefix: HELLO
Start time: 2026-04-06T12:00:00Z

✅ FOUND!
Address: HELLOXYZ...
Mnemonic: ability zone ... (25 words)

Stats:
Total Attempts: 25254774
Total Time: 1m48s
```

## ⚠️ Important Notes

- The mnemonic is the only thing you need to recover the wallet

- Store it securely — this tool does not encrypt or save keys

- Each run is random — times will vary due to probability
