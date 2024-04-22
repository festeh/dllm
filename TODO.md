# Setting API key


You can use Go's linker to set the value of a string variable at compile time. Here's how you can do it:

1. Define a string variable in your Go code that will hold the API key. This variable should be at the package level scope. For example:

```go
package main

import "fmt"

var apiKey string

func main() {
    fmt.Println("API Key: ", apiKey)
}
```

2. When you compile your Go code, use the `-ldflags` option to set the value of the `apiKey` variable. The `-X` flag is used to set the value of a string variable. The syntax is `packageName.variableName=value`. For example:

```bash
go build -ldflags "-X main.apiKey=your_api_key_here"
```

This will set the value of the `apiKey` variable to `your_api_key_here` at compile time. When you run your program, it will print out `API Key: your_api_key_here`.

Remember to replace `your_api_key_here` with your actual API key.

Please note that storing sensitive information like API keys in your binary can be risky. If someone gets access to your binary, they can potentially extract the API key from it. It's generally better to use environment variables or some sort of secure configuration to store sensitive information.

