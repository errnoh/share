// Copyright 2013 errnoh. All rights reserved.
// Use of this source code is governed by
// MIT License that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	path = flag.String("path", "/share", "share file under $HOST/$path")
	port = flag.String("port", "8080", "port")

	verbose = flag.Bool("verbose", false, "")

	cert = flag.String("cert", "", "cert.pem file")
	key  = flag.String("key", "", "key.pem file")
)

func main() {
	flag.Parse()

	var sharing bool
	for _, file := range flag.Args() {
		fi, err := os.Stat(file)
		if err != nil {
			log.Printf("%s: %s", file, err.Error())
			continue
		}
		if fi.Mode()&os.ModeType != 0 {
			log.Printf("%s: not a regular file")
			continue
		}

		path := filepath.Join(*path, filepath.Base(file))
		if *verbose {
			log.Printf("sharing: %s", path)
		}
		sharing = true
		http.Handle(path, FileHandler{file})

	}

	if !sharing {
		log.Println("No files to share")
		return
	}

	if *cert != "" {
		log.Fatal(http.ListenAndServeTLS(":"+*port, *cert, *key, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+*port, nil))
	}
}

type FileHandler struct {
	file string
}

func (fh FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		log.Printf("Request  : %s -> %s", r.RemoteAddr, r.URL)
	}

	if r.Body != nil {
		r.Body.Close()
	}

	f, err := os.Open(fh.file)
	if err != nil {
		if *verbose {
			log.Printf("Failed   : %s -> %s [%s]", r.RemoteAddr, r.URL, err.Error())
		}
		http.NotFound(w, r)
		return
	}
	start := time.Now()
	written, err := io.Copy(w, f)
	if *verbose {
		if err != nil {
			log.Printf("Failed   : %s -> %s (wrote %d bytes) after %s [%s]", r.RemoteAddr, r.URL, written, time.Since(start), err.Error())
		} else {
			log.Printf("Completed: %s -> %s (wrote %d bytes) in %s", r.RemoteAddr, r.URL, written, time.Since(start))
		}
	}

	f.Close()
}
