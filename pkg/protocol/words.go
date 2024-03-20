package protocol

import "math/rand"

var wordsOfWisdom = []string{
	"The only true wisdom is in knowing you know nothing.",
	"The only way to do great work is to love what you do.",
	"Life is what happens when you're busy making other plans.",
	"Get busy living or get busy dying.",
	"You only live once, but if you do it right, once is enough.",
	"Believe you can and you're halfway there.",
	"Turn your wounds into wisdom.",
	"Change the world by being yourself.",
	"Every moment is a fresh beginning.",
	"Never regret anything that made you smile.",
	"Die with memories, not dreams.",
	"Aspire to inspire before we expire.",
}

func getRandomWisdom() string {
	return wordsOfWisdom[rand.Intn(len(wordsOfWisdom))]
}
