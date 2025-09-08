package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"vcon/internal/db"
	"vcon/internal/globalStore"
	"vcon/internal/services"

	"github.com/gin-gonic/gin"
)

type pureDocument struct {
	Title string   `json:"title"`
	Array []string `json:"array"`
}

func main() {

	// connect with DB
	dataBase, err := db.DBConnect()

	if err != nil {
		fmt.Errorf(" DB connection failed ")
	}

	// make gloal store and document Service
	globalStore := globalStore.InitializeStore()
	docService := services.NewDocumentService(dataBase, globalStore)

	// start server and attck controllers to API routes
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/adddocument", func(c *gin.Context) {

		var requestDocument pureDocument

		if err := c.ShouldBindJSON(&requestDocument); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// else we have the document from the frontn d

		docSaved, err := docService.AddDocument(context.Background(), requestDocument.Title, requestDocument.Array)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// else the document is saved and well hydrated
		fmt.Printf(" Saved document : ", docSaved)
		c.JSON(200, gin.H{
			"success": "true",
		})
		return
	})

	router.GET("/getalldoc", func(c *gin.Context) {
		var result []string

		result, err := docService.LoadTitleOfAllDocuments(context.Background())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// succesfull

		fmt.Printf(" Saved document : ", result)
		c.JSON(200, gin.H{
			"docArray": result,
		})
		return

	})

	router.GET("/getbytitle", func(c *gin.Context) {
		// "http://localhost:8080/getbytitle?title=My First Document"
		title := c.Query("title")

		doc, err := docService.GetDocumentByTitle(context.Background(), title)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf(" returning document : ", doc)
		c.JSON(200, gin.H{
			"data": doc,
		})

		return

		// retrn dccument

	})

	type ReqBody struct {
		DocTitle   string
		ParentNode int
		StringArr  []string
	}

	router.POST("/addversiontodocument", func(c *gin.Context) {
		var ReqBody ReqBody

		if err := c.ShouldBindJSON(&ReqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := docService.AddVersionToDocument(context.Background(), ReqBody.DocTitle, ReqBody.ParentNode, ReqBody.StringArr)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// else the document is saved and well hydrated
		fmt.Printf(" Saved version : ")
		c.JSON(200, gin.H{
			"success": "true",
		})
		return

	})

	router.GET("/getversion", func(c *gin.Context) {

		// /getversion?title=My First Document&version=2
		
		title := c.Query("title")
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title query parameter is required"})
			return
		}
		versionStr := c.Query("version")
		if versionStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "version query parameter is required"})
			return
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "version must be a valid number"})
			return
		}
		contentHashes, err := docService.GetVersionFromDocument(c.Request.Context(), version, title)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }

        stringContent, err := docService.ConvertHashesToStrings(contentHashes)

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render document content"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "title":   title,
            "version": version,
            "content": stringContent,
        })


	})

	router.Run()
}
