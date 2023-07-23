package db

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CacheItem struct {
	table          Table
	itemName       string
	value          map[string]types.AttributeValue // current item value
	updateInterval time.Duration
	lastAccess     time.Time // how long it's been since the item value was last accessed
}

func MakeItem(
	table Table,
	itemName string,
	updateInterval time.Duration,
) (CacheItem, error) {
	item := CacheItem{
		table:          table,
		itemName:       itemName,
		updateInterval: updateInterval,
	}
	err := item.update()
	return item, err
}

func (i *CacheItem) update() error {
	newValue, err := getItem(i.table, i.itemName)
	if err != nil {
		return err
	}
	i.value = newValue
	i.lastAccess = time.Now()
	return nil
}

func (i *CacheItem) GetValue() (map[string]types.AttributeValue, error) {
	if time.Now().Sub(i.lastAccess) > i.updateInterval {
		err := i.update()
		return i.value, err
	}
	i.lastAccess = time.Now()
	return i.value, nil
}

type Cache struct {
	items          map[string]CacheItem
	updateInterval time.Duration // minimum time between item updates
	flushPeriod    time.Duration // how long items are kept if they're not changed
}

func (c *Cache) cycleFlush() {
	for {
		time.Sleep(c.flushPeriod)
		now := time.Now()
		for key, item := range c.items {
			if now.Sub(item.lastAccess) > c.flushPeriod {
				delete(c.items, key)
			}
		}
	}
}

func MakeCache(updateInterval time.Duration, flushPeriod time.Duration) Cache {
	cache := Cache{updateInterval: updateInterval, flushPeriod: flushPeriod}
	go cache.cycleFlush()
	return cache
}

func (c *Cache) Get(table Table, itemName string) (map[string]types.AttributeValue, error) {
	itemKey := fmt.Sprintf("%s_%s", table.Name(), itemName)
	// check if this item is already cached
	if item, exists := c.items[itemKey]; exists {
		return item.GetValue()
	}
	item, err := MakeItem(table, itemName, c.updateInterval)
	c.items[itemKey] = item
	return item.value, err
}
