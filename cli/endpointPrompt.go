package cli

import (
	"os"

	"github.com/manifoldco/promptui"
)

func PromptForEndpoint(e string) (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Enter the endpoint for the admin console:",
		Default:   e,
		AllowEdit: true,
		Templates: templates,
		Validate: func(input string) error {
			return nil
		},
	}

	for {
		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				os.Exit(-1)
			}
			continue
		}

		return result, nil
	}
}
