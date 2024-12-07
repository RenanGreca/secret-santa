// Package assert provides assertions checking correct values at runtime.
// Please note that a failed assertion will cause the program to exit.
package assert

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func stringify(item any) string {
	if item == nil {
		return "nil"
	}
	switch t := item.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case int:
		return fmt.Sprintf("%d", item)
	case bool:
		return fmt.Sprintf("%t", item)
	case time.Time:
		return item.(time.Time).Format(time.RFC3339)
	default:
		d, err := json.Marshal(item)
		if err == nil {
			return string(d)
		}
	}
	return fmt.Sprintf("%s", item)
}

// True checks if the given truth value is true.
func True(truth bool, msg string) {
	if !truth {
		log.Println("Assert#True failed")
		log.Fatalln(msg)
	}
}

// Nil checks if the given item is nil.
func Nil(item any, msg string) {
	if item != nil {
		log.Println("Assert#Nil failed")
		log.Printf("%s\n", stringify(item))
		log.Fatalln(msg)
	}
}

// NotNil checks if the given item is not nil.
func NotNil(item any, msg string) {
	if item == nil {
		log.Println("Assert#NotNil failed")
		log.Fatalln(msg)
	}
}

// NoError checks if the given error is nil.
func NoError(err error, msg string, data ...any) {
	if err != nil {
		log.Println("Assert#NoError failed")
		log.Printf("%s\n", stringify(err))
		log.Printf("%s\n", err.Error())
		log.Fatalf(msg, data...)
	}
}

func NotZero(n int, msg string, data ...any) {
	if n == 0 {
		log.Printf("Assert#NotZero failed")
		log.Fatalf(msg, data...)
	}
}

func Equal(a, b any, msg string, data ...any) {
	if a != b {
		log.Println("Assert#Equal failed")
		log.Fatalf(msg, data...)
	}
}

func LogError(msg string, err error) {
	log.Printf(msg, stringify(err))
}

func LogDebug(msg string, data ...any) {
	if os.Getenv("DEBUG") == "true" {
		log.Printf(msg, data...)
	}
}
