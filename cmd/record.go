package cmd

import (
	"log"
	"net/http"

	"github.com/eddiezane/captain-hook/pkg/httptoyaml"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCommand)
}

var recordCommand = &cobra.Command{
	Use:   "record",
	Short: "Listens for an incoming webook and saves it",
	Run:   record,
}

func record(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		h, err := httptoyaml.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/yaml")

		s, err := h.Dump()
		if err != nil {
			log.Fatal(err)
		}

		if len(s) != 0 {
			_, err = w.Write(s)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
