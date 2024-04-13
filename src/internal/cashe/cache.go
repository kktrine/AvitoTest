package cashe

import (
	"banner/models"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	Items             map[string]Item
}

type Item struct {
	BannerID  int32
	FeatureID int32
	TagIDs    []int32
	IsActive  bool
	//UpdatedAt string
	//CreatedAt string
	Content    models.JSONMap
	Expiration time.Time
}

func NewCache() *Cache {

	items := make(map[string]Item)
	exp, err := time.ParseDuration(os.Getenv("CACHE_EXPIRATION"))
	if err != nil {
		panic("Can't parse CACHE_EXPIRATION: " + err.Error())
	}
	clean, err := time.ParseDuration(os.Getenv("CACHE_CLEANUP_INTERVAL"))
	if err != nil {
		panic("Can't parse CACHE_CLEANUP_INTERVAL: " + err.Error())
	}
	cache := Cache{
		Items:             items,
		defaultExpiration: exp,
		cleanupInterval:   clean,
	}
	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	cache.startGC() // данный метод рассматривается ниже

	return &cache
}

// AddOne element to cache
func (c *Cache) AddOne(banner Item) {
	c.Lock()
	defer c.Unlock()
	key := strconv.Itoa(int(banner.BannerID)) + "_" + strconv.Itoa(int(banner.FeatureID))
	banner.Expiration = time.Now().Add(c.defaultExpiration)
	c.Items[key] = banner
}

func (c *Cache) Get(feature, tag int32) (models.JSONMap, bool) {
	c.RLock()
	defer c.RUnlock()
	for key, value := range c.Items {
		parts := strings.Split(key, "_")
		if f, _ := strconv.Atoi(parts[1]); f == int(feature) && slices.Contains(value.TagIDs, tag) {
			if value.Expiration.Before(time.Now()) {
				return nil, false
			}
			return value.Content, value.IsActive
		}
	}
	return nil, false
}

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
func (c *Cache) expiredKeys() map[string]interface{} {

	c.RLock()
	defer c.RUnlock()
	res := make(map[string]interface{})
	for key, value := range c.Items {
		if time.Now().After(value.Expiration) {
			res[key] = nil
		}
	}

	return res
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys map[string]interface{}) {
	c.Lock()

	defer c.Unlock()
	if len(keys) == 0 {
		return
	}
	for key := range keys {
		delete(c.Items, key)
	}
}
