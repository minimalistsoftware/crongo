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
	"flag"
	"fmt"
	"github.com/minimalistsoftware/crongo"
	"log"
	"strings"
)

var cmd = flag.String("cmd", "", "command to run")

//var args = flag.String("args", "", "arguments to the command")
var conf = flag.String("conf", "/etc/crongo.json", "path to crongo.json config file")

// Example 3: A user-defined flag type, a slice of strings
type argsFlag []string

var args argsFlag

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (a *argsFlag) String() string {
	return fmt.Sprint(*a)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a space seperated list, so we split it.
func (a *argsFlag) Set(value string) error {

	for _, arg := range strings.Fields(value) {
		*a = append(*a, arg)
	}
	return nil
}

func init() {
	flag.Var(&args, "arg", "Arguments to pass to the command")
}

func main() {
	flag.Parse()

	if *cmd == "" {
		log.Fatal("ERROR: cmd is empty")
	}

	config := crongo.ReadClientConfig(*conf)

	job := crongo.Run(*cmd, args)
	crongo.PostJob(job, config)

}
