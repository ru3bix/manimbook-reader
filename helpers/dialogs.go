package helpers

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func GetManimBookFile(ctx context.Context) (string, error) {
	result, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select manimbook",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "manimbooks (*.mbook)",
				Pattern:     "*.mbook",
			},
		},
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", fmt.Errorf("no file selected")
	}
	return result, nil
}

func DisplayErrorMsg(ctx context.Context, e error) error {
	_, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.ErrorDialog,
		Title:         "Error",
		Message:       e.Error(),
		Buttons:       []string{"ok"},
		DefaultButton: "ok",
	})
	return err
}
