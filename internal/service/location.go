package service

import (
	"errors"
	"strings"
)

const (
	slashSeparator  = "/"
	schemeDelimiter = "://"
)

// ErrInvalidLocation indicates that stored location metadata cannot be parsed.
var ErrInvalidLocation = errors.New("invalid storage location")

// ParseLocation resolves a stored object location to an S3 bucket and key.
func ParseLocation(location, defaultBucket string) (bucket, key string, err error) {
	trimmedLocation := strings.TrimSpace(location)
	if trimmedLocation == "" {
		return "", "", ErrInvalidLocation
	}
	if hasScheme(trimmedLocation) {
		return parseURILocation(trimmedLocation, defaultBucket)
	}
	return parseKeyLocation(trimmedLocation, defaultBucket)
}

func hasScheme(location string) bool {
	index := strings.Index(location, schemeDelimiter)
	return index > 0 && !strings.Contains(location[:index], slashSeparator)
}

func parseKeyLocation(location, defaultBucket string) (bucket, key string, err error) {
	key = strings.TrimLeft(location, slashSeparator)
	if key == "" {
		return "", "", ErrInvalidLocation
	}
	return defaultBucket, key, nil
}

func parseURILocation(location, defaultBucket string) (bucket, key string, err error) {
	scheme, target, _ := strings.Cut(location, schemeDelimiter)
	if !isSupportedScheme(strings.ToLower(scheme)) {
		return "", "", ErrInvalidLocation
	}
	return parseURITarget(target, defaultBucket)
}

func isSupportedScheme(scheme string) bool {
	return scheme == "s3" || scheme == "gs"
}

func parseURITarget(target, defaultBucket string) (bucket, key string, err error) {
	if target == "" {
		return "", "", ErrInvalidLocation
	}
	if strings.HasPrefix(target, slashSeparator) {
		return parseRootKey(target, defaultBucket)
	}
	return parseBucketKey(target, defaultBucket)
}

func parseRootKey(target, defaultBucket string) (bucket, key string, err error) {
	key = strings.TrimLeft(target, slashSeparator)
	if key == "" {
		return "", "", ErrInvalidLocation
	}
	return defaultBucket, key, nil
}

func parseBucketKey(target, defaultBucket string) (bucket, key string, err error) {
	bucket, objectPath, found := strings.Cut(target, slashSeparator)
	if !found {
		return "", "", ErrInvalidLocation
	}
	key = strings.TrimLeft(objectPath, slashSeparator)
	if key == "" {
		return "", "", ErrInvalidLocation
	}
	return resolveBucket(bucket, defaultBucket), key, nil
}

func resolveBucket(bucket, defaultBucket string) string {
	if bucket == "" {
		return defaultBucket
	}
	return bucket
}
