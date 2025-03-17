// Package styles содержит определения стилей для текстового интерфейса пользователя, используемого в TUI приложении GophKeeper.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Black определяет черный цвет.
	Black = lipgloss.Color("#000000")

	// Red определяет красный цвет.
	Red = lipgloss.Color("#FF5353")

	// Purple определяет фиолетовый цвет.
	Purple = lipgloss.Color("63")

	// LightGreen определяет светло-зеленый цвет.
	LightGreen = lipgloss.Color("47")

	// Blue определяет синий цвет.
	Blue = lipgloss.Color("63")

	// DeepBlue определяет насыщенный голубой цвет.
	DeepBlue = lipgloss.Color("39")

	// Grey определяет серый цвет.
	Grey = lipgloss.Color("#737373")

	// EvenLighterGrey определяет еще более светлый серый цвет.
	EvenLighterGrey = lipgloss.Color("253")

	// DarkGrey определяет темно-серый цвет.
	DarkGrey = lipgloss.Color("#606362")

	// White определяет белый цвет.
	White = lipgloss.Color("#ffffff")

	// Regular определяет базовый стиль текста.
	Regular = lipgloss.NewStyle()

	// Bold определяет жирный стиль текста.
	Bold = Regular.Bold(true)

	// Padded определяет стиль с внутренними отступами по бокам.
	Padded = Regular.Padding(0, 1)

	// Border определяет стиль рамки с закругленными углами.
	Border = Regular.Border(lipgloss.RoundedBorder())

	// ThickBorder определяет стиль толстой рамки с настраиваемым цветом.
	ThickBorder = Regular.Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.AdaptiveColor{Dark: string(DeepBlue), Light: string(DeepBlue)})

	// ActiveBorder определяет стиль активной рамки.
	ActiveBorder = Border.BorderForeground(lipgloss.AdaptiveColor{Dark: string(DeepBlue), Light: string(DeepBlue)})

	// InactiveBorder определяет стиль неактивной рамки.
	InactiveBorder = Border.BorderForeground(lipgloss.AdaptiveColor{Dark: string(White), Light: string(White)})

	// Focused определяет стиль для элементов в фокусе.
	Focused = Regular.Foreground(lipgloss.AdaptiveColor{Dark: "205", Light: "205"})

	// Blurred определяет стиль для не фокусированных элементов.
	Blurred = Regular.Foreground(lipgloss.AdaptiveColor{Dark: "240", Light: "240"})

	// Highlighted определяет стиль для подсвеченных элементов.
	Highlighted = Regular.Foreground(Purple)

	// HeaderStyle определяет стиль заголовка.
	HeaderStyle = Bold.Foreground(lipgloss.Color("#FF79C6"))

	// ContentPaddedStyle определяет стиль контента с дополнительными отступами.
	ContentPaddedStyle = Regular.Padding(1, 4)

	// HelpKeyStyle определяет стиль для клавиш помощи.
	HelpKeyStyle = Bold.Foreground(lipgloss.AdaptiveColor{Dark: "ff", Light: ""}).Margin(0, 1, 0, 0)

	// HelpDescStyle определяет стиль для описания клавиш помощи.
	HelpDescStyle = Regular.Foreground(lipgloss.AdaptiveColor{Dark: "248", Light: "246"})

	// FilePickerBotPadding определяет размер отступа в нижней части селектора файлов.
	FilePickerBotPadding = 10

	// StorageScreenStyle определяет стиль для общего экрана.
	StorageScreenStyle = Regular.PaddingLeft(2)

	// TableStyle определяет стиль для таблиц.
	TableStyle = Border.BorderForeground(lipgloss.Color("240"))

	// TableSelectedStyle определяет стиль для выбранной строки в таблице.
	TableSelectedStyle = Regular.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)

	// TableHeaderStyle определяет стиль для заголовка таблицы.
	TableHeaderStyle = Padded.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)
)
