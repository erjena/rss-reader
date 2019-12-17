import React from 'react';
import LoginForm from './LoginForm.jsx';
import RegisterForm from './RegisterForm.jsx';
import '../public/main.css';

export default class LoginPage extends React.Component {
  constructor(props) {
    super(props)
  }
  render() {
    return (
      <div className="loginPage">
        <LoginForm onLoginSuccess={this.props.onLoginSuccess} />
        <RegisterForm />
      </div>
    )
  }
}