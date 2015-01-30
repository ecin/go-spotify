package spotify

import (
  "fmt"
  "testing"
)

import (
  "golang.org/x/oauth2"
)

var client SpotifyClient
var code string

func setup() {
    credentials := Credentials{"fdb47c4e089942a4b8f0ec82586bc0b7", "37e600e32d9347aab9b8f18ac70670b8"}
    //redirectURL := "http://localhost:8080"
/*
    scopes := []string{
      "playlist-modify-public",
      "user-read-private",
    }
*/
    token := &oauth2.Token{
      RefreshToken: "AQDyrREp_5sHhwbIi3HmM-BwCjWKDyU3I9n2LxBkXPcox6ERwNLRD5ZuKHYu2uKgRnpoABAyK1y3N4Z8hESON3v-HVUiSUtozZ9GFIGyBDpawOMAtcQgQ7uMnEeiofB33U0",
      AccessToken: "BQBaWFcxbc2YHQvqAVshFHc5A-1KjDN3Hc9PwwlE5qX-43E1Lb_7_PwFec3ECP1W9AGanyVdv9tJ6y-EAYIR_8HCEpPfcnrLIuDvdseVcVZ-HBEz4scrt8rVxUH8OGa61B_o0PFx04_iF0cn1D9Sj3JMyCUvA9mm71zI0oO9ci4DuY1u8K2UvNZl6wkqQUpnHmE0bci_",
      TokenType: "Bearer",
    }

    client = NewSpotifyClientWithToken(credentials, token)
/*
    fmt.Printf("Get an authorization code: %s\n", AuthorizationURL(credentials, redirectURL, scopes, ""))

    code := "AQCYd_YCcaoeIhSknGB3dQotrwRrcGSHrOT3gGsXNyBWO2C1e6KDiiF5j5IVbYgZ-s0lKWH7eoWDEeN8cA3iwpiJaiUxJt0gY-1lXtvtNa35BhIhU9Vca2H4dj52xc3DGDP9BmWrV90rZVNWlnN-NAw57WFif2KLl5QcbS7AirgHiLlZX8iGir3gIba7_clRed1Q7nVRVqd_wZL3xGTRwSw9zokmVR77H5Bk-avBN2o55Y0CtiQ"
    _, err := NewSpotifyClientWithCode(credentials, redirectURL, code)

    if err != nil {
      fmt.Print("ERROR ERROR ERROR\n")
      fmt.Print("%q\n", err)
    }*/
}

func TestGetPlaylists(t *testing.T) {
  setup()

  userProfile := client.GetUserProfile()
  playlists := client.GetPlaylists(userProfile)

  fmt.Printf("%v", playlists)
}

func TestGetUserProfile(t *testing.T) {
  setup()

  userProfile := client.GetUserProfile()

  fmt.Printf("%v", userProfile)
}

func TestGetTracksForPlaylist(t *testing.T) {
  setup()

  userProfile := client.GetUserProfile()
  playlist := client.GetPlaylists(userProfile)[0]
  tracks := client.GetTracksForPlaylist(userProfile, playlist)

  fmt.Printf("%v", tracks)
}
