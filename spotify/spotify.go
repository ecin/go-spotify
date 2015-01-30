package spotify

import (
  "encoding/base64"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
)

// Third-party libraries
import (
  "golang.org/x/oauth2"
)

const (
  PLAYLISTS_ENDPOINT = "https://api.spotify.com/v1/users/%s/playlists"
  PLAYLIST_TRACKS_ENDPOINT = "https://api.spotify.com/v1/users/%s/playlists/%s/tracks"
  USER_PROFILE_ENDPOINT = "https://api.spotify.com/v1/me"
)

var Endpoint = oauth2.Endpoint{
  AuthURL: "https://accounts.spotify.com/authorize",
  TokenURL: "https://accounts.spotify.com/api/token",
}

type SpotifyClient struct {
  credentials Credentials
  client *http.Client
}

// TODO: verify these properties need to be public
type TokenResponse struct {
  AccessToken string `json:"access_token"`
  ExpiresIn int64 `json:"expires_in"`
  RefreshToken string `json:"refresh_token"`
  TokenType string `json:"token_type"`
}

type Playlist struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Tracks struct {
    Total int64 `json:"total"`
    Href string `json:"href"`
  } `json:"tracks"`
  Href string `json:"href"`
  URI string `json:"uri"`
  User *UserProfileResponse
}

type PlaylistsResponse struct {
  Playlists []Playlist `json:"items"`
}

type UserProfileResponse struct {
  Id string `json:"id"`
  Name string `json:"display_name"`
  Email string `json:"email"`
  URI string `json:"uri"`
  //Images []Image `json:"images"`
}

type TrackWrapper struct {
  AddedAt string `json:"added_at"`
  AddedBy string `json:"added_by"`
  Track Track `json:"track"`
}

type TracksResponse struct {
  TrackWrappers []TrackWrapper `json:"items"`
  Total int64 `json:"total"`
}

type Artist struct {
  Id string `json:"id"`
  Name string `json:"name"`
}

type Track struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Artists []Artist `json:"artists"`
  URI string `json:"uri"`
  Playlist *Playlist
}

type Credentials struct {
  Id string
  Secret string
}

func (credentials Credentials) getSignature() string {
   data := []byte(fmt.Sprintf("%s:%s", credentials.Id, credentials.Secret))
   signature := base64.URLEncoding.EncodeToString(data)

   return signature
}

func AuthorizationURL(credentials Credentials, redirectURL string, scopes []string, state string) string {
  config := oauth2.Config{
    ClientID: credentials.Id,
    ClientSecret: credentials.Secret,
    Endpoint: Endpoint,
    RedirectURL: redirectURL,
    Scopes: scopes,
  }

  return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
/*
  // TODO: use state parameter by returning two strings
  // https://developer.spotify.com/web-api/authorization-guide/
  params := url.Values{}
  params.Add("client_id", credentials.Id)
  params.Add("response_type", "code")
  params.Add("scope", strings.Join(scopes, " "))
  params.Add("redirect_uri", redirectURL)

  return fmt.Sprintf("%s?%s", AUTHORIZATION_ENDPOINT, params.Encode())
*/
}

func NewSpotifyClientWithCode(credentials Credentials, redirectURL string, code string) (*SpotifyClient, error) {
  config := oauth2.Config{
    ClientID: credentials.Id,
    ClientSecret: credentials.Secret,
    Endpoint: Endpoint,
    RedirectURL: redirectURL,
  }

  token, err := config.Exchange(oauth2.NoContext, code)

  if err != nil {
    return nil, err
  }

  fmt.Printf("TOKEN TOKEN TOKEN\n")
  fmt.Printf("%q", token)

  client := &SpotifyClient{
    credentials: credentials,
    client: config.Client(oauth2.NoContext, token),
  }

  return client, nil
}

func NewSpotifyClientWithToken(credentials Credentials, token *oauth2.Token) SpotifyClient {
  config := oauth2.Config{
    ClientID: credentials.Id,
    ClientSecret: credentials.Secret,
    Endpoint: Endpoint,
  }

  client := SpotifyClient{
    credentials: credentials,
    client: config.Client(oauth2.NoContext, token),
  }

  fmt.Printf("OAUTH TOKEN: %q\n", token)

  return client
}

func (client SpotifyClient) refreshToken() string {
  //params := url.Values{}
  //params.Add("grant_type", "refresh_token")
  //params.Add("refresh_token", client.tokenResponse.RefreshToken)
  //params.Add("client_id", client.credentials.id)
  //params.Add("client_secret", client.credentials.secret)

  //fmt.Printf("TokenResponse: %v\n", client.tokenResponse)

  //httpClient := &http.Client{}
  //request, _ := http.NewRequest("POST", ACCESS_TOKEN_ENDPOINT, nil)
  ////request.Header.Add("Authorization", fmt.Sprintf("Basic %s", client.credentials.getSignature()))

  //response, _ := httpClient.Do(request)
  //body, _ := ioutil.ReadAll(response.Body)

  //fmt.Printf("refreshToken():\n%s\n", body)

  //var tokenResponse TokenResponse
  //json.Unmarshal(body, &tokenResponse)

  //client.tokenResponse = tokenResponse
  //return tokenResponse.AccessToken
  return ""
}

func (client SpotifyClient) GetUserProfile() UserProfileResponse {
  request, _ := http.NewRequest("GET", USER_PROFILE_ENDPOINT, nil)
  response, _ := client.makeRequest(request)

  body, _ := ioutil.ReadAll(response.Body)
  fmt.Printf("%s\n", body)

  var userProfileResponse UserProfileResponse
  json.Unmarshal(body, &userProfileResponse)

  return userProfileResponse
}

func (client SpotifyClient) GetPlaylists(userProfile UserProfileResponse) []Playlist {
  endpoint := fmt.Sprintf(PLAYLISTS_ENDPOINT, userProfile.Id)
  request, _ := http.NewRequest("GET", endpoint, nil)

  response, _ := client.makeRequest(request)

  body, _ := ioutil.ReadAll(response.Body)
  fmt.Printf("%s\n", body)

  var playlistsResponse PlaylistsResponse
  json.Unmarshal(body, &playlistsResponse)

  for _, playlist := range playlistsResponse.Playlists {
    playlist.User = &userProfile
  }

  return playlistsResponse.Playlists
}

func (client SpotifyClient) GetTracksForPlaylist(userProfile UserProfileResponse, playlist Playlist) []Track {
  endpoint := fmt.Sprintf(PLAYLIST_TRACKS_ENDPOINT, userProfile.Id, playlist.Id)
  request, _ := http.NewRequest("GET", endpoint, nil)

  response, _ := client.makeRequest(request)

  body, _ := ioutil.ReadAll(response.Body)
  fmt.Printf("%s\n", body)

  var tracksResponse TracksResponse
  json.Unmarshal(body, &tracksResponse)

  tracks := make([]Track, len(tracksResponse.TrackWrappers))

  for index, trackWrapper := range tracksResponse.TrackWrappers {
    track := trackWrapper.Track
    track.Playlist = &playlist
    tracks[index] = trackWrapper.Track
  }

  return tracks
}

func (client SpotifyClient) makeRequest(request *http.Request) (*http.Response, error) {
  response, err := client.client.Do(request)

  return response, err
}
