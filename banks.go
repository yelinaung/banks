package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	//	"github.com/gin-gonic/gin"
)

var tableName = "currency"
var s *r.Session

func init() {
	var err error

	s, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbName,
		MaxOpen:  10,
	})

	if err != nil {
		fmt.Errorf("failed to connect to database: %v", err)
	}

	_, err1 := r.DB(dbName).TableCreate(tableName).RunWrite(s)

	if err1 == nil {
		fmt.Printf("Error creating table: %s", err1)
	} else {
		r.DB(dbName).TableCreate(tableName).RunWrite(s)
	}
}

func main() {
	fmt.Println("Starting..")

	//c := gron.New()
	//c.AddFunc(gron.Every(2 * time.Minute), func() {
	//	fmt.Println("runs every two minute.")

	//})
	//c.Start()

	// r := gin.New()
	//
	// r.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK,
	// 		"Nothing to see here.Check https://github.com/yelinaung/banks")
	// })

	//var err error

	r.GET("/:bank", func(c *gin.Context) {
		bankName := c.Params.ByName("bank")
		switch bankName {
		case "kbz":
			bank, err = process(scrapKBZ())
			bank.Name = "KBZ"
		case "mab":
			bank, err = process(scrapMAB())
			bank.Name = "MAB"
		case "uab":
			bank, err = process(scrapUAB())
			bank.Name = "UAB"
		case "cbb":
			bank, err = process(scrapCBB())
			bank.Name = "CBB"
		case "agd":
			bank, err = process(scrapAGD())
			bank.Name = "AGD"
		case "aya":
			bank, err = process(scrapAYA())
			bank.Name = "AYA"
		default:
			// TODO	what to reply for default
		}
	
		if err == nil {
			c.JSON(http.StatusOK, bank)
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"message": "Something went wrong!",
				})
		}
	
	})
	r.Run(":" + os.Getenv("PORT"))

}

func filterByBankName(name string) {
	query := r.Table(tableName).Filter(r.Row.Field("name").Eq(name))
	row, err := query.Run(s)
	if err != nil {
		fmt.Print(err)
		return
	}

	var currencies = []Currency{}
	err2 := row.All(&currencies)

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	printObj(currencies)
}
