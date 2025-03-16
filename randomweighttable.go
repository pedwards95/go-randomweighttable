package randomweighttable

import (
	"math"
	"math/rand/v2"
	"sync"
)

// WeightTable cache settings
type WeightTable struct {
	mutex       sync.RWMutex
	items       map[string]*Item
	totalWeight int64
}

// Item carries return objects
type Item struct {
	weight int64
	item   any
}

// New weighttable
func New() (*WeightTable, func()) {
	wt := &WeightTable{
		items:       make(map[string]*Item),
		totalWeight: 0,
	}
	return wt, wt.clean
}

// Clean empties the WeightTable
func (wt *WeightTable) clean() {
	wt.mutex.Lock()
	defer wt.mutex.Unlock()
	clear(wt.items)
}

/*
Add puts an interface into the weight table. Beware, this table never cleans itself, and is unbound in size.

If the object is already in the table, it overrides the item, and adds the 2 weights together.
If item is nil, it will not override the item, but will still add the weights together.

Returns bool of success.
*/
func (wt *WeightTable) Add(key string, item any, weight int64) bool {
	wt.mutex.Lock()
	defer wt.mutex.Unlock()

	if weight <= 0 || len(key) == 0 {
		return false
	}
	if _, ok := wt.items[key]; !ok {
		if item == nil {
			return false
		}
		wt.items[key] = &Item{item: item, weight: weight}
		wt.totalWeight += weight
		return true
	}
	itm := wt.items[key]
	itm.weight += weight
	wt.totalWeight += weight
	if item != nil {
		itm.item = item
	}
	return true
}

// Get a random selection from the WeightTable
func (wt *WeightTable) Get() (key string, item any) {
	wt.mutex.RLock()
	defer wt.mutex.RUnlock()

	if wt.totalWeight <= 0 {
		return "", nil
	}

	randomNumber := int64(rand.IntN(int(wt.totalWeight))) + 1

	for key, item := range wt.items {
		randomNumber -= item.weight
		if randomNumber <= 0 {
			return key, item.item
		}
	}
	return "", nil
}

// Delete an entry from the WeightTable
func (wt *WeightTable) Delete(key string) {
	wt.mutex.Lock()
	defer wt.mutex.Unlock()

	if item, ok := wt.items[key]; ok {
		wt.totalWeight -= item.weight
		delete(wt.items, key)
	}
}

// GetAllProbabilities returns the percentage chance of picking each key, rounded to the nearest hundredth
func (wt *WeightTable) GetAllProbabilities() map[string]float64 {
	wt.mutex.RLock()
	defer wt.mutex.RUnlock()

	probabilities := make(map[string]float64)
	if wt.totalWeight == 0 {
		return probabilities
	}

	for key, item := range wt.items {
		probabilities[key] = math.Round((float64(item.weight)/float64(wt.totalWeight)*100)*100) / 100
	}

	return probabilities
}

// GetProbability returns the percentage chance of picking a specific key, rounded to the nearest hundredth
func (wt *WeightTable) GetProbability(key string) float64 {
	wt.mutex.RLock()
	defer wt.mutex.RUnlock()

	if wt.totalWeight == 0 {
		return 0.0
	}

	if item, ok := wt.items[key]; ok {
		return math.Round((float64(item.weight)/float64(wt.totalWeight)*100)*100) / 100
	}
	return 0.0
}
