package main

func main() {
	engine, err := CreateEngine()
	if err != nil {
		panic(err)
	}
	defer engine.CleanUp()

	engine.Run()
}
