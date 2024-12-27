package main

func main() {
	app := App{}
	app.Initialize(DbUsername, DbPassword, DbName)
	app.Run("localhost:10000")
}