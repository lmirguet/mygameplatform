import { useEffect, useMemo, useState, type FormEvent } from 'react'
import './App.css'

type AuthMode = 'signup' | 'login' | 'profile'

type ApiError = {
  error?: string
  message?: string
}

type Profile = {
  username: string
  avatar_url?: string | null
}

type Lobby = {
  id: string
  game: string
  seats: string
  status: string
  action: string
}

const TOKEN_STORAGE_KEY = 'mygameplatform_token'

const usernameRules = {
  min: 3,
  max: 32,
  pattern: /^[a-z0-9._-]+$/,
}

const defaultAvatar = (username: string) => {
  const initial = username.trim().slice(0, 1).toUpperCase() || 'U'
  return initial
}

const lobbyPreview: Lobby[] = [
  {
    id: 'lobby-1',
    game: 'Connect4',
    seats: '1 / 2',
    status: '1 open seat',
    action: 'Join lobby',
  },
  {
    id: 'lobby-2',
    game: 'Draughts (10x10)',
    seats: '2 / 2',
    status: 'Full',
    action: 'Full',
  },
  {
    id: 'lobby-3',
    game: 'Connect4',
    seats: '0 / 2',
    status: 'New table',
    action: 'Create lobby',
  },
]

function App() {
  const [mode, setMode] = useState<AuthMode>('signup')
  const [token, setToken] = useState<string | null>(() =>
    window.localStorage.getItem(TOKEN_STORAGE_KEY)
  )
  const [profile, setProfile] = useState<Profile | null>(null)
  const [status, setStatus] = useState<{ tone: 'info' | 'error' | 'success'; message: string } | null>(null)

  const [signupEmail, setSignupEmail] = useState('')
  const [signupPassword, setSignupPassword] = useState('')
  const [signupLoading, setSignupLoading] = useState(false)

  const [loginEmail, setLoginEmail] = useState('')
  const [loginPassword, setLoginPassword] = useState('')
  const [loginLoading, setLoginLoading] = useState(false)

  const [profileUsername, setProfileUsername] = useState('')
  const [profileAvatar, setProfileAvatar] = useState('')
  const [profileLoading, setProfileLoading] = useState(false)

  const isAuthed = Boolean(token)

  const apiBase = useMemo(() => {
    const envBase = import.meta.env.VITE_API_BASE as string | undefined
    return envBase ? envBase.replace(/\/$/, '') : ''
  }, [])

  useEffect(() => {
    if (token) {
      setMode('profile')
      void fetchProfile()
    } else {
      setProfile(null)
      setMode('signup')
    }
  }, [token])

  useEffect(() => {
    if (profile) {
      setProfileUsername(profile.username)
      setProfileAvatar(profile.avatar_url ?? '')
    }
  }, [profile])

  const fetchProfile = async () => {
    try {
      const data = await apiRequest<Profile>('GET', '/api/v1/me')
      setProfile(data)
      setStatus(null)
    } catch (err) {
      handleError(err, true)
    }
  }

  const handleError = (err: unknown, isProfile = false) => {
    const fallback = isProfile
      ? 'Unable to load profile. Please sign in again.'
      : 'Something went wrong. Please try again.'
    const message = err instanceof Error ? err.message : fallback
    setStatus({ tone: 'error', message })
  }

  const apiRequest = async <T,>(method: string, path: string, body?: unknown): Promise<T> => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }

    const response = await fetch(`${apiBase}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    })

    const contentType = response.headers.get('content-type') || ''
    if (!response.ok) {
      let message = 'Request failed'
      if (contentType.includes('application/json')) {
        const error = (await response.json()) as ApiError
        message = error.message || error.error || message
      }
      if (response.status === 401 || response.status === 403) {
        clearSession()
        message = 'Session expired. Please sign in again.'
      }
      throw new Error(message)
    }

    if (contentType.includes('application/json')) {
      return (await response.json()) as T
    }
    return {} as T
  }

  const saveToken = (nextToken: string) => {
    window.localStorage.setItem(TOKEN_STORAGE_KEY, nextToken)
    setToken(nextToken)
  }

  const clearSession = () => {
    window.localStorage.removeItem(TOKEN_STORAGE_KEY)
    setToken(null)
    setProfile(null)
    setProfileUsername('')
    setProfileAvatar('')
  }

  const handleSignup = async (event: FormEvent) => {
    event.preventDefault()
    setSignupLoading(true)
    setStatus(null)
    try {
      const data = await apiRequest<{ access_token: string }>('POST', '/api/v1/auth/signup', {
        email: signupEmail,
        password: signupPassword,
      })
      saveToken(data.access_token)
      setSignupEmail('')
      setSignupPassword('')
      setStatus({ tone: 'success', message: 'Account created. You are signed in.' })
    } catch (err) {
      handleError(err)
    } finally {
      setSignupLoading(false)
    }
  }

  const handleLogin = async (event: FormEvent) => {
    event.preventDefault()
    setLoginLoading(true)
    setStatus(null)
    try {
      const data = await apiRequest<{ access_token: string }>('POST', '/api/v1/auth/login', {
        email: loginEmail,
        password: loginPassword,
      })
      saveToken(data.access_token)
      setLoginPassword('')
      setStatus({ tone: 'success', message: 'Welcome back. You are signed in.' })
    } catch (err) {
      handleError(err)
    } finally {
      setLoginLoading(false)
    }
  }

  const handleProfileUpdate = async (event: FormEvent) => {
    event.preventDefault()
    setProfileLoading(true)
    setStatus(null)
    try {
      const payload = {
        username: profileUsername,
        avatar_url: profileAvatar,
      }
      const data = await apiRequest<Profile>('PATCH', '/api/v1/me', payload)
      setProfile(data)
      setStatus({ tone: 'success', message: 'Profile updated.' })
    } catch (err) {
      handleError(err)
    } finally {
      setProfileLoading(false)
    }
  }

  const usernameHint = `Lowercase letters, numbers, dot, underscore, dash. ${usernameRules.min}-${usernameRules.max} chars.`

  return (
    <div className="app">
      <header className="topbar">
        <div className="brand">
          <span className="brand-mark">MG</span>
          <div>
            <p className="brand-title">MyGamePlatform</p>
            <p className="brand-subtitle">Play now. No friction.</p>
          </div>
        </div>
        <nav className="nav">
          <button
            className={mode === 'signup' ? 'nav-button active' : 'nav-button'}
            onClick={() => setMode('signup')}
          >
            Sign up
          </button>
          <button
            className={mode === 'login' ? 'nav-button active' : 'nav-button'}
            onClick={() => setMode('login')}
          >
            Sign in
          </button>
          <button
            className={mode === 'profile' ? 'nav-button active' : 'nav-button'}
            onClick={() => setMode('profile')}
            disabled={!isAuthed}
          >
            Profile
          </button>
        </nav>
      </header>

      <main className="main">
        <section className="hero">
          <div className="hero-card">
            <p className="eyebrow">Instant play</p>
            <h1>From landing to seat in minutes.</h1>
            <p className="hero-copy">
              Create an account, grab a seat, and jump into Connect4 or Draughts. No paywalls. No
              clutter. Just a clear path to play.
            </p>
            <div className="hero-stats">
              <div>
                <p className="stat-title">Target</p>
                <p className="stat-value">2-4 min to play</p>
              </div>
              <div>
                <p className="stat-title">Platforms</p>
                <p className="stat-value">Desktop + Mobile</p>
              </div>
            </div>
            <div className="hero-badges">
              <span>Free forever</span>
              <span>Invite link friendly</span>
              <span>Fair play</span>
            </div>
          </div>
          <div className="gradient-sheen" />
        </section>

        <section className="panel">
          <div className="panel-header">
            <div>
              <h2>{mode === 'profile' ? 'Your profile' : mode === 'login' ? 'Welcome back' : 'Create your account'}</h2>
              <p>
                {mode === 'profile'
                  ? 'Update your public handle before you play.'
                  : 'Minimal fields. Friendly errors. Fast track to the lobby.'}
              </p>
            </div>
            {isAuthed && (
              <button className="ghost" onClick={clearSession}>
                Sign out
              </button>
            )}
          </div>

          {status && (
            <div className={`status ${status.tone}`}>{status.message}</div>
          )}

          {mode === 'signup' && (
            <form className="form" onSubmit={handleSignup}>
              <label>
                Email
                <input
                  type="email"
                  value={signupEmail}
                  onChange={(event) => setSignupEmail(event.target.value)}
                  placeholder="you@example.com"
                  required
                />
              </label>
              <label>
                Password
                <input
                  type="password"
                  value={signupPassword}
                  onChange={(event) => setSignupPassword(event.target.value)}
                  placeholder="At least 8 characters"
                  minLength={8}
                  maxLength={72}
                  required
                />
              </label>
              <button className="primary" type="submit" disabled={signupLoading}>
                {signupLoading ? 'Creating account…' : 'Sign up & play'}
              </button>
              <p className="form-footer">
                Already have an account?{' '}
                <button className="link" type="button" onClick={() => setMode('login')}>
                  Sign in
                </button>
              </p>
            </form>
          )}

          {mode === 'login' && (
            <form className="form" onSubmit={handleLogin}>
              <label>
                Email
                <input
                  type="email"
                  value={loginEmail}
                  onChange={(event) => setLoginEmail(event.target.value)}
                  placeholder="you@example.com"
                  required
                />
              </label>
              <label>
                Password
                <input
                  type="password"
                  value={loginPassword}
                  onChange={(event) => setLoginPassword(event.target.value)}
                  placeholder="Your password"
                  required
                />
              </label>
              <button className="primary" type="submit" disabled={loginLoading}>
                {loginLoading ? 'Signing in…' : 'Sign in'}
              </button>
              <p className="form-footer">
                New here?{' '}
                <button className="link" type="button" onClick={() => setMode('signup')}>
                  Create an account
                </button>
              </p>
            </form>
          )}

          {mode === 'profile' && (
            <div className="profile">
              {!isAuthed ? (
                <div className="empty-state">
                  <p>Sign in to view your profile.</p>
                  <button className="primary" onClick={() => setMode('login')}>
                    Sign in
                  </button>
                </div>
              ) : (
                <>
                  <div className="profile-card">
                    <div className="avatar">
                      {profile?.avatar_url ? (
                        <img src={profile.avatar_url} alt="Avatar" />
                      ) : (
                        <span>{defaultAvatar(profileUsername || 'U')}</span>
                      )}
                    </div>
                    <div>
                      <p className="profile-label">Signed in as</p>
                      <p className="profile-username">{profile?.username || 'loading…'}</p>
                    </div>
                  </div>
                  <form className="form" onSubmit={handleProfileUpdate}>
                    <label>
                      Username
                      <input
                        type="text"
                        value={profileUsername}
                        onChange={(event) => setProfileUsername(event.target.value)}
                        placeholder="player_one"
                        minLength={usernameRules.min}
                        maxLength={usernameRules.max}
                        pattern={usernameRules.pattern.source}
                        required
                      />
                      <span className="hint">{usernameHint}</span>
                    </label>
                    <label>
                      Avatar URL (optional)
                      <input
                        type="url"
                        value={profileAvatar}
                        onChange={(event) => setProfileAvatar(event.target.value)}
                        placeholder="https://..."
                      />
                    </label>
                    <button className="primary" type="submit" disabled={profileLoading}>
                      {profileLoading ? 'Saving…' : 'Save profile'}
                    </button>
                  </form>
                </>
              )}
            </div>
          )}
        </section>
      </main>

      <section className="lobby-preview">
        <div className="lobby-header">
          <div>
            <h2>Lobby table preview</h2>
            <p>Upcoming: live seats, auto-join, and invite links.</p>
          </div>
          <button className="primary" type="button">
            Create lobby
          </button>
        </div>
        <div className="lobby-table">
          <div className="lobby-row lobby-head">
            <span>Game</span>
            <span>Seats</span>
            <span>Status</span>
            <span>Action</span>
          </div>
          {lobbyPreview.map((lobby) => (
            <div key={lobby.id} className="lobby-row">
              <span>{lobby.game}</span>
              <span>{lobby.seats}</span>
              <span>{lobby.status}</span>
              <span>
                <button
                  className={lobby.action === 'Join lobby' ? 'primary small' : 'ghost small'}
                  type="button"
                  disabled={lobby.action === 'Full'}
                >
                  {lobby.action}
                </button>
              </span>
            </div>
          ))}
        </div>
      </section>

      <footer className="footer">
        <p>Designed for fast, friendly play. No payments in MVP.</p>
      </footer>
    </div>
  )
}

export default App
