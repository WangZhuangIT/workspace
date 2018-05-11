package apis

import (
	"fmt"
	. "gin-demo/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexApi(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"msg": "it works",
	})
}

func AddPersonApi(c *gin.Context) {
	firstname := c.Request.FormValue("1")
	lastname := c.Request.FormValue("2")

	person := Person{LastName: lastname, FirstName: firstname}
	id, err := person.AddPerson()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("insert into last id ", id)

	c.JSON(http.StatusOK, gin.H{
		"msg": "123msg",
	})
}

func QueryPersonApi(c *gin.Context) {
	person := Person{}
	persons := make([]Person, 0)
	persons, err := person.GetPerson()
	if err != nil {
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"persons": persons,
	})
}

func QueryPersonApiById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln(err)
	}
	person := Person{Id: id}
	person, err = person.GetPersonOne()
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusOK, gin.H{
			"person": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"person": person,
	})
}

func UpdatePersonApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatal(err)
	}
	person := Person{Id: id}
	err = c.Bind(&person)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err := person.UpdatePerson()
	if err != nil {
		log.Fatalln(err)
	}

	msg := fmt.Sprintf("update person %d successful %d", id, ra)

	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func DelPersonApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	person := Person{Id: id}
	err, id, rowid := person.DelPerson()

	if err != nil {
		log.Fatalln(err)
	}

	msg := fmt.Sprintf("delete person %d successful %d", id, rowid)

	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}
