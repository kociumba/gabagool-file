# Gabagool™️ file format

I lost a full nights of sleep for this and so here it is.

`.gabagool` is a universal file format, right now it can store:
- Text - (like any other text file)
- Images - (in a format that is compatible with the go canvas)
- Binary data - (any data in raw bytes)

## Usage

simply import the gabagool package from this repo like this in go
```go
import "github.com/kociumba/gabagool-file/gabagool"
```

the main package will be a `.gabagool` interpreter for use with desktop environments.
 
