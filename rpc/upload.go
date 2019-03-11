package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maxUploadSize = 100 * 1024 // 100 MB
const uploadPath = "./tmp"

// Upload is an RPC handler for uploading a file
func (r *RPC) Upload(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	file, _, err := req.FormFile("file")
	if err != nil {
		r.logger.WithError(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	res, err := r.app.Upload(ctx, buf)
	if err != nil {
		fmt.Println(err)
		r.logger.WithError(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(res)

	return
}
