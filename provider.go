package provider

import (
	"learngo/provider/settings"
	"math/rand"
	"time"
)

var CodeMin int
var CodeMax int

var idRequest int

func Init() {
	CodeMin, CodeMax = settings.DoCheckAndGetValues()
	rand.Seed(time.Now().Unix())
}

func RequestCode(id int) int {
	if !settings.IsCheckDone() {
		panic("No check was done before! Call provider.Init()!")
	}

	idRequest = id
	code := random(CodeMin, CodeMax)

	return code
}

func random(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
