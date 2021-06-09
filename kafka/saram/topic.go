package saram

import "fmt"

func GetTopics(template string, n int) []string {
	topics := make([]string, 0, n)

	for i := 1; i <= n; i++ {
		topics = append(topics, fmt.Sprintf("%s_%d", template, i))
	}

	return topics
}
