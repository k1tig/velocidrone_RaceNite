func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	//var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.state == tableView {
				m.table.Blur()
				m.state = formView
			} else {
				m.table.Focus()
				m.state = tableView
				return m, cmd
			}
		}
		switch m.state {
		case formView:
			form, cmd := m.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				m.form = f
				m.form = f //something messed up here <~~~~
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case tableView:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	//m.table, cmd = m.table.Update(msg)
	return m, tea.Batch(cmds...)

}
