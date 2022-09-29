package toolkit

import "crypto/rand"

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"
// Tools is the type used to instantiate this module.Any variable of this type will have ccess
// to all the methods with the reciver *Tools
type Tools struct {}

func (t *Tools) RandomString(n int) string {
	s, r := make([]rune,n),[]rune(randomStringSource)
	for i := range s {
		p,_ := rand.Prime(rand.Reader,len(r))
		x, y := p.Uint64(),uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}