// @title           ContamianDOS API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

package api

import (
	command "be/cmd/server/handler/http/command"
	querie "be/cmd/server/handler/http/querie"
	"be/cmd/server/handler/ws"
	_ "be/docs"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ApiServer struct {
	Address    string
	TLSAddress string
	CertFile   string
	KeyFile    string
	server     *http.Server
	tlsServer  *http.Server
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		Address: addr,
	}
}

func NewApiServerWithTLS(addr, tlsAddr, certFile, keyFile string) *ApiServer {
	return &ApiServer{
		Address:    addr,
		TLSAddress: tlsAddr,
		CertFile:   certFile,
		KeyFile:    keyFile,
	}
}

func (apiserver *ApiServer) Run() error {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "UPDATE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token", "X-Max", "password", "player", "owner","name", "Origin", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           24 * time.Hour,
	}))
	//swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//demo
	router.GET("/hello", querie.GetHello)

	/*
		Seccion Players
	*/
	//Gets the game by ID
	router.GET("/api/games/:gameId", querie.GetGame)
	//Joins game
	router.PUT("/api/games/:gameId", command.JoinGame)
	//Starts game
	router.HEAD("/api/games/:gameId/start", command.GameStart)
	//Gets round of a specific game
	router.GET("/api/games/:gameId/rounds", querie.GetRounds)
	//
	router.GET("/api/games/:gameId/rounds/:roundId", querie.ShowRound)
	//Proposes a group in epoch
	router.PATCH("/api/games/:gameId/rounds/:roundId", command.ProposeGroup)
	//Vote for a group in epoch
	router.POST("/api/games/:gameId/rounds/:roundId", command.VoteGroup)
	//Member of proposed group submit action
	router.PUT("/api/games/:gameId/rounds/:roundId", command.SubmitAction)

	/*
		Seccion Public
	*/
	// Creates game
	router.POST("/api/games", command.CreateGame)
	router.POST("/api/games/", command.CreateGame)
	//Search for a game by name and status
	router.GET("/api/games", querie.GameSearch)
	router.GET("/api/games/", querie.GameSearch)

	//ws
	router.GET("/ws", ws.HandleWebSocket)

	srv := &http.Server{
		Addr:    apiserver.Address,
		Handler: router,
	}
	apiserver.server = srv

	if apiserver.TLSAddress != "" && apiserver.CertFile != "" && apiserver.KeyFile != "" {
		tlsSrv := &http.Server{
			Addr:    apiserver.TLSAddress,
			Handler: router,
		}
		apiserver.tlsServer = tlsSrv

		go func() {
			log.Printf("HTTPS server listening on %s", apiserver.TLSAddress)
			if err := tlsSrv.ListenAndServeTLS(apiserver.CertFile, apiserver.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTPS server error: %v", err)
			}
		}()
	}

	log.Printf("HTTP server listening on %s", apiserver.Address)
	return srv.ListenAndServe()
}

func (apiserver *ApiServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down API server...")
	var err error
	if apiserver.tlsServer != nil {
		if e := apiserver.tlsServer.Shutdown(ctx); e != nil {
			log.Println("HTTPS shutdown error:", e)
			err = e
		}
	}
	if apiserver.server != nil {
		if e := apiserver.server.Shutdown(ctx); e != nil {
			log.Println("HTTP shutdown error:", e)
			err = e
		}
	}
	return err
}

// Funcion que devuelve respuesta en la ruta /record
/*
func (apiserver *ApiServer) HandlerRecord(writer http.ResponseWriter, request *http.Request) {
	writer = utils.SetCORSHeaders(writer, request)
	list,err := utils.GetParamsRequest(request)
	if err != nil {
		response := types.Response{
			Status: http.StatusBadRequest,
			Msg:    err.Error(),
			Data:   nil,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, "Error serializing JSON", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(jsonResponse))
		return
	}
	msg := utils.ResponseMessage(len(data))

	response := types.Response{
		Status: http.StatusOK,
		Msg:    msg,
		Data:   data,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Error serializing JSON", http.StatusInternalServerError)
		return
	}
	writer.Write(jsonResponse)
}
*/
