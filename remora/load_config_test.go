package remora

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemora_LoadConfig(t *testing.T) {

	r := &Remora{}
	err := r.LoadConfig([]string{"../test/files/success"})
	assert.Nil(t, err)

	assert.NotNil(t, r.Config)
	assert.Equal(t, "5s", r.Config.AcceptableLag)
	assert.Equal(t, "30s", r.Config.CacheTTL)
	assert.Equal(t, 9258, r.Config.HTTPServe)

	assert.NotNil(t, r.Config.MySQL)
	assert.Equal(t, "localhost", r.Config.MySQL.Host)
	assert.Equal(t, 3306, r.Config.MySQL.Port)
	assert.Equal(t, "root", r.Config.MySQL.User)
	assert.Equal(t, "secret", r.Config.MySQL.Pass)

	// Test for failure on non-existentence of config
	err = r.LoadConfig([]string{"../test/files/notthere"})
	assert.NotNil(t, err)
	assert.Regexp(t, regexp.MustCompile("^Config File \"config\" Not Found in .*"), err)

	// Test invalid syntax
	err = r.LoadConfig([]string{"../test/files/fail"})
	assert.NotNil(t, err)
	assert.Regexp(t, regexp.MustCompile("cannot unmarshal"), err)

}

func TestRemora_Run(t *testing.T) {
	// TODO
}
