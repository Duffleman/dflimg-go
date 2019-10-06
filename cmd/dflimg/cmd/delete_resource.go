package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"dflimg/cmd/dflimg/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteResourceCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "Delete a resource",
	Long:    "Delete a resource",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		urlStr := args[0]

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		err := deleteResource(rootURL, authToken, urlStr)
		if err != nil {
			return err
		}

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s\n", duration)

		return nil
	},
}

func deleteResource(rootURL, authToken, urlStr string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormField("input")
	if err != nil {
		return err
	}

	io.Copy(part, strings.NewReader(urlStr))

	writer.Close()

	c := http.New(rootURL, authToken)

	return c.Request("POST", "delete_resource", body, writer.FormDataContentType(), nil)
}