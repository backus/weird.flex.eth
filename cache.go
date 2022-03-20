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
