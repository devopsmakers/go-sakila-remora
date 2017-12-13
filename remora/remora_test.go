package remora

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemora_LoadConfig(t *testing.T) {

	r := &Remora{}
	err := r.LoadConfig([]string{"../test/files/success"}, "mysql")
	assert.Nil(t, err)

	assert.NotNil(t, r.Config)
	assert.Equal(t, 5, r.Config.AcceptableLag)
	assert.Equal(t, "5s", r.Config.CacheTTL)
	assert.Equal(t, 9258, r.Config.HTTPServe)
	assert.Equal(t, false, r.Config.Maintenance)

	assert.NotNil(t, r.Config.Service)
	assert.Equal(t, "localhost", r.Config.Service.Host)
	assert.Equal(t, 3306, r.Config.Service.Port)
	assert.Equal(t, "root", r.Config.Service.User)
	assert.Equal(t, "secret", r.Config.Service.Pass)
	assert.Equal(t, false, r.Config.Service.Ssl)

	// Test for failure on non-existentence of config
	err = r.LoadConfig([]string{"../test/files/notthere"}, "mysql")
	assert.NotNil(t, err)
	assert.Regexp(t, regexp.MustCompile("^Config File \"config\" Not Found in .*"), err)

	// Test invalid syntax
	err = r.LoadConfig([]string{"../test/files/fail"}, "mysql")
	assert.NotNil(t, err)
	assert.Regexp(t, regexp.MustCompile("cannot unmarshal"), err)

}

func TestRemora_Serve(t *testing.T) {
	// TODO
}

func TestRemora_ServeHTTP(t *testing.T) {
	// TODO
}
