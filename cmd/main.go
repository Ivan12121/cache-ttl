package main

import (
	"cache-ttl/cache"
	"fmt"
	"time"
)

func main() {
	cache := cache.NewCache(2 * time.Second)
	defer cache.Stop()

	cache.Set("user", "Alex", 5 * time.Second)
	
	value, ok := cache.Get("user")
	fmt.Println("Сразу после записи:", value, ok)

	time.Sleep(3 * time.Second)

	value, ok = cache.Get("user")
	fmt.Println("Через 3 секунды:", value, ok)

	time.Sleep(3 * time.Second)

	value, ok = cache.Get("user")
	fmt.Println("Через 6 секунды:", value, ok)
}