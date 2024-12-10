package main

import (
	"os"
	"testing"
)

func TestPositions_FindDeltas(t *testing.T) {
	data, err := os.ReadFile("test_input.txt")
	if err != nil {
		t.Error(err)
	}

	m := LoadLocations(data)
	locations := m.locations

	if len(locations['0']) != 4 {
		t.Errorf("Expected 4 antennas at 0, got %d", len(locations['0']))
	}
	if len(locations['A']) != 3 {
		t.Errorf("Expected 3 antennas at A, got %d", len(locations['A']))
	}

	deltas := locations['0'].FindDeltas()
	if len(deltas) != 6 {
		t.Errorf("Expected 6 deltas, got %d", len(deltas))
	}
	deltas = locations['A'].FindDeltas()
	if len(deltas) != 3 {
		t.Errorf("Expected 5 deltas, got %d", len(deltas))
	}

	if m.CountUniqueAntinodes() != 14 {
		t.Errorf("Expected 14 unique anti nodes, got %d", m.CountUniqueAntinodes())
	}

	if m.CountResonantAntinodes() != 34 {
		t.Errorf("Expected 34 resonant antinodes, got %d", m.CountResonantAntinodes())
	}
}
