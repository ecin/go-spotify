package spotify

import (
  "fmt"
  "testing"
)

var client SpotifyClient
var code string

func setup() {
  if client.redirectURL == "" {
    credentials := Credentials{"fdb47c4e089942a4b8f0ec82586bc0b7", "37e600e32d9347aab9b8f18ac70670b8"}
    redirectURL := "http://localhost:8080"

    scopes := []string{
      "playlist-modify-public",
      "user-read-private",
    }

    fmt.Printf("Get an authorization code: %s\n", AuthorizationURL(credentials, redirectURL, scopes))
    return
    // Regenerate this tokenResponse if it expires (or implement refreshToken())
    tokenResponse := TokenResponse{
      AccessToken: "BQDCwHn4c5h2xFpf2SjOeHbXH1DyxOD1OZ-EZQXT2LPknE1nqXYhdC-GLMRekH_XeRFxYPgtB5IXjxcsQosFF0GX7GpK2Nrjt6XZlkWZfA4CMbVglyjBU-DtEso8OYaSGq_ISA9WF950iNRY5-X3T2lekOiRqoBcQgfyrFcnEsSLxiXqV-E5905yzMzmkA_v",
      ExpiresIn: 3600,
      RefreshToken: "AQBz2wEPsHpzF7DeFQRWA1T7WxG1JU9Sx8vEaz2x0_yEy2619WQast41kV2-pS0HLR0FmQTBaUWzGhnjQqWtkJNIne6mSoJdvCNKbseQL7g5c9SSjhJ81gJzWkfcDP1N2D0",
      TokenType: "Bearer",
    }
    client = NewSpotifyClientWithTokenResponse(credentials, redirectURL, tokenResponse)
  }
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
