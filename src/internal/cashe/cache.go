package cashe

import (
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
	Tags            []int
	Content         map[string]interface{}
	AdminOnlyAccess bool
	Expiration      time.Time
}

func NewCashe(defaultExpiration, cleanupInterval time.Duration) *Cache {

	items := make(map[int][]Item)

	cache := Cache{
		Items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.startGC() // данный метод рассматривается ниже
	}

	return &cache
}

// Add element to cache
func (c *Cache) Add(feature int, tags []int) error {
	c.Lock()
	defer c.Unlock()
	//item, found := c.Items[feature]
	//if found {
	//	for _, existingTags := range item {
	//		for _, tag := range tags {
	//			if slices.Contains(existingTags.Tags, tag) {
	//				return errors.New("banner already exists")
	//			}
	//		}
	//	}
	//}
	itemToAdd := Item{Tags: tags, AdminOnlyAccess: false, Expiration: time.Now().Add(c.defaultExpiration)}
	c.Items[feature] = append(c.Items[feature], itemToAdd)
	return nil
}

func (c *Cache) Get(feature, tag int) (int, []int, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.Items[feature]

	// ключ не найден
	if !found {
		return 0, nil, false
	}

	for _, value := range item {
		if slices.Contains(value.Tags, tag) {

			// Если в момент запроса кеш устарел возвращаем nil
			if time.Now().After(value.Expiration) {
				return 0, nil, false
			}
			return feature, value.Tags, true

		}
	}

	return 0, nil, false

}

//func (c *Cache) Delete(feature, tag int) error {
//	c.Lock()
//	defer c.Unlock()
//
//	item, found := c.Items[feature]
//	if !found {
//		return errors.New("tag not found")
//	}
//	found = false
//	for i, value := range item {
//		if slices.Contains(value.Tags, tag) {
//			found = true
//			c.Items[feature] = append(c.Items[feature][:i], c.Items[feature][i+1:]...)
//			break
//		}
//	}
//	if !found {
//		return errors.New("feature not found")
//	}
//	if len(c.Items[feature]) == 0 {
//		delete(c.Items, tag)
//	}
//	return nil
//}

//func (c *Cache) ChangeAccess(feature, tag int) error {
//	item, found := c.Items[feature]
//
//	// ключ не найден
//	if !found {
//		return errors.New("tag not found")
//	}
//	found = false
//	for i, value := range item {
//		if slices.Contains(value.Tags, tag) {
//			c.Items[feature][i].AdminOnlyAccess = true
//			found = true
//			break
//		}
//	}
//	if !found {
//		return errors.New("feature not found")
//	}
//	return nil
//}

func (c *Cache) startGC() {
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
	for feature, value := range c.Items {
		for i, banner := range value {
			if time.Now().After(banner.Expiration) {
				res[feature] = append(res[feature], i)
			}
		}
	}

	return res
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys map[int][]int) {
	c.Lock()

	defer c.Unlock()

	for feature, nums := range keys {
		if len(nums) == len(c.Items[feature]) {
			delete(c.Items, feature)
			continue
		}
		slices.Reverse(nums)
		fmt.Println(nums)
		for num := range nums {
			c.Items[feature] = append(c.Items[feature][:num], c.Items[feature][num+1:]...)
		}
	}
}
