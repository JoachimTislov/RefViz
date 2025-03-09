package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	yes = "y"
)

func Confirm(msg string) bool {
	p := confirmPrompt(msg)
	v, err := p.Run()
	if err != nil {
		return false
	}
	return v == yes
}

func confirmPrompt(msg string) promptui.Prompt {
	return promptui.Prompt{
		Label:       fmt.Sprintf("%s, are you sure you want to continue", msg),
		IsConfirm:   true,
		HideEntered: true,
	}
}
