package datatypes

import "fmt"

type Inventory[T any] struct {
	Label T
	Stock map[string]int
}

// Copy
func (i *Inventory[T]) Copy() *Inventory[T] {
	stock := make(map[string]int, len(i.Stock))
	for k, v := range i.Stock {
		stock[k] = v
	}
	return &Inventory[T]{
		Label: i.Label,
		Stock: stock,
	}
}

func (i *Inventory[T]) Add(item string) {
	i.AddX(item, 1)
}

func (i *Inventory[T]) AddX(item string, x int) {
	i.Stock[item] += x
}

func (i *Inventory[T]) SetX(item string, x int) {
	i.Stock[item] = x
}

func (i *Inventory[T]) RemoveX(item string, x int) error {
	if i.Stock[item] < x {
		return fmt.Errorf("not enough %s", item)
	}

	i.Stock[item] -= x
	return nil
}

func (i *Inventory[T]) Remove(item string) error {
	return i.RemoveX(item, 1)
}

func (i *Inventory[T]) Count(item string) int {
	return i.Stock[item]
}

func NewInventory[T any](label T) *Inventory[T] {
	return &Inventory[T]{
		Label: label,
		Stock: make(map[string]int),
	}
}
