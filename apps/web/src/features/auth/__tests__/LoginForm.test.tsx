import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { server } from '@/test/msw/server'
import { AuthProvider, LoginForm } from '@/features/auth'

function renderLogin() {
  return render(
    <MemoryRouter initialEntries={['/login']}>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginForm />} />
          <Route path="/" element={<div>Tableau de bord</div>} />
        </Routes>
      </AuthProvider>
    </MemoryRouter>,
  )
}

describe('LoginForm', () => {
  it('connecte et redirige vers le tableau de bord', async () => {
    server.use(
      http.post('/api/auth/login', () =>
        HttpResponse.json({ username: 'graefik' }),
      ),
    )
    const user = userEvent.setup()
    renderLogin()

    await user.type(screen.getByLabelText(/mot de passe/i), 'secret')
    await user.click(screen.getByRole('button', { name: /se connecter/i }))

    expect(await screen.findByText('Tableau de bord')).toBeInTheDocument()
  })

  it('affiche une erreur sur identifiants invalides', async () => {
    server.use(
      http.post('/api/auth/login', () => new HttpResponse(null, { status: 401 })),
    )
    const user = userEvent.setup()
    renderLogin()

    await user.type(screen.getByLabelText(/mot de passe/i), 'mauvais')
    await user.click(screen.getByRole('button', { name: /se connecter/i }))

    expect(await screen.findByText(/identifiants invalides/i)).toBeInTheDocument()
  })
})
