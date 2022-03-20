package main

import (
	"os"
	"path"
)

type FileSystemCache struct {
	dir string
}

type Cacheable interface {
	CacheKey() string
}

// Caching functionality concerned with the filesystem

func NewFileSystemCache(dir string) FileSystemCache {
	cacheDir, err := JoinProjectPath(dir)

	check(err)
	check(EnsureDirExists(cacheDir))

	return FileSystemCache{cacheDir}
}

func (cache FileSystemCache) CachePath(object Cacheable) string {
	targetFile := path.Join(cache.dir, object.CacheKey())

	return targetFile
}

func (cache FileSystemCache) IsCached(object Cacheable) bool {
	targetFile := cache.CachePath(object)
	pathType, err := checkPathType(targetFile)
	check(err)

	return pathType == IsFile
}

func (cache FileSystemCache) WriteCache(object Cacheable, serialized []byte) {
	os.WriteFile(cache.CachePath(object), serialized, 0777)
}

func (cache FileSystemCache) ReadCache(object Cacheable) []byte {
	buffer, err := os.ReadFile(cache.CachePath(object))
	check(err)

	return buffer
}

// Caching functionality concerned with JSON (de)serialization

type JSONSerializable interface{}

type WithCacheCallback func() ([]byte, error)

func (cache FileSystemCache) WithRawCache(subject Cacheable, callback WithCacheCallback) ([]byte, error) {
	if cache.IsCached(subject) {
		return cache.ReadCache(subject), nil
	}

	liveResult, err := callback()

	if err != nil {
		return nil, err
	}

	cache.WriteCache(subject, liveResult)
	return liveResult, nil
}
