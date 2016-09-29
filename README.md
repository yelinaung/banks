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

You will need [forego](https://github.com/ddollar/forego) to run two processes.
And also install [rethinkdb](https://rethinkdb.com/docs/install/) for storing the json.

```bash
$ go get -u github.com/ddollar/forego
$ go get -u github.com/yelinaung/banks
$ cd $GOPATH/src/github.com/yelinaung/banks
$ go get
$ rethinkdb
$ ./run # or ./run-prod if you want to run in production
```

## Usage 

#### Getting latest rates by bank 

e.g put the bank name after the base url. For example, 

```http
GET 

/api/v1/b/[bank name]

```

Example available at : [http://c.yelinaung.com/api/v1/b/kbz](http://c.yelinaung.com/api/v1/b/kbz)

#### Getting latest rates

During debug mode, scraper runs **every 20 seconds** and during prod mode, the scraper runs **every 2 hours**.

```http
GET 

/api/v1/latest
```

Example available at : [http://c.yelinaung.com/api/v1/latest](http://c.yelinaung.com/api/v1/latest)

## Contributing

  1. Fork it
  2. Create your feature branch (`git checkout -b my-new-feature`)
  3. Commit your changes (`git commit -am 'Added some feature'`)
  4. Push to the branch (`git push origin my-new-feature`)


## License

MIT

