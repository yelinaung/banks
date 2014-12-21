# Banks

Just scrap-and-show API for currency exchange rates of Banks in Myanmar.

Supported banks
- [KBZ](kbzbank.com)
- [CB](www.cbbank.com.mm)


## How to run

```bash
$ go get github.com/yelinaung/banks
$ go run banks.go
```

You have to put the bank name as parameters in path

e.g For KBZ, put the bank name after the base url. Same for CB.

`localhost:3001/kbz`

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

## TOOD 

- Deploy to Heroku
- To add AGD 
- To scrap periodically ?
- To cache with the dates ?


## Lincese
MIT

