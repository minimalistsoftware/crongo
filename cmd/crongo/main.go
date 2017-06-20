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
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/minimalistsoftware/crongo"
	"log"
	"net/http"
)

var toRun = flag.String("cmd", "", "command to run")
var conf = flag.String("conf", "/etc/crongo.json", "path to crongo.json config file")

func main() {
	flag.Parse()

	if *toRun == "" {
		log.Fatal("ERROR: cmd is empty")
	}

	log.Println(*conf)

	config := crongo.ReadConfig(*conf)

	j := crongo.Run(*toRun)
	b, _ := json.Marshal(j)
	log.Printf("%s", b)

	endpoint := config.Server.String() + "/api/jobs"

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("ERROR: Unable to send job\n")
		log.Panic(err)
	}

	log.Println(resp)
}
