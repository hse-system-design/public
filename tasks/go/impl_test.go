package ordcol

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"sort"
	"testing"
)

func TestCollection_AddAndGet(t *testing.T) {
	collection := NewCollection()

	collection.Add(NewItem(2, 4))
	collection.Add(NewItem(5, 9))

	t.Run("existing keys", func(t *testing.T) {
		item1, err := collection.At(2)
		require.NoError(t, err)
		require.Equal(t, 4, item1.Value())

		item2, err := collection.At(5)
		require.NoError(t, err)
		require.Equal(t, 9, item2.Value())
	})

	t.Run("non-existing keys", func(t *testing.T) {
		_, err := collection.At(12)
		require.Error(t, err)
	})
}

func TestCollection_Empty(t *testing.T) {
	collection := NewCollection()

	t.Run("at", func(t *testing.T) {
		_, err := collection.At(1)
		require.Error(t, err)
	})

	t.Run("iter by insert", func(t *testing.T) {
		iter := collection.IterateBy(ByInsertion)
		require.False(t, iter.HasNext())
	})

	t.Run("iter by keys", func(t *testing.T) {
		iter := collection.IterateBy(ByKey)
		require.False(t, iter.HasNext())
	})
}

func TestCollection_IterateByInsertion(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	collection := NewCollection()

	for _, key := range keys {
		collection.Add(NewItem(key, rand.Int()))
	}

	iter := collection.IterateBy(ByInsertion)

	t.Run("equal order", func(t *testing.T) {
		var lookupKeys []int
		for iter.HasNext() {
			item, err := iter.Next()
			require.NoError(t, err)

			lookupKeys = append(lookupKeys, item.Key())
		}

		require.Equal(t, keys, lookupKeys)
	})

	t.Run("error on empty iterator", func(t *testing.T) {
		_, err := iter.Next()
		require.Error(t, err)
	})
}

func TestCollection_IterateByKeys(t *testing.T) {
	var keys []int
	for i := 0; i < 10; i++ {
		keys = append(keys, rand.Int())
	}

	collection := NewCollection()

	for _, key := range keys {
		collection.Add(NewItem(key, rand.Int()))
	}

	iter := collection.IterateBy(ByKey)

	t.Run("equal order", func(t *testing.T) {
		var lookupKeys []int
		for iter.HasNext() {
			iter, err := iter.Next()
			require.NoError(t, err)

			lookupKeys = append(lookupKeys, iter.Key())
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		require.Equal(t, keys, lookupKeys)
	})

	t.Run("error on empty iterator", func(t *testing.T) {
		_, err := iter.Next()
		require.Error(t, err)
	})
}

func TestPanicOnWrongIterationOrder(t *testing.T) {
	collection := NewCollection()

	require.Panics(t, func() {
		collection.IterateBy(IterationOrder(100))
	})
}