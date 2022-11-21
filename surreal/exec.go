package surreal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/buonotti/bus-stats-api/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/*
Starts the surrealDb database in a goroutine. The executable is searched in ./bin. If the database is not reachable the
api tries three times to connect to it. If the database schema is not created the cli also creates the table
*/
func Exec() {
	surrealExe := viper.GetString("database.executable")

	if runtime.GOOS == "windows" {
		surrealExe = surrealExe + ".exe"
	}
	log.Debug(fmt.Sprintf("Db executable is: %s", surrealExe))
	cmd := exec.Command(surrealExe)
	mode := viper.GetString(config.Get("database.{env}.mode"))
	user := viper.GetString(config.Get("database.{env}.user"))
	pass := viper.GetString(config.Get("database.{env}.pass"))
	cmd.Args = []string{surrealExe, "start", "--user", user, "--pass", pass, mode}
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()

	for i := 3; i >= 1 && !isDbOnline(); i-- {
		log.Warn(fmt.Sprintf("database seems not to be online. retrying to connect. tries left: %d", i))
		cmd := exec.Command("sleep", "2")
		err := cmd.Run()
		if err != nil {
			log.Error("error while waiting for db")
			os.Exit(1)
		}
	}

	if !isDbOnline() {
		log.Error("could not read database")
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("started database with authentication in %s", mode))
	isDefined := viper.GetBool("database.generated")
	if !isDefined {
		err := ScaffoldDB()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("generated database tables")

		viper.Set("database.generated", true)
		err = viper.WriteConfig()
		if err != nil {
			log.Error(err)
		}
	}
}

func isDbOnline() bool {
	return PingDB() == nil
}
