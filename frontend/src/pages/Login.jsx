import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import './Auth.css'

export default function Login() {
  const [isRegisterMode, setIsRegisterMode] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    retypePassword: ''
  })
  const [rememberMe, setRememberMe] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)
  
  const { login, register } = useAuth()
  const navigate = useNavigate()

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
    setError('')
  }

  const handleLoginSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await login(formData.email, formData.password)
      navigate('/dashboard')
    } catch (err) {
      setError(err.response?.data?.error || 'Login failed. Please check your credentials.')
    } finally {
      setLoading(false)
    }
  }

  const handleRegisterSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    setLoading(true)

    if (formData.password !== formData.retypePassword) {
      setError('Passwords do not match')
      setLoading(false)
      return
    }

    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters')
      setLoading(false)
      return
    }

    try {
      await register(formData.name, formData.email, formData.password)
      setSuccess('Registration successful! Please check your email to verify your account.')
      setFormData({ name: '', email: '', password: '', retypePassword: '' })
      
      // Switch to login mode after 3 seconds
      setTimeout(() => {
        setIsRegisterMode(false)
        setSuccess('')
      }, 3000)
    } catch (err) {
      setError(err.response?.data?.error || 'Registration failed. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app-container">
      <div className="auth-card">
        <h2>Login/Register</h2>

        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}

        {!isRegisterMode ? (
          // Login Form
          <div className="auth-layout">
            <div className="left-section">
              <p className="new-user-text">New User?</p>
              <button 
                className="register-btn"
                onClick={() => {
                  setIsRegisterMode(true)
                  setError('')
                  setSuccess('')
                }}
              >
                Register
              </button>
            </div>

            <div className="right-section">
              <form onSubmit={handleLoginSubmit}>
                <div className="form-group">
                  <label>Email address</label>
                  <input
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    placeholder="Enter your email"
                  />
                </div>

                <div className="form-group">
                  <label>Password</label>
                  <input
                    type="password"
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    required
                    placeholder="Enter your password"
                  />
                </div>

                <div className="form-actions">
                  <button type="submit" className="login-btn" disabled={loading}>
                    {loading ? 'Logging in...' : 'Log In'}
                  </button>
                  <label className="remember-me">
                    <input
                      type="checkbox"
                      checked={rememberMe}
                      onChange={(e) => setRememberMe(e.target.checked)}
                    />
                    remember me
                  </label>
                </div>
              </form>
            </div>
          </div>
        ) : (
          // Register Form
          <div className="register-form">
            <form onSubmit={handleRegisterSubmit}>
              <div className="form-row">
                <div className="form-group">
                  <label>Name</label>
                  <input
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    placeholder="Enter your name"
                  />
                </div>

                <div className="form-group">
                  <label>Email address</label>
                  <input
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    placeholder="Enter your email"
                  />
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>Password</label>
                  <input
                    type="password"
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    required
                    placeholder="Enter password"
                  />
                </div>

                <div className="form-group">
                  <label>Retype Password</label>
                  <input
                    type="password"
                    name="retypePassword"
                    value={formData.retypePassword}
                    onChange={handleChange}
                    required
                    placeholder="Retype password"
                  />
                </div>
              </div>

              <div className="register-actions">
                <button type="submit" className="register-submit-btn" disabled={loading}>
                  {loading ? 'Registering...' : 'Register'}
                </button>
                <button 
                  type="button" 
                  className="back-to-login-btn"
                  onClick={() => {
                    setIsRegisterMode(false)
                    setError('')
                    setSuccess('')
                  }}
                >
                  Back to Login
                </button>
              </div>
            </form>
          </div>
        )}
      </div>
    </div>
  )
}
