package api

import "fmt"

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  []PokemonStat
	Types  []PokemonType
}

type PokemonStat struct {
	Name  string
	Value int
}

type PokemonType struct {
	Name string
}

func (p Pokemon) Print() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Name, s.Value)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Name)
	}
}
