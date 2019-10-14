package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/eddiezane/captain-hook/pkg/hook"

	"github.com/spf13/cobra"
)

var port string

var recordCommand = &cobra.Command{
	Use:   "record",
	Short: "Listens for an incoming webook and saves it",
	RunE:  record,
}

func record(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		cmd.Usage()
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO(eddiezane): Handle http error response

		// TODO(eddiezane): Log body? If so need to clone the readcloser
		log.Printf("method: %s, headers: %v", r.Method, r.Header)

		h, err := hook.NewFromRequest(r)
		if err != nil {
			log.Fatal(err)
		}

		s, err := h.Dump()
		if err != nil {
			log.Fatal(err)
		}
		// TODO(eddiezane): Would this ever happen?
		if len(s) != 0 {
			// TODO(eddiezane): Create dirs
			if err := ioutil.WriteFile(args[0], s, 0644); err != nil {
				log.Fatal(err)
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Printf("starting server on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}

func init() {
	recordCommand.Flags().StringVar(&port, "port", "8080", "Port to listen on")
	rootCmd.AddCommand(recordCommand)
}
