package shlink_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/loganintech/shlink-client/shlink"
)

func TestClient_doRequest(t *testing.T) {
	assert.New(t)
	assert.NoError(t, godotenv.Load("../.env"))

	client, err := shlink.NewClient(context.Background(), os.Getenv("API_KEY"), os.Getenv("API_URL"))
	if err != nil {
		return
	}
	linkResp, err := client.CreateShortlink(&shlink.CreateShortlinkRequest{
		LongUrl:      "https://twitch.tv/rocketleague",
		FindIfExists: true,
		ForwardQuery: true,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, linkResp.ShortUrl)
}
