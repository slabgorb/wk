#Installation

I used Go 1.8 for this, I have not tested it with earlier versions of Go, but I suspect it would be compatible with many earlier versions, as I am not using any recent stdlibs like `context`.

Probably the easiest way to install the program into the proper place, if you have Go already set up is to do a

`go get github.com/slabgorb/wk`


#Running

`cd $GOPATH/src/github.com/slabgorb/wk`

`go install && wk -file ./input.wk`

#Testing

`cd $GOPATH/src/github.com/slabgorb/wk`

`go test -v ./...`


#File format for input file
A workflow file is expected to be in the following format:

    command:name # comments are denoted by '#' and anything on a line after a # is ignored.
      argumentname:value
      argumentname:value
      argumentname:value


    parallelcommand:name
     argumentname:value
     argumentname:value

    parallelcommand:name
     argumentname:value
     argumentname:value
