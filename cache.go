package main

import (
	"encoding/json"
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

type WithRawCacheCallback func() ([]byte, error)
type WithJSONCacheCallback[ResultType JSONSerializable] func() (ResultType, error)

func (cache FileSystemCache) WithRawCache(subject Cacheable, callback WithRawCacheCallback) ([]byte, error) {
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

func WithJSONCache[Deserialized JSONSerializable](cache FileSystemCache, subject Cacheable, callback WithJSONCacheCallback[Deserialized]) (Deserialized, error) {
	var deserialized Deserialized

	if cache.IsCached(subject) {
		logger.Debug("Cache hit (%s)", subject.CacheKey())

		data := cache.ReadCache(subject)
		json.Unmarshal(data, &deserialized)
		return deserialized, nil
	}

	deserialized, err := callback()

	if err != nil {
		return deserialized, err
	}

	serialized, err := json.MarshalIndent(deserialized, "", "  ")
	if err != nil {
		return deserialized, err
	}

	cache.WriteCache(subject, serialized)

	return deserialized, nil
}
