# Cachey

Cachey is a simple and flexible caching library for Go, inspired by Laravel's cache library. It provides a straightforward interface for caching data with support for multiple providers.

## Features

- Easy-to-use API for caching data.
- Memory store provider with TTL support.
- Inspired by Laravel's Cache library.

## Installation

To install Cachey, you can use `go get`:

```bash
go get github.com/codemaestro64/cachey
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/codemaestro64/cachey"
)

func main() {
    // Create a new cache instance with the memory provider
    cache, err := cachey.New(cachey.MemoryStore)
    if err != nil {
        panic(err)
    }

    // Set a value in the cache
    cache.Put("key", "value", 5*time.Second)

    // Retrieve the value from the cache
    value := cache.Get("key")
    fmt.Println(value) // Output: value

    // Wait for expiration
    time.Sleep(6 * time.Second)
    value = cache.Get("key")
    fmt.Println(value) // Output: <nil>
}
```

### Cache Methods

Cachey provides a range of methods for interacting with the cache:

- **Has(key string) bool**: Checks if a value exists in the cache for the given key.
- **Get(key string) any**: Retrieves the value associated with the given key from the cache.
- **GetOrDefault(key string, defaultFunc func() any) any**: Retrieves the value for the specified key, or returns the result of `defaultFunc` if the key does not exist.
- **Remember(key string, duration time.Duration, rememberFunc func() any) any**: Retrieves the value for the specified key, or calls `rememberFunc` to generate the value and store it in the cache.
- **RememberForever(key string, rememberFunc func() any) any**: Similar to `Remember`, but stores the value indefinitely.
- **Pull(key string) any**: Retrieves the value for the specified key and removes it from the cache.
- **Put(key string, data any, duration time.Duration)**: Stores the given data in the cache with the specified duration.
- **Forever(key string, data any)**: Stores the given data indefinitely.
- **Add(key string, data any, duration time.Duration)**: Stores the given data only if the key does not already exist.
- **Forget(key string)**: Removes the value associated with the specified key from the cache.
- **Flush()**: Removes all values from the cache.

### Registering Additional Providers

You can register additional cache providers by using the `RegisterProvider` function. The following providers are planned for future implementation:

- Redis
- File-based caching
- Memcached

To register a provider, use the following syntax:

```go
err := cachey.RegisterProvider("providerName", providerConstructor)
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Todo

- Implement Redis cache provider.
- Implement file-based cache provider.
- Implement Memcached cache provider.

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## Contact

For questions or feedback, please open an issue in the repository.
```