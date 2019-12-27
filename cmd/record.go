package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/eddiezane/hook/pkg/hook"

	"github.com/spf13/cobra"
)

var (
	// Flags
	port   string
	base64 []string

	recordCommand = &cobra.Command{
		Use:     "record",
		Short:   "Listens for an incoming webook and saves it",
		Long:    "records starts up a local HTTP server and saves a request made against it into a YAML serialization at the provided path",
		Example: "hook record --port 9000 path/to/webhook.yml",
		RunE:    record,
	}
)

type recorder struct {
	mu sync.Mutex
	f  *os.File

	opts []hook.Option
}

func newRecorder(path string, opts ...hook.Option) (*recorder, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &recorder{
		f:    f,
		opts: opts,
	}, nil
}

func (r *recorder) close() error {
	return r.f.Close()
}

func (r *recorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO(eddiezane): Handle http error response

	// TODO(eddiezane): Log body? If so need to clone the readcloser
	log.Printf("method: %s, headers: %v, params: %v", req.Method, req.Header, req.URL.Query())

	h, err := hook.NewFromRequest(req, r.opts...)
	if err != nil {
		log.Fatal(err)
	}

	s, err := h.Dump()
	if err != nil {
		log.Println("error dumping hook:", err, h)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO(eddiezane): Would this ever happen?
	if len(s) != 0 {
		// TODO(eddiezane): Create dirs
		r.mu.Lock()
		defer r.mu.Unlock()
		fw := bufio.NewWriter(r.f)
		if fi, err := r.f.Stat(); err == nil && fi.Size() > 0 {
			// If file has data in it already, append doc separator.
			if _, err := fw.Write([]byte("---\n")); err != nil {
				log.Println("error writing doc separator:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if _, err := fw.Write(s); err != nil {
			log.Println("error writing file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fw.Flush()
	}

	w.WriteHeader(http.StatusOK)
}

func record(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments provided. expected %d", 1)
	}

	var opts []hook.Option
	b64t := &hook.Base64Transformer{}
	for _, f := range base64 {
		opts = append(opts, hook.DecodeOption(b64t, f))
	}

	r, err := newRecorder(args[0], opts...)
	if err != nil {
		return err
	}
	defer r.close()

	log.Printf("starting server on port %s", port)
	return http.ListenAndServe(":"+port, r)
}

func init() {
	recordCommand.Flags().StringVar(&port, "port", "8080", "Port to listen on")
	recordCommand.Flags().StringArrayVar(&base64, "base64", nil, "comma separated list of fields to base64 decode")
	rootCmd.AddCommand(recordCommand)
}
