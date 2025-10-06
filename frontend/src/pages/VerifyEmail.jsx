import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom' //useNavigate hook for navigation, useSearchParams hook for search parameters in the url
import axios from 'axios'
import './Auth.css'

export default function VerifyEmail() {
  const [searchParams] = useSearchParams()
  const [status, setStatus] = useState('verifying') // verifying, success, error
  const [message, setMessage] = useState('Verifying your email...')
  const navigate = useNavigate()

  useEffect(() => {
    const token = searchParams.get('token')
    
    if (!token) {
      setStatus('error')
      setMessage('Invalid verification link')
      return
    }

    verifyEmail(token)
  }, [searchParams])

  const verifyEmail = async (token) => {
    try {
      const response = await axios.get(`/api/auth/verify-email?token=${token}`)
      setStatus('success')
      setMessage(response.data.message || 'Email verified successfully!')
      
      // Redirect to login after 3 seconds
      setTimeout(() => {
        navigate('/login')
      }, 3000)
    } catch (err) {
      setStatus('error')
      setMessage(err.response?.data?.error || 'Verification failed. Link may be expired.')
    }
  }

  return (
    <div className="app-container">
      <div className="auth-card">
        <h2>Email Verification</h2>
        
        {status === 'verifying' && (
          <div className="verification-status">
            <div className="spinner"></div>
            <p>{message}</p>
          </div>
        )}

        {status === 'success' && (
          <div className="success-message">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
            <p>{message}</p>
            <p style={{fontSize: '14px', color: '#666'}}>Redirecting to login...</p>
          </div>
        )}

        {status === 'error' && (
          <div className="error-message">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="15" y1="9" x2="9" y2="15"></line>
              <line x1="9" y1="9" x2="15" y2="15"></line>
            </svg>
            <p>{message}</p>
            <button 
              onClick={() => navigate('/login')}
              className="btn-primary"
              style={{marginTop: '20px', padding: '12px 30px', cursor: 'pointer'}}
            >
              Go to Login
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
