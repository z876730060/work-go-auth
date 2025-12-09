package cloud

import "github.com/z876730060/auth/internal/service"

var (
	RegisterManagerInstance = &RegisterManager{}
)

type Register interface {
	Register(cfg service.Config) error
	Unregister(cfg service.Config) error
}

type RegisterManager struct {
	registers []Register
}

func (m *RegisterManager) AddRegister(r Register) {
	m.registers = append(m.registers, r)
}

func (m *RegisterManager) Register(cfg service.Config) error {
	for _, r := range m.registers {
		if err := r.Register(cfg); err != nil {
			return err
		}
	}
	return nil
}

func (m *RegisterManager) Unregister(cfg service.Config) error {
	for _, r := range m.registers {
		if err := r.Unregister(cfg); err != nil {
			return err
		}
	}
	return nil
}
