package app

import (
	"fmt"
	"github.com/spf13/viper"
	"lab/pkg"
	"strconv"
	"strings"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	var (
		err error
	)
	err = initConfigs()
	if err != nil {
		panic(err)
	}

	fsm := pkg.NewFSM(viper.GetInt("init.n"), createArray(viper.GetString("init.phi")), createArray(viper.GetString("init.psi")))
	fsm.Init()

	connectivityComponents := fsm.GetConnectivityComponents()
	for _, v := range connectivityComponents {
		fmt.Println("Connected components:", v)
	}
	fmt.Printf("Number of components: %d\n", len(connectivityComponents))

	strongConnectivityComponents := fsm.GetStrongConnectivityComponents()
	for _, v := range strongConnectivityComponents {
		fmt.Println("Strong connected components:", v)
	}
	fmt.Printf("Number of strong components: %d\n", len(strongConnectivityComponents))

	equivalenceClass := pkg.NewClassEquivalence()
	equivalenceClasses := equivalenceClass.GetEquivalenceClasses(fsm.PhiTable, fsm.PsiTable)
	for i, class := range equivalenceClasses {
		fmt.Printf("Equivalence class %d\n", i)
		for idx, v := range class {
			fmt.Printf("Subclass %d: %s\n", idx, v.Class)
		}
	}
}

func initConfigs() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}

func createArray(str string) []int {
	tmp := strings.Split(str, "")
	result := make([]int, len(tmp))

	for i, v := range tmp {
		result[i], _ = strconv.Atoi(v)
	}
	return result
}
