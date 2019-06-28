package main

type IDType int32

func generateID() func() IDType {
	var id IDType = 1
	return func() IDType {
		id++
		return id
	}
}

var GenerateID = generateID()

type Msg struct {
	UserID IDType
	Content string
}
