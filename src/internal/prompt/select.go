package prompt

import "github.com/manifoldco/promptui"

func SelectPrompt(label string, items any) promptui.Select {
	return promptui.Select{
		Label: label,
		Items: items,
	}
}
