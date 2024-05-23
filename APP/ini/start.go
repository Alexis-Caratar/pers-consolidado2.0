package main


func main() {

	router := vestigo.NewRouter()
	config.Load(router) // se esta refiriendo al paquete config / load (el archivo en realidad se llama router)

	router.Get("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("../../PUBLIC/"))).ServeHTTP)

	conf, err := dal.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server Start SISNOVA-PERS", conf.Port)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(conf.Port), router))
}
