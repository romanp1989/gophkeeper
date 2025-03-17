// Package tui предоставляет инструменты для управления состоянием пользовательского интерфейса в приложениях с текстовым интерфейсом пользователя.
package tui

import "github.com/charmbracelet/bubbletea"

// Cache управляет кэшем моделей интерфейса пользователя для оптимизации их повторного использования.
type Cache struct {
	cache map[Page]TeaLike
}

// NewCache создает и возвращает новый экземпляр Cache.
func NewCache() *Cache {
	return &Cache{
		cache: make(map[Page]TeaLike),
	}
}

// Get возвращает модель TeaLike, связанную с указанной страницей. Если модель отсутствует, возвращается nil.
func (c *Cache) Get(page Page) TeaLike {
	return c.cache[page]
}

// Put помещает модель TeaLike в кэш с указанным ключом страницы.
func (c *Cache) Put(page Page, model TeaLike) {
	c.cache[page] = model
}

// UpdateAll выполняет обновление всех кэшированных моделей, используя предоставленное сообщение msg.
func (c *Cache) UpdateAll(msg tea.Msg) []tea.Cmd {
	commands := make([]tea.Cmd, len(c.cache))
	var i int
	for k := range c.cache {
		commands[i] = c.Update(k, msg)
		i++
	}
	return commands
}

// Update обновляет модель, ассоциированную с указанным ключом страницы, используя предоставленное сообщение msg.
func (c *Cache) Update(key Page, msg tea.Msg) tea.Cmd {
	return c.cache[key].Update(msg)
}
