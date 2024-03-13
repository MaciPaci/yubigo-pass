package test

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TeaModelMock is a mock of Bubbletea Model interface
type TeaModelMock struct {
	returnedModels []tea.Model
	iterations     int
}

// NewTeaModelMock returns pointer to the new instance of TeaModelMock
func NewTeaModelMock() *TeaModelMock {
	return &TeaModelMock{}
}

// WillReturnOnce mocks one return of a BubbleTea model's Update function
func (m *TeaModelMock) WillReturnOnce(model tea.Model) *TeaModelMock {
	m.returnedModels = append(m.returnedModels, model)
	return m
}

// Init mocks Bubbletea Model's init function
func (m *TeaModelMock) Init() tea.Cmd {
	return func() tea.Msg {
		return struct{}{}
	}
}

// Update mocks Bubbletea Model's update function
func (m *TeaModelMock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	returnModel := m.returnedModels[m.iterations]
	m.iterations++
	return returnModel, tea.Quit
}

// View mocks Bubbletea Model's view function
func (m *TeaModelMock) View() string {
	return ""
}
