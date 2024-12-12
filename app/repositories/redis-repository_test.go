package repositories

import (
	"creditlimit-connector/app/models"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_Save(t *testing.T) {
	// Initialize a mock Redis client and repository
	mockRedisClient, mockClient := redismock.NewClientMock()
	repo := &RedisRepositoryImp{
		redisClient: mockRedisClient,
	}

	// Define test data
	key := "user:123"
	value := models.FnGetOnshoreBusinessDateModel{
		Date:              "2021-01-01",
		IsOperationalDate: true,
	}
	ttlInSecond := 60

	valueBytes, _ := json.Marshal(value)

	// Setup the mock to expect the Set command and return nil error (success)
	mockClient.ExpectSet(key, valueBytes, time.Duration(ttlInSecond)*time.Second).SetVal("OK")

	// Call the Save method
	err := repo.Save(key, value, ttlInSecond)

	// Assertions
	assert.NoError(t, err) // Assert no error occurred
}

func TestRedisRepository_Save_MarshalError(t *testing.T) {
	// Initialize a mock Redis client and repository
	mockRedisClient, _ := redismock.NewClientMock()
	repo := &RedisRepositoryImp{
		redisClient: mockRedisClient,
	}

	// Test case where Marshal fails (invalid value)
	key := "user:123"
	value := make(chan int) // Channels cannot be marshaled
	ttlInSecond := 60

	// Call the Save method, which should return an error due to marshal failure
	err := repo.Save(key, value, ttlInSecond)

	// Assertions
	assert.Error(t, err)                                                // Assert an error occurred
	assert.Contains(t, err.Error(), "json: unsupported type: chan int") // Specific error message
}

func TestRedisRepository_Find(t *testing.T) {
	// Initialize a mock Redis client and repository
	mockRedisClient, mockClient := redismock.NewClientMock()
	repo := &RedisRepositoryImp{
		redisClient: mockRedisClient,
	}

	// Define test data
	key := "user:123"
	expectedValue := models.FnGetOnshoreBusinessDateModel{
		Date:              "2021-01-01",
		IsOperationalDate: true,
	}

	expectedValueBytes, _ := json.Marshal(expectedValue)

	// Setup the mock to expect the Get command and return a serialized JSON string
	mockClient.ExpectGet(key).SetVal(string(expectedValueBytes))

	// Create a variable to hold the result
	var result models.FnGetOnshoreBusinessDateModel

	// Call the Find method
	err := repo.Find(key, &result)

	// Assertions
	assert.NoError(t, err)                                                     // Assert no error occurred
	assert.Equal(t, expectedValue.Date, result.Date)                           // Assert the result is as expected
	assert.Equal(t, expectedValue.IsOperationalDate, result.IsOperationalDate) // Assert the result is as expected

}

func TestRedisRepository_Find_Error(t *testing.T) {
	// Initialize a mock Redis client and repository
	mockRedisClient, mockClient := redismock.NewClientMock()
	repo := &RedisRepositoryImp{
		redisClient: mockRedisClient,
	}

	// Define test data
	key := "user:123"

	// Setup the mock to simulate a Redis error
	mockClient.ExpectGet(key).SetErr(errors.New("Redis error"))

	// Create a variable to hold the result
	var result models.FnGetOnshoreBusinessDateModel

	// Call the Find method
	err := repo.Find(key, result)

	// Assertions
	assert.Error(t, err)                        // Assert an error occurred
	assert.Equal(t, "Redis error", err.Error()) // Assert that the error message is correct
	assert.Empty(t, result)                     // Assert that the value returned is empty
}
