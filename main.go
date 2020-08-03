/*
 * Fd
 * Copyright (C) 2020  Nicolò Santamaria
 *
 * Fd is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Fd is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	pattern  string
	print0   bool
	useRegex bool
	re       *regexp.Regexp
)

func usage() {
	fmt.Printf(`fd - Find all files mathing a pattern.
Fd recursively finds all the files whose names match a pattern provided in input.
Usage:
    %s [options] [pattern] [path]
Options:
    -r    Use a regex instead of the shell file name pattern.
    -0    Separate search results by the null character (instead of newlines). Useful for piping results to 'xargs'.
`, os.Args[0])
}

func die(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func exists(fname string) (bool, error) {
	_, err := os.Stat(fname)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func matches(fname string) (bool, error) {
	if useRegex {
		return re.MatchString(fname), nil
	}
	return filepath.Match(pattern, fname)
}

func checkFile(fpath string, info os.FileInfo, err error) error {
	isMatch, err := matches(info.Name())
	if err != nil {
		return err
	}
	if isMatch {
		fmt.Println(fpath)
	}
	return nil
}

func main() {
	var root string

	switch argc := flag.NArg(); {
	case argc == 1:
		pattern = flag.Arg(0)
		root = "."
	case argc == 2:
		pattern = flag.Arg(0)
		root = flag.Arg(1)
	default:
		flag.Usage()
		os.Exit(1)
	}

	ok, err := exists(root)
	if err != nil {
		die(err)
	}

	if useRegex {
		var err error
		if re, err = regexp.Compile(pattern); err != nil {
			die(err)
		}
	}

	if ok {
		if err := filepath.Walk(root, checkFile); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("'%s': No such file or directory.\n", root)
	}
}

func init() {
	flag.BoolVar(&useRegex, "r", false, "Use a regex instead of the shell file name pattern.")
	flag.BoolVar(&print0, "0", false, "Separate search results by the null character (instead of newlines). Useful for piping results to 'xargs'.")
	flag.Usage = usage
	flag.Parse()
}
