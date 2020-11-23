# dropf

*Drop your files*

<img src="https://raw.githubusercontent.com/jubalh/dropf/master/static/images/dropf.png" height="100">

dropf can be used to quickly upload files to your server and share them with others.  
It is meant as a private replacement for things like MegaUpload, imgur or Dropbox.

The goal is to have this easily managable over ssh on the server side.
Configuration should be in a simple file.

It is for people who like to share files by uploading them (with scp, for example) but sometimes not having a terminal at hand, or having friends who don't know how to use ssh.

For this reason files will be directly stored in a users directory for association.

## Installation

```
go get github.com/jubalh/dropf
```

## Usage

Use the example `config.json` and place it next to the binary. Edit it to add users.
Place the `static` and `templates` directories next to the binary.
Then run it.

