module server

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/mebusy/goweb v0.0.0-20210603083155-272a5f2d3e92
	github.com/mebusy/gowebfront v0.0.0
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
)

replace github.com/mebusy/gowebfront v0.0.0 => ../
