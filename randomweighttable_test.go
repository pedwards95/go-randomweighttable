package randomweighttable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	wt, clean := New()
	assert.NotNil(t, wt, "TestNew; Expected new WeightTable, got nil")
	clean()
}

func TestAdd(t *testing.T) {
	wt, clean := New()
	defer clean()

	// Add a new item
	assert.True(t, wt.Add("item1", "data1", 10), "TestAdd; Expected Add to return true for new item")

	// Check if the item was added correctly
	key, item := wt.Get()
	assert.Equal(t, "item1", key, "TestAdd; Expected key to be 'item1'")
	assert.Equal(t, "data1", item, "TestAdd; Expected item to be 'data1'")

	// Add weight to existing key without changing item
	assert.True(t, wt.Add("item1", nil, 5), "TestAdd; Expected Add to return true when adding weight to existing key")
	key, item = wt.Get()
	assert.Equal(t, "item1", key, "TestAdd; Expected key to be 'item1' after update")
	assert.Equal(t, "data1", item, "TestAdd; Expected item to not be updated'")

	assert.True(t, wt.Add("item1", "updatedData", 5), "TestAdd; Expected Add to return true when updating item and adding weight")
	key, item = wt.Get()
	assert.Equal(t, "item1", key, "TestAdd; Expected key to be 'item1' after update")
	assert.Equal(t, "updatedData", item, "TestAdd; Expected item to be updated to 'updatedData'")

	// Fail cases
	assert.False(t, wt.Add("", "data2", 10), "TestAdd; Expected Add to return false for empty key")
	assert.False(t, wt.Add("item2", "data2", 0), "TestAdd; Expected Add to return false for zero weight")
	assert.False(t, wt.Add("item2", "data2", -5), "TestAdd; Expected Add to return false for negative weight")
	assert.False(t, wt.Add("item3", nil, 10), "TestAdd; Expected Add to return false for nil item when key does not exist")
}

func TestGet(t *testing.T) {
	wt, clean := New()
	defer clean()

	wt.Add("item1", "data1", 10)
	wt.Add("item2", "data2", 20)

	key, item := wt.Get()
	assert.NotEmpty(t, key, "TestGet; Expected to retrieve a non-empty key")
	assert.NotNil(t, item, "TestGet; Expected to retrieve a non-nil item")
}

func TestDelete(t *testing.T) {
	wt, clean := New()
	defer clean()

	wt.Add("item1", "data1", 10)
	wt.Delete("item1")

	key, item := wt.Get()
	assert.NotEqual(t, "item1", key, "TestDelete; Expected item1 to be deleted")
	assert.NotEqual(t, "data1", item, "TestDelete; Expected item1's data to be deleted")
}

func TestClean(t *testing.T) {
	wt, clean := New()
	wt.Add("item1", "data1", 10)
	clean()

	key, item := wt.Get()
	assert.Empty(t, key, "TestClean; Expected WeightTable to be empty after clean")
	assert.Nil(t, item, "TestClean; Expected WeightTable to return nil item after clean")
}

func TestGetAllProbabilities(t *testing.T) {
	wt, clean := New()
	defer clean()

	// Test when table is empty
	probabilities := wt.GetAllProbabilities()
	assert.Empty(t, probabilities, "TestGetAllProbabilities; Expected empty map when table is empty")

	// Test with actual values
	wt.Add("item1", "data1", 10)
	wt.Add("item2", "data2", 30)
	wt.Add("item3", "data3", 60)

	expected := map[string]float64{
		"item1": 10.00,
		"item2": 30.00,
		"item3": 60.00,
	}

	probabilities = wt.GetAllProbabilities()
	assert.Equal(t, expected, probabilities, "TestGetAllProbabilities; Expected probabilities to match calculated values")
}

func TestGetProbability(t *testing.T) {
	wt, clean := New()
	defer clean()

	// Test when table is empty
	assert.Equal(t, 0.00, wt.GetProbability("item1"), "TestGetProbability; Expected 0.00 probability for any key in empty table")

	wt.Add("item1", "data1", 10)
	wt.Add("item2", "data2", 30)
	wt.Add("item3", "data3", 60)

	assert.Equal(t, 10.00, wt.GetProbability("item1"), "TestGetProbability; Expected probability for item1 to be 10.00")
	assert.Equal(t, 30.00, wt.GetProbability("item2"), "TestGetProbability; Expected probability for item2 to be 30.00")
	assert.Equal(t, 60.00, wt.GetProbability("item3"), "TestGetProbability; Expected probability for item3 to be 60.00")
	assert.Equal(t, 0.00, wt.GetProbability("item4"), "TestGetProbability; Expected probability for non-existent item4 to be 0.00")
}

func TestProbabilityRounding(t *testing.T) {
	wt, clean := New()
	defer clean()

	wt.Add("item1", "data1", 1)
	wt.Add("item2", "data2", 1000)
	wt.Add("item3", "data3", 500)

	probabilities := wt.GetAllProbabilities()
	assert.Equal(t, 0.07, probabilities["item1"], "TestProbabilityRounding; Expected probability for item1 to be 0.07")
	assert.Equal(t, 66.62, probabilities["item2"], "TestProbabilityRounding; Expected probability for item2 to be 66.62")
	assert.Equal(t, 33.31, probabilities["item3"], "TestProbabilityRounding; Expected probability for item3 to be 33.31")
}
