package structs

import (
	"errors"
	"time"

	"github.com/shokHorizon/jsonRunner/comands"
)

type Action struct {
	Name      string            `json:"name"`
	Action    string            `json:"action"`
	Params    map[string]string `json:"params"`
	Result    string            `json:"results"`
	Next      []string          `json:"next"`
	PrevNodes []*Action         `json:"-"`
	NextNodes []*Action         `json:"-"`
}

type Config struct {
	Actions    []Action `json:"actions"`
	Conditions []Action `json:"conditions"`
}

func (act *Action) Exec(nodes map[string]*Action) error {
	params := make(map[string]string)
	for _, prevNode := range act.PrevNodes {
		for k, v := range prevNode.Params {
			params[k] = v
		}
	}
	for k, v := range act.Params {
		params[k] = v
	}

	act.Params = params

	var ok bool
	var err error

	if act.Result != "" {
		return nil
	}

	switch act.Action {
	case "createFile":
		if path, ok := params["file"]; ok {
			return comands.CreateFile(path)
		}
		return errors.New("parameter 'file' not found")
	case "renameFile":
		var (
			path   string
			rename string
		)
		if path, ok = params["file"]; !ok {
			return errors.New("parameter 'file' not found")
		}
		if rename, ok = params["rename"]; !ok {
			return errors.New("parameter 'rename' not found")
		}
		return comands.RenameFile(path, rename)
	case "appendString":
		var (
			path string
			text string
		)
		if path, ok = params["file"]; !ok {
			return errors.New("parameter 'file' not found")
		}
		if text, ok = params["text"]; !ok {
			return errors.New("parameter 'text' not found")
		}
		return comands.AppendFile(path, text)
	case "getCreationTime":
		var (
			path string
		)
		if path, ok = params["file"]; !ok {
			return errors.New("parameter 'file' not found")
		}
		if act.Result, err = comands.GetCreationTime(path); err != nil {
			return err
		}
		return nil
	case "timeFromString":
		var (
			time string
		)
		if time, ok = params["time"]; !ok {
			return errors.New("parameter 'time' not found")
		}
		act.Result = time
		return nil
	case "ifTime":
		var (
			firstTime    time.Time
			secondTime   time.Time
			firstAction  *Action
			secondAction *Action
			operator     string
		)
		if operator, ok = act.Params["operator"]; !ok {
			return errors.New("cannot find 'operator' param")
		}
		if nodeName, ok := act.Params["first_arg"]; ok {
			if firstAction, ok = nodes[nodeName]; !ok {
				return errors.New("first node for if statement is not initialized")
			}
		} else {
			return errors.New("forget about 'first_arg' param")
		}
		if nodeName, ok := act.Params["second_arg"]; ok {
			if secondAction, ok = nodes[nodeName]; !ok {
				return errors.New("first node for if statement is not initialized")
			}
		} else {
			return errors.New("forget about 'second_arg' param")
		}
		if firstAction.Result == "" || secondAction.Result == "" {
			return nil
		}
		if firstTime, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", firstAction.Result); err != nil {
			return err
		}
		if secondTime, err = time.Parse("15:04 02.01.2006", secondAction.Result); err != nil {
			return err
		}
		if ok, err = comands.CompareTimes(firstTime, secondTime, operator); err != nil {
			return err
		}
		if ok {
			act.Result = act.Next[0]
		} else {
			act.Result = act.Next[1]
		}
	}
	return nil
}
