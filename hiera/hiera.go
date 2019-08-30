package hiera

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cast"
)

type Hiera struct {
	Bin    string
	Config string
	Scope  map[string]interface{}
	Merge  string
}

func NewHiera(bin string, config string, scope map[string]interface{}, merge string) Hiera {
	return Hiera{
		Bin:    bin,
		Config: config,
		Scope:  scope,
		Merge:  merge,
	}
}

func (h *Hiera) Exec(arg ...string) ([]byte, error) {
	var scope []string
	var args []string
	var out []byte
	var err error

	for key, value := range h.Scope {
		scope = append(scope, strings.Join([]string{key, value.(string)}, "="))
	}

	for _, c := range [][]string{[]string{"--render-as", "json", "--config", h.Config, "--merge", h.Merge}, arg, scope} {
		args = append(args, c...)
	}

	log.Printf("[DEBUG] Exec args are %v", args)
	if out, err = exec.Command(h.Bin, args...).Output(); err != nil {
		log.Printf("[DEBUG] out is %s", string(out))
		return out, err
	}

	return out, nil
}

func (h *Hiera) Array(key string) ([]interface{}, error) {
	var f interface{}
	var e []interface{}

	out, err := h.Exec("--type", "Array", key)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		return nil, err
	}
	if _, ok := f.([]interface{}); ok {
		for _, v := range f.([]interface{}) {
			e = append(e, cast.ToString(v))
		}
	} else {
		return nil, fmt.Errorf("Key '%s' does not return a valid array", key)
	}

	return e, nil
}

func (h *Hiera) Hash(key string) (map[string]interface{}, error) {
	var f interface{}

	e := make(map[string]interface{})

	out, err := h.Exec("--type", "Hash", key)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		return nil, err
	}

	if _, ok := f.(map[string]interface{}); ok {
		for k, v := range f.(map[string]interface{}) {
			e[k] = cast.ToString(v)
		}
	} else {
		return nil, fmt.Errorf("Key '%s' does not return a valid hash", key)
	}
	return e, nil
}

func (h *Hiera) Value(key string) (string, error) {
	var f interface{}

	out, err := h.Exec(key)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(out, &f)
	if err != nil {
		return "", err
	}

	return cast.ToString(f), nil
}
