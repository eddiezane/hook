package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/eddiezane/captain-hook/pkg/httptoyaml"

	"github.com/spf13/cobra"
)

var port string

var recordCommand = &cobra.Command{
	Use:   "record",
	Short: "Listens for an incoming webook and saves it",
	Run:   record,
}

func record(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Usage()
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO(eddiezane): Handle http error response

		// TODO(eddiezane): Log body? If so need to clone the readcloser
		log.Printf("method: %s, headers: %v", r.Method, r.Header)

		h, err := httptoyaml.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}

		s, err := h.Dump()
		if err != nil {
			log.Fatal(err)
		}
		// TODO(eddiezane): Would this ever happen?
		if len(s) != 0 {
			// TODO(eddiezane): Path traversal
			// TODO(eddiezane): Create dirs
			err = ioutil.WriteFile(args[0], s, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Printf("starting server on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	recordCommand.Flags().StringVar(&port, "port", "8080", "Port to listen on")
	rootCmd.AddCommand(recordCommand)
}
