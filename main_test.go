package main

import "testing"

func TestIsTechGig(t *testing.T) {
	gig := Gig{
		Title:       "Python Developer Needed",
		Description: "Looking for someone to build a scraping bot.",
		Skills:      []string{"python", "scraping"},
	}
	if !isTechGig(gig) {
		t.Errorf("Expected gig to be tech-related")
	}

	gig2 := Gig{
		Title:       "Logo Designer",
		Description: "Design a logo for a bakery.",
		Skills:      []string{"design", "photoshop"},
	}
	if isTechGig(gig2) {
		t.Errorf("Expected gig to NOT be tech-related")
	}
}
