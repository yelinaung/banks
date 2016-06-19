package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting..")

	Run()

	//c := gron.New()
	//c.AddFunc(gron.Every(2 * time.Minute), func() {
	//	fmt.Println("runs every two minute.")

	//})
	//c.Start()

	//r := gin.New()
	//
	//r.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK,
	//		"Nothing to see here.Check https://github.com/yelinaung/banks")
	//})

	//var err error

	//r.GET("/:bank", func(c *gin.Context) {
	//	bankName := c.Params.ByName("bank")
	//	switch bankName {
	//case "kbz":
	//	//bank, err = process(scrapKBZ())
	//	//bank.Name = "KBZ"
	//case "mab":
	//	bank, err = process(scrapMAB())
	//	bank.Name = "MAB"
	//case "uab":
	//	bank, err = process(scrapUAB())
	//	bank.Name = "UAB"
	//case "cbb":
	//	bank, err = process(scrapCBB())
	//	bank.Name = "CBB"
	//case "agd":
	//	bank, err = process(scrapAGD())
	//	bank.Name = "AGD"
	//case "aya":
	//	bank, err = process(scrapAYA())
	//	bank.Name = "AYA"
	//default:
	// TODO	what to reply for default
	//}

	//if err == nil {
	//	c.JSON(http.StatusOK, bank)
	//} else {
	//c.JSON(http.StatusInternalServerError,
	//	gin.H{
	//		"message": "Something went wrong!",
	//	})
	//}

	//})
	//r.Run(":" + os.Getenv("PORT"))
}
