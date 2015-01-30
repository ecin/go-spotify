package spotify

// Standard libraries
import (
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

func AuthorizationURL(credentials Credentials, redirectURL string, scopes []string, state string) string {
  config := oauth2.Config{
    ClientID: credentials.Id,
    ClientSecret: credentials.Secret,
    Endpoint: Endpoint,
    RedirectURL: redirectURL,
    Scopes: scopes,
  }

  return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func NewSpotifyClientWithCode(credentials Credentials, redirectURL string, code string) (SpotifyClient, error) {
  config := oauth2.Config{
    ClientID: credentials.Id,
    ClientSecret: credentials.Secret,
    Endpoint: Endpoint,
    RedirectURL: redirectURL,
  }

  token, err := config.Exchange(oauth2.NoContext, code)

  if err != nil {
    return SpotifyClient{}, err
  }

  client := NewSpotifyClientWithToken(credentials, token)
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

  return client
}

func (client SpotifyClient) GetUserProfile() (UserProfileResponse, error) {
  request, err := http.NewRequest("GET", USER_PROFILE_ENDPOINT, nil)

  if err != nil {
    return UserProfileResponse{}, err
  }

  response, err := client.makeRequest(request)

  if err != nil {
    return UserProfileResponse{}, err
  }

  body, err := ioutil.ReadAll(response.Body)

  if err != nil {
    return UserProfileResponse{}, err
  }

  var userProfileResponse UserProfileResponse
  json.Unmarshal(body, &userProfileResponse)

  return userProfileResponse, nil
}

func (client SpotifyClient) GetPlaylists(userProfile UserProfileResponse) ([]Playlist, error) {
  endpoint := fmt.Sprintf(PLAYLISTS_ENDPOINT, userProfile.Id)
  request, err := http.NewRequest("GET", endpoint, nil)

  if err != nil {
    return []Playlist{}, err
  }

  response, err := client.makeRequest(request)

  if err != nil {
    return []Playlist{}, err
  }

  body, err := ioutil.ReadAll(response.Body)

  if err != nil {
    return []Playlist{}, err
  }

  var playlistsResponse PlaylistsResponse
  json.Unmarshal(body, &playlistsResponse)

  for _, playlist := range playlistsResponse.Playlists {
    playlist.User = &userProfile
  }

  return playlistsResponse.Playlists, nil
}

func (client SpotifyClient) GetTracksForPlaylist(userProfile UserProfileResponse, playlist Playlist) ([]Track, error) {
  endpoint := fmt.Sprintf(PLAYLIST_TRACKS_ENDPOINT, userProfile.Id, playlist.Id)
  request, err := http.NewRequest("GET", endpoint, nil)

  if err != nil {
    return []Track{}, err
  }

  response, err := client.makeRequest(request)

  if err != nil {
    return []Track{}, err
  }

  body, err := ioutil.ReadAll(response.Body)

  if err != nil {
    return []Track{}, err
  }

  var tracksResponse TracksResponse
  json.Unmarshal(body, &tracksResponse)

  tracks := make([]Track, len(tracksResponse.TrackWrappers))

  for index, trackWrapper := range tracksResponse.TrackWrappers {
    track := trackWrapper.Track
    track.Playlist = &playlist
    tracks[index] = trackWrapper.Track
  }

  return tracks, nil
}

func (client SpotifyClient) makeRequest(request *http.Request) (*http.Response, error) {
  response, err := client.client.Do(request)

  return response, err
}
