# imagext Package

[![GoDoc](https://godoc.org/github.com/vbsw/imagext?status.svg)](https://godoc.org/github.com/vbsw/imagext) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/imagext)](https://goreportcard.com/report/github.com/vbsw/imagext) [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## About
This package provides functions for images. It's written for fast prototyping. Results may depend on format.

This Package is published on <https://github.com/vbsw/imagext>.

## Copyright
Copyright 2021, Vitali Baumtrok (vbsw@mailbox.org).

imagext Package is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

imagext Package is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Example

	package main
	import (
		"fmt"
		"github.com/vbsw/imagext"
	)

	func main() {
		imgColor := imagext.LoadImage("/home/alice/pictures/example.png")
		imgGray := imagext.NewGray(imgColor)
		err := imagext.SaveImage("/home/alice/pictures/example-gray.png", imgGray)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

## References
- https://golang.org/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
- https://en.wikipedia.org/wiki/Grayscale
- https://www.rapidtables.com/convert/color/cmyk-to-rgb.html
- https://www.sciencedirect.com/topics/computer-science/median-filter
