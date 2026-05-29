package command

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"be/internal/logic"
	"be/util/types"

	"github.com/gin-gonic/gin"
)

func setupPublicRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestCreateGame_PublicSuccess(t *testing.T) {
	logic.NewGamePool()
	r := setupPublicRouter()
	r.POST("/api/games", CreateGame)

	body := types.CreateRequest{Name: "PublicGame", Owner: "Owner1", Password: "pass123"}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCreateGame_PublicMissingName(t *testing.T) {
	logic.NewGamePool()
	r := setupPublicRouter()
	r.POST("/api/games", CreateGame)

	body := types.CreateRequest{Name: "", Owner: "Owner1", Password: ""}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing name, got %d", w.Code)
	}
}

func TestCreateGame_PublicMissingOwner(t *testing.T) {
	logic.NewGamePool()
	r := setupPublicRouter()
	r.POST("/api/games", CreateGame)

	body := types.CreateRequest{Name: "Game1", Owner: "", Password: ""}
	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/games", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing owner, got %d", w.Code)
	}
}
