package main

import (
  "gopkg.in/kataras/iris.v6"
  "gopkg.in/kataras/iris.v6/adaptors/httprouter"
  "gopkg.in/kataras/iris.v6/adaptors/view"
  "fmt"
  "github.com/culqi/culqi-go"
  "github.com/culqi/culqi-go/charge"
  "encoding/json"

)

func main() {
  app := iris.New()
  app.Adapt(iris.DevLogger())
  app.Adapt(httprouter.New())

  tmpl := view.HTML("./views", ".html")
  tmpl.Layout("index.html")


  app.Adapt(tmpl)



  app.Get("/", func(ctx *iris.Context) {
    ctx.Render("index.html", iris.Map{"gzip": true})
  })

  app.Post("/cargo", func(ctx *iris.Context) {
      fmt.Print("hola cargo")

    token := ctx.PostValue("token")
		fmt.Printf("\nResponse Status Code: %v", token)


    // 1. Configuración
    config := &culqi.Config{
      MerchantCode:   "pk_test_Rp2uV5dXI3quFq2X",  // Código de tu Comercio
      ApiKey:   "sk_test_8GC9UJfifciOurwW", // API Key
      //ApiVersion: "v2",   // No es un parametro necesario, por defecto es la v2
    }

    // 2. Crea un nuevo cliente
    client := culqi.New(config)
    // 3. Parametros de creación de cargo
    params := &charge.ChargeParams{
      TokenId: token ,
      Email: "liz@mail.com",
      CurrencyCode: "PEN",
      Amount: 100,
    }

    // 4. Crear Cargo
    resp, err := charge.Create(params, client)

    if err != nil {
        fmt.Printf(err.Error())
    }
    fmt.Printf("\nResponse Body: %v", resp)
    fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())

    // 5. convertir response en variable de go

    type Outcome struct {
      UserMessage string `json:"user_message"`
      Code string `json:"code"`
      MerchantMessage string  `json:"merchant_message"`
    }

		type TokenResponse struct {

			Object string `json:"object"`
			Id string `json:"id"`
			Email string `json:"email"`
      Outcome Outcome `json:"outcome"`
		}



		var jsontype TokenResponse

		json.Unmarshal([]byte(resp.Body()), &jsontype)


		fmt.Printf("\nResponse Body Object: %v", jsontype.Object)
    fmt.Printf("\nResponse Body Object: %v", jsontype)

		//6. response json

		ctx.JSON(200,jsontype)

  })


  app.StaticWeb("/static", "./assets")
  app.Listen(":8082")
}
