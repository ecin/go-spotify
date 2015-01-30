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
    redirectURL := "http://localhost:8080"
    scopes := []string{
      "playlist-modify-public",
      "user-read-private",
    }

    fmt.Printf("Get an authorization code: %s\n", AuthorizationURL(credentials, redirectURL, scopes, ""))
    return

    token := &oauth2.Token{
      RefreshToken: "AQDyrREp_5sHhwbIi3HmM-BwCjWKDyU3I9n2LxBkXPcox6ERwNLRD5ZuKHYu2uKgRnpoABAyK1y3N4Z8hESON3v-HVUiSUtozZ9GFIGyBDpawOMAtcQgQ7uMnEeiofB33U0",
      AccessToken: "BQBaWFcxbc2YHQvqAVshFHc5A-1KjDN3Hc9PwwlE5qX-43E1Lb_7_PwFec3ECP1W9AGanyVdv9tJ6y-EAYIR_8HCEpPfcnrLIuDvdseVcVZ-HBEz4scrt8rVxUH8OGa61B_o0PFx04_iF0cn1D9Sj3JMyCUvA9mm71zI0oO9ci4DuY1u8K2UvNZl6wkqQUpnHmE0bci_",
      TokenType: "Bearer",
    }

    client = NewSpotifyClientWithToken(credentials, token)

    code := "AQB1pSq2-gHFmHm8aDitKUaDX03VJRNsiw0sA-N90q4rQM2ehxYzWnFxB1f69hXBrSP1vLIna47qdNBZnP6uLx9Ey-J73ysnr0HO7Wi-V2SSg71xMaymd2CthWm8YFpMl-_InGrbj5LmVNo419N6JQubtX4vNy0p2D6vZHO_FfrRlmuoYHnsQCnX8bSgNiTLEHEW6hv-1lZjdNiUUYGLL0fHt7ir_Zan-XIGnZf2Ry3V_xX7alQ"
    c, err := NewSpotifyClientWithCode(credentials, redirectURL, code)

    if err != nil {
      fmt.Print("ERROR ERROR ERROR\n")
      fmt.Print("%q\n", err)
    }

    client = c
}

func TestGetPlaylists(t *testing.T) {
  setup()

  userProfile, _ := client.GetUserProfile()
  playlists, _ := client.GetPlaylists(userProfile)

  fmt.Printf("%v", playlists)
}

func TestGetUserProfile(t *testing.T) {
  setup()

  userProfile, _ := client.GetUserProfile()

  fmt.Printf("%v", userProfile)
}

func TestGetTracksForPlaylist(t *testing.T) {
  setup()

  userProfile, _ := client.GetUserProfile()
  playlists, _ := client.GetPlaylists(userProfile)
  tracks, _ := client.GetTracksForPlaylist(userProfile, playlists[0])

  fmt.Printf("%v", tracks)
}
