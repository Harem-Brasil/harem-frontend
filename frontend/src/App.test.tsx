import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import App from './App'

describe('App', () => {
  it('renders the Harém Brasil heading', () => {
    render(<MemoryRouter><App /></MemoryRouter>)
    expect(screen.getByRole('heading', { name: /Harém Brasil/i })).toBeInTheDocument()
  })

  it('renders sidebar navigation links', () => {
    render(<MemoryRouter><App /></MemoryRouter>)
    expect(screen.getByRole('link', { name: /Início/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Fórum/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Mensagens/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Assinaturas/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Perfil/i })).toBeInTheDocument()
  })

  it('renders the publish button', () => {
    render(<MemoryRouter><App /></MemoryRouter>)
    expect(screen.getByRole('button', { name: /Nova publicação/i })).toBeInTheDocument()
  })

  it('renders the search input', () => {
    render(<MemoryRouter><App /></MemoryRouter>)
    expect(screen.getByPlaceholderText(/Procurar criadoras, tópicos e posts/i)).toBeInTheDocument()
  })
})
