package cli

import (
	"errors"
	"fmt"
	"strings"
	"yubigo-pass/internal/app/common"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/app/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// AppModel is the main application model responsible for managing different views (sub-models)
// and orchestrating the overall application flow based on messages.
type AppModel struct {
	width       int
	height      int
	activeModel tea.Model
	session     utils.Session
	container   services.Container
	lastError   error
}

// NewAppModel creates the initial state of the top-level application model.
// It now creates the initial CLI model directly.
func NewAppModel(container services.Container) AppModel {
	// Create the initial login model here, passing the store from the container
	initialModel := NewLoginModel(container.Store)

	return AppModel{
		container:   container,
		session:     utils.NewEmptySession(),
		activeModel: initialModel, // Start with the login model
	}
}

// Init initializes the application model by initializing the currently active sub-model.
func (m AppModel) Init() tea.Cmd {
	if m.activeModel != nil {
		return m.activeModel.Init()
	}
	return nil
}

// Update handles incoming messages, updates the active sub-model, and manages transitions
// between different application views based on custom messages. It creates new models as needed.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case common.ErrorMsg:
		m.lastError = msg.Err
		return m, nil

	case common.StateMsg:
		m.lastError = nil
		switch msg.State {
		case common.StateGoToCreateUser:
			m.activeModel = NewCreateUserModel(m.container.Store)
			return m, m.activeModel.Init()
		case common.StateGoToMainMenu:
			if !m.session.IsAuthenticated() {
				cmds = append(cmds, common.ErrCmd(errors.New("cannot enter main menu: not authenticated")))
				m.activeModel = NewLoginModel(m.container.Store)
				return m, tea.Batch(m.activeModel.Init(), tea.Batch(cmds...))
			}
			m.activeModel = NewMainMenuModel()
			return m, m.activeModel.Init()
		case common.StateGoToAddPassword:
			if !m.session.IsAuthenticated() {
				cmds = append(cmds, common.ErrCmd(errors.New("cannot add password: not authenticated")))
				m.activeModel = NewLoginModel(m.container.Store)
				return m, tea.Batch(m.activeModel.Init(), tea.Batch(cmds...))
			}
			m.activeModel = NewAddPasswordModel(m.container.Store, m.session)
			return m, m.activeModel.Init()

		case common.StateGoBack:
			switch m.activeModel.(type) {
			case AddPasswordModel:
				m.activeModel = NewMainMenuModel()
			case CreateUserModel:
				m.activeModel = NewLoginModel(m.container.Store)
			default:
				m.activeModel = NewLoginModel(m.container.Store)
			}
			return m, m.activeModel.Init()

		case common.StateLogout:
			m.session.Clear()
			m.activeModel = NewLoginModel(m.container.Store)
			return m, m.activeModel.Init()

		case common.StateQuit:
			return m, tea.Quit

		case common.StateUserCreated:
			m.activeModel = NewLoginModel(m.container.Store)
			return m, m.activeModel.Init()
		case common.StatePasswordAdded:
			m.activeModel = NewMainMenuModel()
			return m, m.activeModel.Init()
		}

	case common.LoginSuccessMsg:
		m.lastError = nil
		m.session = msg.Session
		m.activeModel = NewMainMenuModel()
		return m, m.activeModel.Init()

	case common.UserToCreateMsg:
		m.lastError = nil
		err := m.createNewUser(msg.Username, msg.Password)
		if err != nil {
			return m, common.ErrCmd(fmt.Errorf("failed to create user: %w", err))
		}
		return m, common.ChangeStateCmd(common.StateUserCreated)

	case common.PasswordToAddMsg:
		m.lastError = nil
		err := m.addNewPassword(msg.Data.Title, msg.Data.Username, msg.Data.Password, msg.Data.Url)
		if err != nil {
			return m, common.ErrCmd(fmt.Errorf("failed to add password: %w", err))
		}
		return m, common.ChangeStateCmd(common.StatePasswordAdded)

	default:
		if m.activeModel != nil {
			var updatedModel tea.Model
			updatedModel, cmd = m.activeModel.Update(msg)
			m.activeModel = updatedModel
			cmds = append(cmds, cmd)
		}
	}

	sizeMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
	if m.activeModel != nil {
		var updatedModel tea.Model
		updatedModel, cmd = m.activeModel.Update(sizeMsg)
		m.activeModel = updatedModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the view of the currently active sub-model, optionally prepending an error message.
func (m AppModel) View() string {
	var viewBuilder strings.Builder

	if m.lastError != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).MarginBottom(1)
		viewBuilder.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.lastError)))
		viewBuilder.WriteString("\n")
	}

	if m.activeModel != nil {
		viewBuilder.WriteString(m.activeModel.View())
	} else {
		viewBuilder.WriteString("Error: No active model to display.")
	}

	return viewBuilder.String()
}

// createNewUser handles the logic for creating a new user entry in the database.
func (m *AppModel) createNewUser(username, password string) error {
	userUUID := uuid.New().String()
	salt, err := crypto.NewSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}
	passwordHash := crypto.HashPasswordWithSalt(password, salt)

	createUserInput := model.NewUser(userUUID, username, passwordHash, salt)
	err = m.container.Store.CreateUser(createUserInput)
	if err != nil {
		var userExistsError *model.UserAlreadyExistsError
		if errors.As(err, &userExistsError) {
			return userExistsError
		}
		return fmt.Errorf("database error creating user: %w", err)
	}
	return nil
}

// addNewPassword handles the logic for encrypting and adding a new password entry to the database.
func (m *AppModel) addNewPassword(title, username, password, url string) error {
	if !m.session.IsAuthenticated() {
		return errors.New("cannot add password: no active user session")
	}

	encryptionKey := crypto.DeriveAESKey(m.session.GetPassphrase(), m.session.GetSalt())

	encryptedPassword, nonce, err := crypto.EncryptAES(encryptionKey, []byte(password))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	ciphertext := append(nonce, encryptedPassword...)

	addPasswordInput := model.NewPassword(
		m.session.GetUserID(),
		title,
		username,
		string(ciphertext),
		url,
		nonce,
	)

	err = m.container.Store.AddPassword(addPasswordInput)
	if err != nil {
		var passExistsError *model.PasswordAlreadyExistsError
		if errors.As(err, &passExistsError) {
			return passExistsError
		}
		return fmt.Errorf("database error adding password: %w", err)
	}

	return nil
}
