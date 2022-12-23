package cli

import (
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
)

func PromptForTlsVerify(tls bool) (bool, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Verify TLS:",
		Default:   strconv.FormatBool(tls),
		AllowEdit: true,
		Templates: templates,
		Validate: func(input string) error {
			_, err := strconv.ParseBool(input)
			return err
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

		return strconv.ParseBool(result)
	}
}
