import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';
import LoginPage from './LoginPage.jsx';
import UserPage from './UserPage.jsx';
import '../public/main.css';

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoggedIn: false
    }
    this.handleLoginSuccess = this.handleLoginSuccess.bind(this);
    this.handleRegisterSuccess = this.handleRegisterSuccess.bind(this);
    this.handleLogoutSuccess = this.handleLogoutSuccess.bind(this);
  }

  componentDidMount() {
    axios.get('/api/checkLoggedIn')
      .then((response) => {
        this.setState({ isLoggedIn: true });
      })
      .catch((error) => {
        console.log(error);
        this.setState({ isLoggedIn: false });
      })
  }

  handleLoginSuccess() {
    this.setState({ isLoggedIn: true });
  }

  handleRegisterSuccess() {
    this.setState({ isLoggedIn: true });
  }

  handleLogoutSuccess() {
    this.setState({ isLoggedIn: false });
  }

  render() {
    if (this.state.isLoggedIn) {
      return <UserPage onLogout={this.handleLogoutSuccess} />
    } else {
      return <LoginPage onLoginSuccess={this.handleLoginSuccess} onRegisterSuccess={this.handleRegisterSuccess} />
    }
  }
}

ReactDOM.render(<App />, document.getElementById('app'))