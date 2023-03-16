package util

import (
	"catching-pokemons/models"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePokemonSuccess(t *testing.T) {
	c := require.New(t)

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	var response models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &response)
	c.NoError(err)

	parsedPokemon, err := ParsePokemon(response)
	c.NoError(err)

	body, err = ioutil.ReadFile("samples/api_response.json")
	c.NoError(err)

	var expected models.Pokemon

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)

	c.Equal(expected, parsedPokemon)
}

func TestParsePokemonTypeNotFound(t *testing.T) {
	c := require.New(t)

	body, err := ioutil.ReadFile("samples/pokeapi_response.json")
	c.NoError(err)

	var response models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &response)
	c.NoError(err)

	response.PokemonType = []models.PokemonType{}

	_, err = ParsePokemon(response)
	c.NotNil(err)
	c.EqualError(ErrNotFoundPokemonType, err.Error())
}
