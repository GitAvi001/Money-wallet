import { useState, useEffect } from 'react'
import { useAuth } from '../context/AuthContext'
import axios from 'axios'
import './Dashboard.css'

export default function Dashboard() {
  const { user, logout } = useAuth()
  const [wallet, setWallet] = useState(null)
  const [users, setUsers] = useState([])
  const [transactions, setTransactions] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  
  // Transfer form
  const [showTransferModal, setShowTransferModal] = useState(false)
  const [transferData, setTransferData] = useState({
    receiver_id: '',
    amount: '',
    description: ''
  })

  // Add funds form
  const [showAddFundsModal, setShowAddFundsModal] = useState(false)
  const [addAmount, setAddAmount] = useState('')

  useEffect(() => {
    fetchWallet()
    fetchUsers()
    fetchTransactions()
  }, [])

  const fetchWallet = async () => {
    try {
      console.log('Fetching wallet...')
      const response = await axios.get('/api/wallet')
      console.log('Wallet response:', response.data)
      setWallet(response.data)
    } catch (err) {
      console.error('Failed to fetch wallet:', err.response?.data || err.message)
    } finally {
      setLoading(false)
    }
  }

  const fetchUsers = async () => {
    try {
      const response = await axios.get('/api/auth/users')
      setUsers(response.data.filter(u => u.id !== user.id))
    } catch (err) {
      console.error('Failed to fetch users:', err)
    }
  }

  const fetchTransactions = async () => {
    try {
      console.log('Fetching transactions...')
      const response = await axios.get('/api/transactions')
      console.log('Transactions response:', response.data)
      setTransactions(response.data)
    } catch (err) {
      console.error('Failed to fetch transactions:', err.response?.data || err.message)
    }
  }

  const handleTransfer = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    // Validate amount
    const amount = parseFloat(transferData.amount)
    if (isNaN(amount) || amount <= 0) {
      setError('Please enter a valid amount')
      return
    }

    // Check if amount exceeds balance
    if (amount > wallet.balance) {
      setError(`Insufficient funds. Your current balance is $${wallet.balance.toFixed(2)}`)
      return
    }

    // Check if recipient is selected
    if (!transferData.receiver_id) {
      setError('Please select a recipient')
      return
    }

    try {
      await axios.post('/api/transactions/transfer', {
        receiver_id: parseInt(transferData.receiver_id),
        amount: amount,
        description: transferData.description
      })
      
      setSuccess('Transfer successful!')
      setShowTransferModal(false)
      setTransferData({ receiver_id: '', amount: '', description: '' })
      
      fetchWallet()
      fetchTransactions()
      
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err.response?.data?.error || 'Transfer failed')
    }
  }

  const handleAddFunds = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    // Validate amount
    const amount = parseFloat(addAmount)
    if (isNaN(amount) || amount <= 0) {
      setError('Please enter a valid amount greater than 0')
      return
    }

    if (amount > 100000) {
      setError('Maximum amount per transaction is $100,000')
      return
    }

    try {
      console.log('Adding funds:', amount)
      const response = await axios.post('/api/wallet/add', {
        amount: amount
      })
      console.log('Add funds response:', response.data)
      
      setSuccess('Funds added successfully!')
      setShowAddFundsModal(false)
      setAddAmount('')
      
      fetchWallet()
      fetchTransactions()
      
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      console.error('Add funds error:', err.response?.data || err.message)
      setError(err.response?.data?.error || err.message || 'Failed to add funds')
      setShowAddFundsModal(false)
    }
  }

  const getReceiverName = (receiverId) => {
    const receiver = users.find(u => u.id === receiverId)
    return receiver ? receiver.name : `User #${receiverId}`
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Money Transfer Dashboard</h1>
        <div className="user-info">
          <span>Welcome, {user?.name}</span>
          <button onClick={logout} className="logout-btn">Logout</button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <div className="dashboard-content">
        {/* Wallet Card */}
        <div className="card wallet-card">
          <h2>My Wallet</h2>
          <div className="wallet-balance">
            <span className="balance-label">Current Balance</span>
            <span className="balance-amount">${wallet?.balance?.toFixed(2) || '0.00'}</span>
          </div>
          <div className="wallet-actions">
            <button onClick={() => setShowAddFundsModal(true)} className="btn btn-primary">
              Add Funds
            </button>
          </div>
        </div>

        {/* Transfer Money Card */}
        <div className="card transfer-card">
          <h2>Send Money</h2>
          <p className="card-description">Transfer money to another user</p>
          <button onClick={() => setShowTransferModal(true)} className="btn btn-success">
            Make Transfer
          </button>
        </div>
      </div>

      {/* Transaction History */}
      <div className="card transactions-card">
        <h2>Transaction History</h2>
        <div className="transactions-list">
          {transactions.length === 0 ? (
            <p className="no-transactions">No transactions yet</p>
          ) : (
            <table className="transactions-table">
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Type</th>
                  <th>Description</th>
                  <th>Amount</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map(tx => (
                  <tr key={tx.id}>
                    <td>{new Date(tx.created_at).toLocaleDateString()}</td>
                    <td className={`tx-type tx-type-${tx.transaction_type}`}>
                      {tx.transaction_type}
                    </td>
                    <td>
                      {tx.transaction_type === 'transfer' && tx.sender_id === user.id
                        ? `To ${getReceiverName(tx.receiver_id)}`
                        : tx.transaction_type === 'transfer' && tx.receiver_id === user.id
                        ? `From User #${tx.sender_id}`
                        : tx.description}
                    </td>
                    <td className={tx.sender_id === user.id && tx.transaction_type === 'transfer' ? 'amount-debit' : 'amount-credit'}>
                      {tx.sender_id === user.id && tx.transaction_type === 'transfer' ? '-' : '+'}
                      ${tx.amount.toFixed(2)}
                    </td>
                    <td>
                      <span className={`status status-${tx.status}`}>{tx.status}</span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {/* Transfer Modal */}
      {showTransferModal && (
        <div className="modal-overlay" onClick={() => setShowTransferModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Transfer Money</h3>
            <form onSubmit={handleTransfer}>
              <div className="form-group">
                <label>Select Recipient</label>
                <select
                  value={transferData.receiver_id}
                  onChange={(e) => setTransferData({...transferData, receiver_id: e.target.value})}
                  required
                >
                  <option value="">Choose recipient...</option>
                  {users.map(u => (
                    <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Amount</label>
                <input
                  type="number"
                  step="0.01"
                  min="0.01"
                  value={transferData.amount}
                  onChange={(e) => setTransferData({...transferData, amount: e.target.value})}
                  required
                  placeholder="Enter amount"
                />
              </div>

              <div className="form-group">
                <label>Description (optional)</label>
                <input
                  type="text"
                  value={transferData.description}
                  onChange={(e) => setTransferData({...transferData, description: e.target.value})}
                  placeholder="Payment for..."
                />
              </div>

              <div className="modal-actions">
                <button type="submit" className="btn btn-success">Send Money</button>
                <button type="button" onClick={() => setShowTransferModal(false)} className="btn btn-secondary">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Add Funds Modal */}
      {showAddFundsModal && (
        <div className="modal-overlay" onClick={() => setShowAddFundsModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>Add Funds to Wallet</h3>
            <form onSubmit={handleAddFunds}>
              <div className="form-group">
                <label>Amount</label>
                <input
                  type="number"
                  step="0.01"
                  min="0.01"
                  value={addAmount}
                  onChange={(e) => setAddAmount(e.target.value)}
                  required
                  placeholder="Enter amount"
                />
              </div>

              <div className="modal-actions">
                <button type="submit" className="btn btn-primary">Add Funds</button>
                <button type="button" onClick={() => setShowAddFundsModal(false)} className="btn btn-secondary">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
