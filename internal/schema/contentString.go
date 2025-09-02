package schema


type StringInteger struct{
	ID int 	`bson:"_id"` // identifier will be mongoDB id as th would provide fater 
	Content string	`bson:"content"`;
}       
/*

EXAMPLE 


{ "_id": 45,  "content": "import React from 'react';" }
{ "_id": 101, "content": "function App() {" }
{ "_id": 800, "content": "  return <h1>Hello, World!</h1>;" }
{ "_id": 2,   "content": "}" }

*/
