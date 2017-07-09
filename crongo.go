// Copyright (C) 2017  Karl Cordes

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package crongo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// Job contains details about a command that was executed
type Job struct {
	Start    time.Time
	End      time.Time
	Pid      int
	Command  string
	Output   string
	Status   string
	Success  bool
	Hostname string
}

//@TODO sort by job age
type ByAge []Job

func (j ByAge) Len() int {
	return len(j)
}

func (j ByAge) Less(a, b int) bool {
	return false
}

// Run executes a command and captures its output
// Returns a Job
func Run(command string, args string) Job {
	var j Job

	cmd := exec.Command(command, args)
	j.Command = command
	j.Start = time.Now()

	out, err := cmd.Output()
	j.Output = string(out)
	j.End = time.Now()
	// command returned non-zero exit
	if err != nil {
		j.Status = err.Error()
		log.Printf("ERROR %s\n", err)
	}

	j.Success = cmd.ProcessState.Success()
	j.Pid = cmd.ProcessState.Pid()

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname, somehow..")
		hostname = "UNKNOWN"
	}
	j.Hostname = hostname

	return j
}

func PostJob(j Job, config ClientConfig) {

	b, _ := json.Marshal(j)

	endpoint := config.Server + "/api/jobs"

	_, err := http.Post(endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("ERROR: Unable to send job\n")
		log.Panic(err)
	}
}

func ListJobs() []Job {
	//TODO read from config
	files, err := ioutil.ReadDir("./output")
	if err != nil {
		log.Printf("ERROR unable to read jobs output directory: %s\n", err)
	}
	jobs := make([]Job, len(files))
	for _, file := range files {
		job, err := ReadJob(file.Name())
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func ReadJob(filename string) (Job, error) {
	var j Job
	//@TODO read from Config
	b, err := ioutil.ReadFile("./output/" + filename)
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
