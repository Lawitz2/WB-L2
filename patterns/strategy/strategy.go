package main

import "fmt"

/*
Стратегия применяется при наличии большого количества схожих страктов, поведение которых должно отличаться
*/

const cacheCap = 2

type cache struct {
	evictionStrategy IcacheEvictionStrategy
	storage          map[string]string
	len              int
	cap              int
}

func initCache(e IcacheEvictionStrategy) *cache {
	storage := make(map[string]string)
	return &cache{
		evictionStrategy: e,
		storage:          storage,
		len:              0,
		cap:              cacheCap,
	}
}

func (c *cache) setEvictionStrategy(es IcacheEvictionStrategy) {
	c.evictionStrategy = es
}

func (c *cache) add(key, val string) {
	if c.len == c.cap {
		c.evict()
	}
	c.len++
	c.storage[key] = val
}

func (c *cache) del(key string) {
	delete(c.storage, key)
}

func (c *cache) evict() {
	c.evictionStrategy.evict(c)
	c.len--
}

type IcacheEvictionStrategy interface {
	evict(c *cache)
}

type fifoStrat struct{}

func (f fifoStrat) evict(c *cache) {
	//determines what to delete, calls c.del(key)
	fmt.Println("used First In First Out to evict something from cache")
}

type lruStrat struct{}

func (l lruStrat) evict(c *cache) {
	//determines what to delete, calls c.del(key)
	fmt.Println("used Least Recently Used strategy to evict something from cache")
}

type lfuStrat struct{}

func (l lfuStrat) evict(c *cache) {
	//determines what to delete, calls c.del(key)
	fmt.Println("used Least Frequently Used strategy to evict something from cache")
}

func main() {
	lru := lruStrat{}
	lfu := lfuStrat{}
	fifo := fifoStrat{}

	c := initCache(fifo)

	c.add("1", "a")
	c.add("2", "b")
	c.add("3", "v")
	c.add("4", "g")
	c.setEvictionStrategy(lru)
	c.add("5", "d")
	c.setEvictionStrategy(lfu)
	c.add("6", "e")
}
