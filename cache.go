package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

type FileSystemCache struct {
	dir string
}

type Cacheable interface {
	CacheKey() string
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewFileSystemCache(dir string) FileSystemCache {
	projectDir, err := os.Getwd()
	check(err)

	cacheDir := path.Join(projectDir, dir)
	cacheDirExists, err := fileExists(cacheDir)
	check(err)

	if !cacheDirExists {
		err := os.Mkdir(cacheDir, 0777)
		check(err)
	}

	return FileSystemCache{cacheDir}
}

type DemoGetRandNum struct {
	min int
	max int
}

func (d DemoGetRandNum) CacheKey() string {
	return fmt.Sprintf("min-%d_max-%d", d.min, d.max)
}

func (cache FileSystemCache) CachePath(object Cacheable) string {
	targetFile := path.Join(cache.dir, object.CacheKey())

	return targetFile
}

func (cache FileSystemCache) IsCached(object Cacheable) bool {
	targetFile := cache.CachePath(object)
	targetExists, err := fileExists(targetFile)
	check(err)

	return targetExists
}

func (cache FileSystemCache) WriteCache(object Cacheable, serialized []byte) {
	os.WriteFile(cache.CachePath(object), serialized, 0777)
}

func (cache FileSystemCache) ReadCache(object Cacheable) []byte {
	buffer, err := os.ReadFile(cache.CachePath(object))
	check(err)

	return buffer
}

func debugCache() {
	rand.Seed(time.Now().UnixNano())

	cache := NewFileSystemCache("data")

	generator := DemoGetRandNum{1, 5000}

	var value int

	if cache.IsCached(generator) {
		bytes := cache.ReadCache(generator)
		deserialized, err := strconv.Atoi(string(bytes))
		check(err)
		value = deserialized
		fmt.Printf("From cache%d\n", deserialized)
	} else {
		value = rand.Intn(generator.max-generator.min) + generator.min
		fmt.Printf("Random number = %d\n", value)
		serialized := strconv.Itoa(value)
		fmt.Printf("Writing to cache: %s\n", serialized)
		cache.WriteCache(generator, []byte(serialized))
	}

	fmt.Printf("Hello! Cached thingy produced %d\n", value)
}
