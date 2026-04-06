package main

import (
	"crypto/sha512"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

var attempts uint64 = 0

// Algorand base32 alphabet
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// Map char → 5-bit value
var charMap = func() map[rune]uint8 {
	m := make(map[rune]uint8)
	for i, c := range alphabet {
		m[c] = uint8(i)
	}
	return m
}()

// Convert prefix string → bit array
func prefixToBits(prefix string) []bool {
	var bits []bool
	for _, c := range prefix {
		val := charMap[c]
		for i := 4; i >= 0; i-- {
			bit := (val >> i) & 1
			bits = append(bits, bit == 1)
		}
	}
	return bits
}

// Compute Algorand checksum (last 4 bytes of SHA512/256)
func checksum(pk []byte) []byte {
	hash := sha512.Sum512_256(pk)
	return hash[len(hash)-4:]
}

// Bit comparison without encoding
func matchesPrefix(pk []byte, targetBits []bool) bool {
	chk := checksum(pk)

	// Combine pk + checksum (36 bytes total)
	full := make([]byte, 36)
	copy(full[:32], pk)
	copy(full[32:], chk)

	bitIndex := 0

	for _, b := range full {
		for i := 7; i >= 0; i-- {
			if bitIndex >= len(targetBits) {
				return true
			}

			bit := (b >> i) & 1
			if (bit == 1) != targetBits[bitIndex] {
				return false
			}

			bitIndex++
		}
	}

	return true
}

func worker(targetBits []bool, found *atomic.Bool, wg *sync.WaitGroup, resultChan chan<- crypto.Account) {
	defer wg.Done()

	for !found.Load() {
		account := crypto.GenerateAccount()

		atomic.AddUint64(&attempts, 1)

		if matchesPrefix(account.Address[:], targetBits) {
			if found.CompareAndSwap(false, true) {
				resultChan <- account
			}
			return
		}
	}
}

func main() {
	// ✅ Prefix
	targetPrefix := "ABCD"
	if len(os.Args) > 1 {
		targetPrefix = strings.ToUpper(os.Args[1])
	}

	targetBits := prefixToBits(targetPrefix)

	start := time.Now()

	fmt.Println("Searching for prefix:", targetPrefix)
	fmt.Println("Start time:", start.Format(time.RFC3339))

	// ✅ Max CPU usage
	numWorkers := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Workers:", numWorkers)

	var wg sync.WaitGroup
	var found atomic.Bool
	resultChan := make(chan crypto.Account, 1)

	// ✅ Workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(targetBits, &found, &wg, resultChan)
	}

	// ✅ Progress
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("Attempts:", atomic.LoadUint64(&attempts))
		}
	}()

	// ✅ Wait for result
	account := <-resultChan

	end := time.Now()
	duration := end.Sub(start)

	mn, err := mnemonic.FromPrivateKey(account.PrivateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n✅ FOUND!")
	fmt.Println("Address:", account.Address.String())
	fmt.Println("Mnemonic:", mn)

	fmt.Println("\nStats:")
	fmt.Println("Total Attempts:", atomic.LoadUint64(&attempts))
	fmt.Println("Start Time:", start.Format(time.RFC3339))
	fmt.Println("End Time:", end.Format(time.RFC3339))
	fmt.Println("Total Time:", duration)

	found.Store(true)
	wg.Wait()
}
