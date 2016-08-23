# Banks

Just scrap-and-show API for currency exchange rates of Banks in Myanmar.

Supported banks

- [KBZ Bank](http://www.kbzbank.com)
- [CB Bank](http://www.cbbank.com.mm)
- [AYA Bank](http://ayabank.com)
- [AGD Bank](http://www.agdbank.com)
- [MAB](http://www.mabbank.com)
- [UAB](http://www.unitedamarabank.com)


## How to run

```bash
$ go get github.com/yelinaung/banks
$ cd $GOPATH/src/github.com/yelinaung/banks
$ go get
$ export PORT="8080" && go build && ./banks
```

## Usage 

Getting latest rates by bank 
e.g put the bank name after the base url. For example, 

```http
GET 

/b/[bank name]

```

Getting latest rates

```http
GET 

/latest
```

## Contributing

  1. Fork it
  2. Create your feature branch (`git checkout -b my-new-feature`)
  3. Commit your changes (`git commit -am 'Added some feature'`)
  4. Push to the branch (`git push origin my-new-feature`)


## Lincese
MIT

