package handler

import (
	"blendcube/config"
	"blendcube/cube"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type (
	errorResponse struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	}
)

const (
	glTF = "gltf"
	glb  = "glb"
)

var (
	faces = [18]string{
		"U", "D", "F", "B", "L", "R",
		"U'", "D'", "F'", "B'", "L'", "R'",
		"U2", "D2", "F2", "B2", "L2", "R2",
	}
)

func (er errorResponse) Write(w http.ResponseWriter) {
	body, err := json.Marshal(er)
	if err != nil {
		return
	}
	w.WriteHeader(er.StatusCode)
	_, _ = w.Write(body)
}

// New ルーターを生成する
func New() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(
		config.Shared.Server.Timeout * time.Second,
	))

	r.Get("/status", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.With(middleware.URLFormat).Get("/cube", generateCube)

	return r
}

func generateCube(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = r.Context()
		query     = r.URL.Query()
		algorithm []string
		isBinary  bool
		isUnlit   bool
	)

	format, ok := ctx.Value(middleware.URLFormatCtxKey).(string)
	if ok {
		if format != glTF && format != glb {
			errorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "the following extension are supported: .gltf, .glb",
			}.Write(w)
			return
		}
		isBinary = format == glb
	}

	if alg := query.Get("alg"); alg != "" {
		algs := strings.Split(alg, " ")
		for _, a := range algs {
			hit := false
			for _, f := range faces {
				if a == f {
					hit = true
				}
			}
			if !hit {
				errorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    `alg must only use "U D F B L R U' D' F' B' L' R' U2 D2 F2 B2 L2 R2"`,
				}.Write(w)
				return
			}
		}
		algorithm = algs
	}

	if isUnlitStr := query.Get("is_unlit"); isUnlitStr != "" {
		if isUnlitStr != "true" && isUnlitStr != "false" {
			errorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "is_unlit must be true or false",
			}.Write(w)
			return
		}
		isUnlit = isUnlitStr == "true"
	}

	body, err := cube.Generate(algorithm, isBinary, isUnlit)
	if err != nil {
		panic(err)
	}

	if isBinary {
		w.Header().Set("Content-Type", "model/gltf-binary")
	} else {
		w.Header().Set("Content-Type", "model/gltf+json")
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
