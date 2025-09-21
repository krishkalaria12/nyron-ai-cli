package dialogs

import "github.com/krishkalaria12/nyron-ai-cli/util"

type DialogId string

type DialogModel interface {
	util.Model
}
