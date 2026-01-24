package main

import "testing"

func TestIsTechGig(t *testing.T) {
	config := Config{
		Keywords: []string{"python", "scraping", "design"},
	}

	tests := []struct {
		name     string
		gig      Gig
		expected bool
	}{
		{
			name: "Tech Gig",
			gig: Gig{
				Title:       "Python Developer Needed",
				Description: "Looking for someone to build a scraping bot.",
				Skills:      []string{"python", "scraping"},
			},
			expected: true,
		},
		{
			name: "Non-Tech Gig (if keywords match)",
			gig: Gig{
				Title:       "Logo Designer",
				Description: "Design a logo for a bakery.",
				Skills:      []string{"design", "photoshop"},
			},
			expected: true, // "design" is in our mock config keywords
		},
		{
			name: "Irrelevant Gig",
			gig: Gig{
				Title:       "Chef Needed",
				Description: "Cook food.",
				Skills:      []string{"cooking"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTechGig(tt.gig, config); got != tt.expected {
				t.Errorf("isTechGig() = %v, want %v", got, tt.expected)
			}
		})
	}
}
