// Copyright (C) 2017  Karl Cordes

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details./ You should have received a copy of the GNU General Public License
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses

package crongo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Api")
}

type JobsHandler struct {
	Config ServerConfig
}

func (jh JobsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var j Job
		err = json.Unmarshal(body, &j)
		if err != nil {
			log.Printf("ERROR: Unable to decode JSON")
		}
		jh.SaveJob(j)
	}

	if r.Method == "GET" {
		jobs := jh.ListJobs()
		b, err := json.Marshal(jobs)
		if err != nil {
			log.Printf("ERROR: Unable to marshal Jobs array: %s\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		w.Write(b)

	}
}

// Write the job JSON to the output directory
func (jh JobsHandler) SaveJob(j Job) {

	now := time.Now()
	nowFormatted := now.Format("20060102T15_04_05")

	filename := nowFormatted + "_" + j.Hostname + "_" + path.Base(j.Command) + ".json"
	fullPath := path.Join(jh.Config.OutputDir, filename)

	b, _ := json.Marshal(j)
	ioutil.WriteFile(fullPath, b, 0666)
}

func (jh JobsHandler) ListJobs() []Job {
	files, err := ioutil.ReadDir(jh.Config.OutputDir)
	if err != nil {
		log.Printf("ERROR unable to read jobs output directory: %s\n", err)
	}
	jobs := make([]Job, len(files))

	for i, file := range files {
		job, err := jh.ReadJob(file.Name())
		if err != nil {
			log.Println(err)
			continue
		}
		jobs[i] = job
	}
	return jobs
}

func (jh JobsHandler) ReadJob(filename string) (Job, error) {
	var j Job
	f := path.Join(jh.Config.OutputDir, filename)
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("ERROR unable to read job file: %s\n", err)
		return j, err
	}
	err = json.Unmarshal(b, &j)
	if err != nil {
		log.Printf("ERROR unable to read job file json: %s\n", err)
		return j, err
	}
	return j, nil
}

func ServeAPI(config ServerConfig) {
	var jh JobsHandler
	jh.Config = config
	mux := http.NewServeMux()
	mux.Handle("/api/jobs", jh)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome. This will eventually display success/failures of jobs")
	})

	log.Printf("Listening on %s\n", jh.Config.ListenAddress)
	log.Fatal(http.ListenAndServe(jh.Config.ListenAddress, mux))
}
