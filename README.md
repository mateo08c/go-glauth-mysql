# go-glauth-mysql
Interacts with the [Glauth](https://github.com/glauth/glauth) MySQL database in Go !

- [go-glauth-mysql](#go-glauth-mysql)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Environment Variables For Testing](#environment-variables-for-testing)

## Installation
```bash
go get github.com/mateo08c/go-glauth-mysql
```

## Usage

```go
package main

import (
  "github.com/mateo08c/go-glauth-mysql/glauth"
  "log"
)

func main() {
  client, err := glauth.New(&glauth.Context{
    Username: "glauth",
    Password: "password",
    Hostname: "example.com",
    Port:     "3306",
    Database: "glauth",
  })
  if err != nil {
    log.Fatal(err)
  }

  users, err := client.GetUsers()
  if err != nil {
    log.Fatal(err)
  }

  for _, user := range users {
    log.Println(user)
  }
}

```

### Environment Variables For Testing
| Variable | Description |
|----------|-------------|
| DB_USERNAME | MySQL username |
| DB_PASSWORD | MySQL password |
| DB_HOSTNAME | MySQL hostname | 
| DB_PORT | MySQL port |
| DB_NAME | MySQL database name |
