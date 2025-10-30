import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const App = () => {
  const [order, setOrder] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [orderId, setOrderId] = useState('');

  const searchOrder = async (e) => {
    e.preventDefault();
    if (!orderId.trim()) return;

    setLoading(true);
    setError('');
    setOrder(null);

    try {
      const response = await fetch(`${API_URL}/api/order/${orderId}`);
      if (!response.ok) {
        throw new Error('Order not found');
      }
      const orderData = await response.json();
      setOrder(orderData);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount) => {
    return `$${(amount / 100).toFixed(2)}`;
  };

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h1 style={styles.title}>📦 Order Service</h1>
        <p style={styles.subtitle}>Find your order by ID</p>
        
        <form onSubmit={searchOrder} style={styles.form}>
          <div style={styles.inputGroup}>
            <input
              type="text"
              value={orderId}
              onChange={(e) => setOrderId(e.target.value)}
              placeholder="Enter order ID (e.g., b563feb7b2b84b6test)"
              style={styles.input}
              required
            />
            <button 
              type="submit" 
              disabled={loading}
              style={styles.button}
            >
              {loading ? 'Searching...' : 'Search'}
            </button>
          </div>
        </form>

        {error && (
          <div style={styles.error}>
            ❌ {error}
          </div>
        )}

        {order && (
          <div style={styles.orderCard}>
            <div style={styles.orderHeader}>
              <h2>Order Details</h2>
              <span style={styles.orderId}>{order.order_uid}</span>
            </div>

            <div style={styles.grid}>
              <div style={styles.section}>
                <h3 style={styles.sectionTitle}>📋 Order Info</h3>
                <div style={styles.detailItem}>
                  <strong>Track Number:</strong> {order.track_number}
                </div>
                <div style={styles.detailItem}>
                  <strong>Customer:</strong> {order.customer_id}
                </div>
                <div style={styles.detailItem}>
                  <strong>Entry:</strong> {order.entry}
                </div>
                <div style={styles.detailItem}>
                  <strong>Created:</strong> {new Date(order.date_created).toLocaleString()}
                </div>
              </div>

              <div style={styles.section}>
                <h3 style={styles.sectionTitle}>🚚 Delivery</h3>
                <div style={styles.detailItem}>
                  <strong>Name:</strong> {order.delivery.name}
                </div>
                <div style={styles.detailItem}>
                  <strong>Phone:</strong> {order.delivery.phone}
                </div>
                <div style={styles.detailItem}>
                  <strong>Address:</strong> {order.delivery.city}, {order.delivery.address}
                </div>
                <div style={styles.detailItem}>
                  <strong>Email:</strong> {order.delivery.email}
                </div>
              </div>

              <div style={styles.section}>
                <h3 style={styles.sectionTitle}>💳 Payment</h3>
                <div style={styles.detailItem}>
                  <strong>Amount:</strong> {formatCurrency(order.payment.amount)}
                </div>
                <div style={styles.detailItem}>
                  <strong>Provider:</strong> {order.payment.provider}
                </div>
                <div style={styles.detailItem}>
                  <strong>Bank:</strong> {order.payment.bank}
                </div>
                <div style={styles.detailItem}>
                  <strong>Currency:</strong> {order.payment.currency}
                </div>
              </div>
            </div>

            <div style={styles.itemsSection}>
              <h3 style={styles.sectionTitle}>🛍️ Items ({order.items.length})</h3>
              <div style={styles.itemsGrid}>
                {order.items.map((item, index) => (
                  <div key={index} style={styles.itemCard}>
                    <div style={styles.itemHeader}>
                      <strong>{item.name}</strong>
                      <span style={styles.brand}>{item.brand}</span>
                    </div>
                    <div style={styles.itemDetails}>
                      <div>Price: {formatCurrency(item.price)}</div>
                      <div>Sale: {item.sale}%</div>
                      <div>Total: {formatCurrency(item.total_price)}</div>
                      <div>Status: {item.status}</div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

const styles = {
  container: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  card: {
    background: 'white',
    borderRadius: '20px',
    padding: '40px',
    boxShadow: '0 20px 40px rgba(0,0,0,0.1)',
    width: '100%',
    maxWidth: '900px',
  },
  title: {
    textAlign: 'center',
    color: '#333',
    marginBottom: '10px',
    fontSize: '2.5rem',
  },
  subtitle: {
    textAlign: 'center',
    color: '#666',
    marginBottom: '30px',
    fontSize: '1.1rem',
  },
  form: {
    marginBottom: '30px',
  },
  inputGroup: {
    display: 'flex',
    gap: '10px',
  },
  input: {
    flex: 1,
    padding: '15px',
    border: '2px solid #e1e5e9',
    borderRadius: '10px',
    fontSize: '16px',
    outline: 'none',
    transition: 'border-color 0.3s',
  },
  button: {
    padding: '15px 30px',
    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    color: 'white',
    border: 'none',
    borderRadius: '10px',
    fontSize: '16px',
    cursor: 'pointer',
    fontWeight: 'bold',
    transition: 'transform 0.2s',
  },
  error: {
    background: '#fee',
    color: '#c33',
    padding: '15px',
    borderRadius: '10px',
    marginBottom: '20px',
    textAlign: 'center',
  },
  orderCard: {
    border: '2px solid #f0f0f0',
    borderRadius: '15px',
    padding: '30px',
  },
  orderHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '30px',
    paddingBottom: '20px',
    borderBottom: '2px solid #f0f0f0',
  },
  orderId: {
    background: '#f0f0f0',
    padding: '8px 15px',
    borderRadius: '20px',
    fontSize: '14px',
    fontFamily: 'monospace',
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
    gap: '30px',
    marginBottom: '30px',
  },
  section: {
    padding: '20px',
    background: '#f8f9fa',
    borderRadius: '10px',
  },
  sectionTitle: {
    marginBottom: '15px',
    color: '#333',
    fontSize: '1.2rem',
  },
  detailItem: {
    marginBottom: '8px',
    color: '#555',
  },
  itemsSection: {
    padding: '20px',
    background: '#f8f9fa',
    borderRadius: '10px',
  },
  itemsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
    gap: '15px',
  },
  itemCard: {
    background: 'white',
    padding: '15px',
    borderRadius: '8px',
    border: '1px solid #e1e5e9',
  },
  itemHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: '10px',
  },
  brand: {
    background: '#667eea',
    color: 'white',
    padding: '2px 8px',
    borderRadius: '10px',
    fontSize: '12px',
  },
  itemDetails: {
    fontSize: '14px',
    color: '#666',
  },
};

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);