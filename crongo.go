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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"time"
)

type Job struct {
	Start   time.Time
	End     time.Time
	Pid     int
	Command string
	Output  string
	Status  string
	Success bool
}

type Config struct {
	Server url.URL
	Token  string
}

func Run(command string) Job {
	var j Job

	cmd := exec.Command(command)
	j.Command = command
	j.Start = time.Now()

	out, err := cmd.Output()
	j.Output = string(out)
	j.End = time.Now()

	j.Success = cmd.ProcessState.Success()
	j.Pid = cmd.ProcessState.Pid()

	//command returned non-zero exit
	//Unsure exactly how to handle this
	if err != nil {
		j.Status = err.Error()
		log.Printf("ERROR: %s\n", err)
	} else {
		log.Printf("OK\n")
	}

	return j
}

func ReadConfig(confPath string) Config {
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Printf("ERROR: Unable to read config: %s", confPath)
		log.Panic()
	}
	var c Config
	err = json.Unmarshal(b, &c)

	if err != nil {
		log.Printf("ERROR: Config is invalid JSON")
		log.Printf("%s\n", err)
		log.Panic()
	}
	return c
}
