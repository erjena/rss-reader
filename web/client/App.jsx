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
      loginPage: false,
      data: []
    }
    this.requestFeed = this.requestFeed.bind(this);
    this.handleLoginSuccess = this.handleLoginSuccess.bind(this);
    this.handleLogoutSuccess = this.handleLogoutSuccess.bind(this);
  }

  componentDidMount(event) {
    this.requestFeed();
  }

  requestFeed() {
    axios.get('/api/list')
      .then((response) => {
        console.log(response.data)
        this.setState({ data: response.data, loginPage: false })
      })
      .catch((error) => {
        if (error.response.status === 401) {
          this.setState({ loginPage: true })
        } else {
          console.log(error)
        }
      })
  }

  handleLoginSuccess() {
    this.setState({ loginPage: false });
  }


  handleLogoutSuccess() {
    this.setState({ loginPage: true });
  }

  render() {
    let renderPage;
    if (this.state.loginPage) {
      renderPage = <LoginPage onLoginSuccess={this.handleLoginSuccess} />
    } else {
      renderPage = <UserPage onRequestFeed={this.requestFeed} data={this.state.data} onLogout={this.handleLogoutSuccess} />
    }
    return (
      renderPage
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'))