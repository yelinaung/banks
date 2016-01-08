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
$ export PORT="8080" && go run banks.go
```

You have to put the bank name as parameters in path

e.g For KBZ, put the bank name after the base url. Same for other banks.

- KBZ - `localhost:8080/kbz`
- CB  - `localhost:8080/cb`
- AYA - `localhost:8080/aya`
- AGD - `localhost:8080/agd`
- MAB - `localhost:8080/mab`
- UAB - `localhost:8080/uab`

You can try deploying to [Heroku](www.heroku.com) too

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## Sample response

```json
{
  "name":"KBZ",
  "base":"MMK",
  "time":"2014-12-21 14:51:06.97045683 +0630 MMT",
  "rates":[
    {
      "USD":{
        "buy":"1025",
        "sell":"1034"
      }
    },
    {
      "EUR":{
        "buy":"1249",
        "sell":"1268"
      }
    },
    {
      "SGD":{
        "buy":"774",
        "sell":"786"
      }
    }
  ]
}
```

## TODO

- ~~Deploy to Heroku~~
- ~~To add AGD & AYA~~
- To scrap periodically ?
- To cache with the dates ?

## Contributing

  1. Fork it
  2. Create your feature branch (`git checkout -b my-new-feature`)
  3. Commit your changes (`git commit -am 'Added some feature'`)
  4. Push to the branch (`git push origin my-new-feature`)


## Lincese
MIT

