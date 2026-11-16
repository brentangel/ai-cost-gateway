package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type CacheEntry struct {
	Response  []byte
	ExpiresAt time.Time
}

type Metrics struct {
	TotalRequests    uint64  `json:"total_requests"`
	CacheHits        uint64  `json:"cache_hits"`
	CacheMisses      uint64  `json:"cache_misses"`
	EstimatedSavings float64 `json:"estimated_savings_usd"`
	EnterpriseActive bool    `json:"enterprise_license_active"`
}

var (
	cacheMap      = make(map[string]CacheEntry)
	cacheMu       sync.RWMutex
	totalRequests uint64
	cacheHits     uint64
	cacheMisses   uint64
	isEnterprise  bool
)

// Average cost saved per cached enterprise LLM call
const avgCostPerCall = 0.03

func main() {
	// Check for Proprietary Enterprise License Key
	entKey := os.Getenv("ENTERPRISE_LICENSE_KEY")
	if len(entKey) > 10 && entKey[:4] == "ENT-" {
		isEnterprise = true
		log.Println("💎 ENTERPRISE LICENSE VERIFIED: Advanced TTL & Audit Logs Unlocked.")
	} else {
		log.Println("🚀 OPEN CORE EDITION: Running under Apache 2.0 License.")
	}

	http.HandleFunc("/v1/chat/completions", handleProxy)
	http.HandleFunc("/metrics", handleMetrics)

	log.Println("AI Cost-Control Gateway running on :8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	hits := atomic.LoadUint64(&cacheHits)
	savings := float64(hits) * avgCostPerCall

	m := Metrics{
		TotalRequests:    atomic.LoadUint64(&totalRequests),
		CacheHits:        hits,
		CacheMisses:      atomic.LoadUint64(&cacheMisses),
		EstimatedSavings: savings,
		EnterpriseActive: isEnterprise,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	atomic.AddUint64(&totalRequests, 1)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	hash := sha256.Sum256(bodyBytes)
	cacheKey := hex.EncodeToString(hash[:])

	cacheMu.RLock()
	entry, found := cacheMap[cacheKey]
	cacheMu.RUnlock()

	if found && time.Now().Before(entry.ExpiresAt) {
		atomic.AddUint64(&cacheHits, 1)
		log.Println("⚡ Cache HIT! Saving API costs.")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Gateway-Cache", "HIT")
		w.Write(entry.Response)
		return
	}

	atomic.AddUint64(&cacheMisses, 1)
	log.Println("❌ Cache MISS. Forwarding request to OpenAI...")

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	targetURL := "https://openai.com"
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Proxy creation error", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Upstream API error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading upstream response", http.StatusInternalServerError)
		return
	}

	// Enterprise tier gets 24-hour caching; Open Core gets 1-hour caching
	ttl := 1 * time.Hour
	if isEnterprise {
		ttl = 24 * time.Hour
	}

	cacheMu.Lock()
	cacheMap[cacheKey] = CacheEntry{
		Response:  respBytes,
		ExpiresAt: time.Now().Add(ttl),
	}
	cacheMu.Unlock()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("X-Gateway-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBytes)
}
