package foo

import (
	"fmt"
	"github.com/lsgrep/sak/ucfg"
	"github.com/lsgrep/sak/ucfg/example/bar"
	"github.com/spf13/viper"
)

type Foo struct {
	Url string `mapstructure:"url"`
}

var foo Foo

func init() {
	ucfg.Register("example.foo", initFoo)
}

func initFoo(vp *viper.Viper) {
	err := vp.Unmarshal(&foo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("init Foo done")

	url := vp.GetString("url")
	if url == "" {
		panic("foo.url not found")
	}
}

func Work() {
	bar.Work()
}
