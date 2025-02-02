package helpers

import "math/rand"

type Dice struct {
	Value int
	Emoji string
}

var DiceEmojiMap = map[int]string{
	1: "1️⃣",
	2: "2️⃣",
	3: "3️⃣",
	4: "4️⃣",
	5: "5️⃣",
	6: "6️⃣",
}

func GetRandomD6() Dice {
	min := 1
	max := 6
	random := int32(rand.Intn(max-min) + min)
	return Dice{int(random), DiceEmojiMap[int(random)]}
}
