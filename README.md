## Overview

This is simple SDK implementation for the [Pocket](https://getpocket.com/developer/?src=footer_v2)
Before using it, you must create a [Pocket Application](https://getpocket.com/developer/apps/new)

This SDK implements all [Developer API features](https://getpocket.com/developer/docs/overview):

- [Authentication](https://getpocket.com/developer/docs/authentication)
- [Modify](https://getpocket.com/developer/docs/v3/modify)
- [Retrieve](https://getpocket.com/developer/docs/v3/retrieve)
- [Add](https://getpocket.com/developer/docs/v3/add)

## Content

- [Installation](#installation)
- [Authentication](#authentication)
  - [Create a pocket object](#create-a-pocket-object)
  - [Generate a request token](#generate-a-request-token)
  - [User authorization](#user-authorization)
  - [Generate an access token](#generate-an-access-token)

## Installation

To install the package run:

```bash
go get -u github.com/VladimirStepanov/pocket-golang-sdk
```

## Authentication

Authentication performs in 4 steps:
- create a Pocket object
- get a request token
- user authorization
- convert a request token into a Pocket access token

### Create a pocket object

For creating a Pocket object, you need a consumer key, you've gotten after app registration

```go
p := pocket.New("consumer-key")
```

### Generate a request token

Method 1 
```go
err := p.AuthApp(context.Background(), "redirect-url")

if err != nil {
    log.Fatal(err)
}

fmt.Println("Request token", p.GetRequestToken())
```

Method 2
```go
res, err := p.GenerateRequestToken(context.Background(), "redirect-url")

if err != nil {
    log.Fatal(err)
}

p.SetRequestToken(res.Code)
fmt.Println("Request token", p.GetRequestToken())
```

### User authorization

Once you have a request token, you need to redirect the user to Pocket to authorize your
application's request token. For getting authorize link use this method:

```go
link := p.MakeAuthUrl("redirect-url")

// result example: https://getpocket.com/auth/authorize?redirect_uri=https%3A%2F%2Fgoogle.com&request_token=ffffcc4e-ffff-ffff-ffff-f7f68f 
```

### Generate an access token

After successfully user authorization, you can get an access token

Method 1
```go
err := p.AuthUser(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Println("Access token", p.GetAccessToken())
```

Method 2
```go
at, err := p.GenerateAccessToken(context.Background())
if err != nil {
    log.Fatal(err)
}
p.SetAccessToken(at.AccessToken)
fmt.Println("Access token", p.GetAccessToken())
```