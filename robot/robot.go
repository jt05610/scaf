package robot

type BotKind string

const (
	BotKindPhysical BotKind = "physical"
	BotKindVirtual  BotKind = "virtual"
	BotKindLibrary  BotKind = "library"
)

type Bot struct {
	Name string
	Kind BotKind
}

func NewBot(kind BotKind) *Bot {
	return nil
}

func DefaultPhysicalBot() *Bot {
	return nil
}

func DefaultVirtualBot() *Bot {
	return nil
}
