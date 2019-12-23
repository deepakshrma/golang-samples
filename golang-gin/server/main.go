package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

type person struct {
	fname string
	lname string
}
type employee struct {
	person
	developer bool
}
type Member struct {
	Age int `something:"1"`
}

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	// Disable log's color
	// gin.DisableConsoleColor()
	// Logging to a file.
	employee := employee{
		person: person{
			"Deepak",
			"Sharma",
		},
		// developer: true,
	}
	fmt.Printf("New Employee: %v\n", employee)
	f, _ := os.Create("gin.log")

	member := Member{34}
	t := reflect.TypeOf(member)
	field := t.Field(0)
	something := field.Tag.Get("something")
	//field, _ := t.FieldByName("Age") //alternative
	fmt.Println(something)
	// gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(Logger())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Static("/assets", "./assets")
	// router.StaticFS("/more_static", http.Dir("my_file_system"))
	// router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	dir, _ := os.Getwd()
	router.LoadHTMLGlob(path.Join(dir, "templates/**/*"))

	router.SetHTMLTemplate(html)

	// Custom template loader
	// t, err := loadTemplate()
	// if err != nil {
	// 	panic(err)
	// }
	// r.SetHTMLTemplate(t)
	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	// r.Use(gin.Logger())

	// // Recovery middleware recovers from any panics and writes a 500 if there was one.
	// r.Use(gin.Recovery())

	// // Per route middleware, you can add as many as you desire.
	// r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// // Authorization group
	// // authorized := r.Group("/", AuthRequired())
	// // exactly the same as:
	// authorized := r.Group("/")
	// // per group middleware! in this case we use the custom created
	// // AuthRequired() middleware just in the "authorized" group.
	// authorized.Use(AuthRequired())
	// {
	// 	authorized.POST("/login", loginEndpoint)
	// 	authorized.POST("/submit", submitEndpoint)
	// 	authorized.POST("/read", readEndpoint)

	// 	// nested group
	// 	testing := authorized.Group("testing")
	// 	testing.GET("/analytics", analyticsEndpoint)
	// }
	// gin.H is a shortcut for map[string]interface{}
	/*
	   router.GET("/someJSON", func(c *gin.Context) {
	   	c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	   })

	   router.GET("/moreJSON", func(c *gin.Context) {
	   	// You also can use a struct
	   	var msg struct {
	   		Name    string `json:"user"`
	   		Message string
	   		Number  int
	   	}
	   	msg.Name = "Lena"
	   	msg.Message = "hey"
	   	msg.Number = 123
	   	// Note that msg.Name becomes "user" in the JSON
	   	// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
	   	c.JSON(http.StatusOK, msg)
	   })

	   router.GET("/someXML", func(c *gin.Context) {
	   	c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	   })

	   router.GET("/someYAML", func(c *gin.Context) {
	   	c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	   })

	   router.GET("/someProtoBuf", func(c *gin.Context) {
	   	reps := []int64{int64(1), int64(2)}
	   	label := "test"
	   	// The specific definition of protobuf is written in the testdata/protoexample file.
	   	data := &protoexample.Test{
	   		Label: &label,
	   		Reps:  reps,
	   	}
	   	// Note that data becomes binary data in the response
	   	// Will output protoexample.Test protobuf serialized data
	   	c.ProtoBuf(http.StatusOK, data)
	   })
	*/
	router.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
	})

	router.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	router.GET("/secureJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// Will output  :   while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})
	// Example for binding JSON ({"user": "manu", "password": "123"})

	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>user</user>
	//		<password>123</password>
	//	</root>)
	router.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	router.GET("/JSONP?callback=x", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}

		//callback is x
		// Will output  :   x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	router.POST("/post", func(c *gin.Context) {

		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")
		test := c.PostFormArray("test")

		fmt.Printf("\n###############\nids: %v; names: %v, test: %v\n###############\n", ids, names, test)
	})

	router.GET("/http2", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// use pusher.Push() to do server push
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
	})
	router.GET("/uri/:name/:id", bindUri)
	router.GET("/person", getPerson)
	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
			"now":   time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
		})
	})
	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	router.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note that you are using the copied context "cCp", IMPORTANT
			log.Println("Done! in path " + cCp.Request.URL.Path)

		}()
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})

	router.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// since we are NOT using a goroutine, we do not have to copy the context
		log.Println("Done! in path " + c.Request.URL.Path)
	})
	// 	$ curl "localhost:8085/bookable?check_in=2018-04-16&check_out=2018-04-17"
	// {"message":"Booking dates are valid!"}

	// $ curl "localhost:8085/bookable?check_in=2018-03-10&check_out=2018-03-09"
	// {"error":"Key: 'Booking.CheckOut' Error:Field validation for 'CheckOut' failed on the 'gtfield' tag"}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	router.GET("/bookable", getBookable)

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", sample)
		v1.POST("/submit", sample)
		v1.GET("/read", sample)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", sample)
		v2.POST("/submit", sample)
		v2.GET("/read", sample)
	}
	// router.Run(":8080")
	// http.ListenAndServe(":8080", router)

	// Advance
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServeTLS("./testdata/server.pem", "./testdata/server.key")
}
func sample(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"success": "Success"})
}

//$ curl -X GET "localhost:8085/testing?name=appleboy&address=xyz&birthday=1992-03-15"
type BindURI struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

// $ curl -v localhost:8088/uri/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3
// $ curl -v localhost:8088/uri/thinkerou/not-uuid

func bindUri(c *gin.Context) {
	var bindUri BindURI
	if err := c.ShouldBindUri(&bindUri); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"name": bindUri.Name, "uuid": bindUri.ID})
}

// loadTemplate loads templates embedded by go-assets-builder
// func loadTemplate() (*template.Template, error) {
// 	t := template.New("")
// 	for name, file := range Assets.Files {
// 		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
// 			continue
// 		}
// 		h, err := ioutil.ReadAll(file)
// 		if err != nil {
// 			return nil, err
// 		}
// 		t, err = t.New(name).Parse(string(h))
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return t, nil
// }

func getPerson(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
	}

	c.String(200, "Success")
}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

// Booking contains binded and validated data.
type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
