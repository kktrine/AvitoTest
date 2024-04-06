package cashe

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	Items             map[int][]Item
}

type Item struct {
	Value           []int
	AdminOnlyAccess bool
	Expiration      time.Time
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[int][]Item)

	cache := Cache{
		Items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC() // данный метод рассматривается ниже
	}

	return &cache
}

// Add element to cache
func (c *Cache) Add(tag int, features []int) error {
	c.Lock()
	defer c.Unlock()
	item, found := c.Items[tag]
	if found {
		for _, existingFeatures := range item {
			for _, feature := range features {
				if slices.Contains(existingFeatures.Value, feature) {
					return errors.New("banner already exists")
				}
			}
		}

	}
	itemToAdd := Item{Value: features, AdminOnlyAccess: false, Expiration: time.Now().Add(c.defaultExpiration)}
	c.Items[tag] = append(c.Items[tag], itemToAdd)
	return nil
}

func (c *Cache) Get(tag, feature int) (int, []int, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.Items[tag]

	// ключ не найден
	if !found {
		return 0, nil, false
	}

	var features []int
	for _, value := range item {
		if slices.Contains(value.Value, feature) {
			features = value.Value
			// Если в момент запроса кеш устарел возвращаем nil
			if time.Now().After(value.Expiration) {
				return 0, nil, false
			}

		}
	}

	if len(features) == 0 {
		return 0, nil, false
	}

	return tag, features, true
}

func (c *Cache) Delete(tag, feature int) error {

	c.Lock()

	defer c.Unlock()

	item, found := c.Items[tag]
	if !found {
		return errors.New("tag not found")
	}
	found = false
	for i, value := range item {
		if slices.Contains(value.Value, feature) {
			found = true
			c.Items[tag] = append(c.Items[tag][:i], c.Items[tag][i+1:]...)
			break
		}
	}
	if !found {
		return errors.New("feature not found")
	}
	if len(c.Items[tag]) == 0 {
		delete(c.Items, tag)
	}
	return nil
}

func (c *Cache) ChangeAccess(tag, feature int) error {
	item, found := c.Items[tag]

	// ключ не найден
	if !found {
		return errors.New("tag not found")
	}
	found = false
	for i, value := range item {
		if slices.Contains(value.Value, feature) {
			c.Items[tag][i].AdminOnlyAccess = true
			found = true
			break
		}
	}
	if !found {
		return errors.New("feature not found")
	}
	return nil
}

func (c *Cache) StartGC() {
	go c.gC()
}

func (c *Cache) gC() {
	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.Items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() map[int][]int {

	c.RLock()

	defer c.RUnlock()
	res := make(map[int][]int)
	for tag, value := range c.Items {
		for i, feature := range value {
			if time.Now().After(feature.Expiration) {
				res[tag] = append(res[tag], i)
			}
		}
	}

	return res
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys map[int][]int) {
	c.Lock()

	defer c.Unlock()

	for tag, nums := range keys {
		if len(nums) == len(c.Items[tag]) {
			delete(c.Items, tag)
			continue
		}
		slices.Reverse(nums)
		fmt.Println(nums)
		for num := range nums {
			c.Items[tag] = append(c.Items[tag][:num], c.Items[tag][num+1:]...)
		}
	}
}
