package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/eddiezane/captain-hook/pkg/httptoyaml"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fireCommand)

	// TODO(eddiezane): Slurp up all hooks here?
}

var fireCommand = &cobra.Command{
	Use:   "fire",
	Short: "Fires the selected webhook at a given url",
	Run:   fire,
}

func fire(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Usage()
		os.Exit(1)
	}
	hookName := args[0]
	bs, err := ioutil.ReadFile(hookName)
	if err != nil {
		panic(err)
	}
	h, err := httptoyaml.Slurp(bs)
	if err != nil {
		panic(err)
	}
	r, err := httptoyaml.Unmarshal(h)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(args[1])
	if err != nil {
		panic(err)
	}
	r.URL = u
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
