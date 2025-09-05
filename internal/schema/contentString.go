package schema

type ContentString struct {
	ID      string `bson:"_id"` // identifier ( SHA 256 Hash ) will be mongoDB id as th would provide fater
	Content string `bson:"content"`
}

/*

EXAMPLE


{ "_id": "45",  "content": "import React from 'react';" }
{ "_id": "101", "content": "function App() {" }
{ "_id": "800", "content": "  return <h1>Hello, World!</h1>;" }
{ "_id": "2",   "content": "}" }

*/
