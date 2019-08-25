package cmd

import (
	"bytes"
	"dflimg/dflerr"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ShortenURLCmd = &cobra.Command{
	Use:     "shorten",
	Aliases: []string{"s"},
	Short:   "Shorten a URL",
	Long:    "Shorten a URL",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		urlStr := args[0]
		shortcuts := cmd.Flag("shortcuts")
		nsfw := cmd.Flag("nsfw")

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		body, err := shortenURL(rootURL, authToken, urlStr, shortcuts, nsfw)
		if err != nil {
			return err
		}

		r, err := parseResponse(body)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(r.URL)
		if err != nil {
			fmt.Println("Could not copy to clipboard. Please copy the URL manually")
		}

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s: %s\n", duration, r.URL)

		return nil
	},
}

func shortenURL(rootURL, authToken, urlStr string, shortcuts, nsfw *pflag.Flag) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if shortcuts != nil {
		shortcutsStr := shortcuts.Value.String()
		part, err := writer.CreateFormField("shortcuts")
		if err != nil {
			return nil, err
		}

		io.Copy(part, strings.NewReader(shortcutsStr))
	}

	if nsfw != nil {
		nsfwStr := nsfw.Value.String()
		part, err := writer.CreateFormField("nsfw")
		if err != nil {
			return nil, err
		}

		io.Copy(part, strings.NewReader(nsfwStr))
	}

	part, err := writer.CreateFormField("url")
	if err != nil {
		return nil, err
	}

	io.Copy(part, strings.NewReader(urlStr))
	writer.Close()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/shorten_url", rootURL), body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", authToken)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("empty response body")
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var dflE dflerr.E
		err := json.Unmarshal(content, &dflE)
		if err != nil {
			return nil, err
		}

		return nil, dflE
	}

	return content, nil
}
