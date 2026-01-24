package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Gig struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
}

func main() {
	data, err := os.ReadFile("gigs.json")
	if err != nil {
		panic(err)
	}

	var gigs []Gig
	if err := json.Unmarshal(data, &gigs); err != nil {
		panic(err)
	}

	skillCounts := make(map[string]int)
	titleKeywords := make(map[string]int)

	fmt.Println("=== ALL GIG TITLES ===")
	for _, g := range gigs {
		fmt.Printf("- %s\n", g.Title)
		
		for _, s := range g.Skills {
			skillCounts[s]++
		}
		
		// Simple title keyword extraction
		words := strings.Fields(strings.ToLower(g.Title))
		for _, w := range words {
			if len(w) > 3 { // skip small words
				titleKeywords[w]++
			}
		}
	}

	fmt.Println("\n=== TOP SKILLS ===")
	printTopMap(skillCounts, 15)

	fmt.Println("\n=== TOP TITLE KEYWORDS ===")
	printTopMap(titleKeywords, 15)
}

func printTopMap(m map[string]int, n int) {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	if n > len(ss) {
		n = len(ss)
	}
	for i := 0; i < n; i++ {
		fmt.Printf("%s: %d\n", ss[i].Key, ss[i].Value)
	}
}
