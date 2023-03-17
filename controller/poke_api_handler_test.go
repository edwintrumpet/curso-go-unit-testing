package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
	c := require.New(t)

	pokemon, err := GetPokemonFromPokeAPI("pikachu")
	c.NoError(err)

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)
}

func TestGetPokemonFromPokeApiSuccessWithMocks(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	id := "pikachu"

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder(http.MethodGet, request, httpmock.NewStringResponder(http.StatusOK, string(body)))

	pokemon, err := GetPokemonFromPokeAPI(id)
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, pokemon)
}

func TestGetPokemonFromPokeApiInternalServerError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	id := "pikachu"

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder(http.MethodGet, request, httpmock.NewStringResponder(http.StatusInternalServerError, string(body)))

	_, err = GetPokemonFromPokeAPI(id)
	c.NotNil(err)
	c.EqualError(ErrPokeApiFailure, err.Error())
}

func TestGetPokemon(t *testing.T) {
	c := require.New(t)

	r, err := http.NewRequest(http.MethodGet, "/pokemon/{id}", nil)
	c.NoError(err)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": "pikachu",
	}

	r = mux.SetURLVars(r, vars)

	GetPokemon(w, r)

	expectedBodyResponse, err := ioutil.ReadFile("samples/api_response.json")
	c.NoError(err)

	var expectedPokemon models.Pokemon

	err = json.Unmarshal([]byte(expectedBodyResponse), &expectedPokemon)
	c.NoError(err)

	var currentPokemon models.Pokemon

	err = json.Unmarshal([]byte(w.Body.Bytes()), &currentPokemon)
	c.NoError(err)

	c.Equal(http.StatusOK, w.Code)
	c.Equal(expectedPokemon, currentPokemon)
}
