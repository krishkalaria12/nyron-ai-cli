package components

import "github.com/charmbracelet/bubbles/key"

type EditorKeyMap struct {
	SendMessage key.Binding
	Newline     key.Binding
}

func DefaultEditorKeyMap() EditorKeyMap {
	return EditorKeyMap{
		SendMessage: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "send"),
		),
		Newline: key.NewBinding(
			key.WithKeys("shift+enter", "ctrl+j"),
			key.WithHelp("ctrl+j", "newline"),
		),
	}
}

func (k EditorKeyMap) KeyBindings() []key.Binding {
	return []key.Binding{
		k.SendMessage,
		k.Newline,
	}
}
