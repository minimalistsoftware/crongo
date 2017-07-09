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

var OutputDir = "./output/"

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Api")
}

func JobsHandler(w http.ResponseWriter, r *http.Request) {
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
		SaveJob(j)
	}

	if r.Method == "GET" {
		jobs := ListJobs()
		//TODO handle this error
		b, _ := json.Marshal(jobs)
		w.Write(b)

	}
}

func ServeAPI() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/jobs", JobsHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome. This will eventually display success/failures of jobs")
	})

	//@TODO read listen address from config
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Write the job JSON to the output directory
func SaveJob(j Job) {

	now := time.Now()
	nowFormatted := now.Format("20060102T15_04_05")
	log.Println(nowFormatted)

	filename := nowFormatted + "_" + j.Hostname + "_" + path.Base(j.Command) + ".json"
	fullPath := path.Join(OutputDir, filename)

	b, _ := json.Marshal(j)
	ioutil.WriteFile(fullPath, b, 0666)

}
