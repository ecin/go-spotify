package spotify

import (
  "fmt"
  "os"
  "testing"
)

import (
  "golang.org/x/oauth2"
)

var client SpotifyClient

func setup() (err error) {
  if client.Token != nil {
    return
  }

  id := os.Getenv("SPOTIFY_ID")
  secret := os.Getenv("SPOTIFY_SECRET")
  credentials := Credentials{id, secret}
  token := &oauth2.Token{
    AccessToken: "BQDXahWIzeCQdqrQorQAh5v1GxZZwK4oi7z12qx11Iz6a1KN-xsLJUPksyKwcNAtfQ7SYPlKZCQNY_eGCKFEC_qI5-XjLcOPALm0fjAEA4hZgVSkcIWt8OAbEhA9i9cliZ3Ovm2-CT8hV-CXNbdwiCjuLRnSo2PssUmANAcZBPnSPUU7yD8QwVzJEOYhE3UXIoxrJY8H",
    RefreshToken: "AQCGoqe5ZoB9_rlCtnBKabZvfPedxvgG8ldPpxjnoNDzsJ7LNprEAJMhJGgj-WIO1NPdyURQBvJKn15lCGSvSvB9Pq9BEFF_QnZAGJHpOFc015Le4BGqdXiFt_hp-j2BcBI",
    TokenType: "Bearer",
  }

  client = NewSpotifyClientWithToken(credentials, token)
  return


  code := os.Getenv("SPOTIFY_OAUTH_CODE")
  redirectURL := "http://localhost:8080"

  if code == "" {
    scopes := []string{
      "playlist-modify-public",
      "user-read-private",
    }

    fmt.Printf("Get an authorization code: %s\n", AuthorizationURL(credentials, redirectURL, scopes, ""))
    os.Exit(1)

    return
  } else {
    client, err = NewSpotifyClientWithCode(credentials, redirectURL, code)
    fmt.Printf("Here's your token:\n")
    fmt.Printf("%q", client.Token)

    return
  }
}

func TestGetPlaylists(t *testing.T) {
  setup()

  userProfile, _ := client.GetUserProfile()
  playlists, _ := client.GetPlaylists(userProfile)

  fmt.Println("Here are your playlists:")
  fmt.Printf("%q\n\n", playlists)
}

func TestGetUserProfile(t *testing.T) {
  setup()

  userProfile, _ := client.GetUserProfile()

  fmt.Println("Here is your user profile:")
  fmt.Printf("%q\n\n", userProfile)
}

func TestGetTracksForPlaylist(t *testing.T) {
  setup()

  userProfile, _ := client.GetUserProfile()
  playlists, _ := client.GetPlaylists(userProfile)
  tracks, _ := client.GetTracksForPlaylist(userProfile, playlists[0])

  fmt.Println("Here are your tracks:")
  fmt.Printf("%q\n\n", tracks)
}
