package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"os"
	"testing"
)

// resetEnv resets the environment and viper state between tests
func resetEnv() {
	viper.Reset()
	Config = nil
}

func TestConfig(t *testing.T) {
	Convey("Config Tests", t, func() {

		Convey("When no environment variables are set, it should use defaults", func() {
			resetEnv()
			os.Clearenv()

			config := GetConfig()

			So(config.Port, ShouldEqual, "8080")
			So(config.Host, ShouldEqual, "0.0.0.0")
			So(config.ChromeSetup, ShouldBeEmpty)
		})

		Convey("When environment variables are set, they should override defaults", func() {
			resetEnv()

			t.Setenv("PORT", "9090")
			t.Setenv("HOST", "local")
			t.Setenv("CHROME_SETUP", "chrome")

			config := GetConfig()

			So(config.Port, ShouldEqual, "9090")
			So(config.Host, ShouldEqual, "local")
			So(config.ChromeSetup, ShouldEqual, "chrome")
		})
	})
}
