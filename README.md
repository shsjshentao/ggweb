# ggweb
##简易的web框架
···
func main() {

	route := ggweb.NewRoute()
	
	route.AddBefore(Exec)
	
	g1 := route.AddGroup("/a")
	
	g1.Handle("/b", Exec)
	
	g1.Handle("c", Exec)
	
	route.AddRoute("/aa", Exec)
	
	http.ListenAndServe(":80", route)
	
}

func Exec(c *ggweb.Context) {

}···
