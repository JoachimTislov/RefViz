package ops

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func selectPrompt(label string, items interface{}) promptui.Select {
	return promptui.Select{
		Label: label,
		Items: items,
	}
}

func confirmPrompt(msg string) promptui.Prompt {
	return promptui.Prompt{
		Label:       fmt.Sprintf("%s, are you sure you want to continue", msg),
		IsConfirm:   true,
		HideEntered: true,
	}
}

func confirm(msg string) bool {
	p := confirmPrompt(msg)
	v, err := p.Run()
	if err != nil {
		return false
	}
	return v == yes
}
