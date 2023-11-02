package util

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// Create a new source with a seed of current time in unix format to truly randomize the shuffle
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
}

// Create a random float value between min and max
func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Round a float value to a certain precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money
func RandomMoney() float64 {
	return roundFloat(randomFloat(0.0, 1000.99), 2)
}

// RandomCurrency returns a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "ZAR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
