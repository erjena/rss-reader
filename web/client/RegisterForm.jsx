import React from 'react';
import axios from 'axios';
import '../public/main.css';

export default class RegisterForm extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      email: '',
      password: '',
      verifyPassword: '',
      fields: {},
      errors: {}
    }
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
    this.handleVerifyPasswordChange = this.handleVerifyPasswordChange.bind(this);
    this.handleValidation = this.handleValidation.bind(this);
    this.handleClick = this.handleClick.bind(this);
  }

  handleEmailChange(e) {
    e.preventDefault();
    this.setState({ email: e.target.value });
  }

  handlePasswordChange(e) {
    e.preventDefault();
    this.setState({ password: e.target.value });
  }

  handleVerifyPasswordChange(e) {
    e.preventDefault();
    this.setState({ verifyPassword: e.target.value });
  }

  handleValidation() {
    let errors = {};
    let formIsValid = true;

    if (!this.state.email) {
      formIsValid = false;
      errors["email"] = "Please enter email";
    }
    const re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    if (re.test(String(this.state.email.toLowerCase()))) {
      formIsValid = false;
      errors["email"] = "Invalid email";
    }
    if (!this.state.password) {
      formIsValid = false;
      errors["password"] = "Please enter password";
    }
    if (this.state.password !== this.state.verifyPassword) {
      formIsValid = false;
      errors["verifyPassword"] = "Verify password";
    }
    this.setState({ errors: errors });
    return formIsValid;
  }

  handleClick(event) {
    event.preventDefault();
    if (this.handleValidation) {
      axios.post('/api/register', {
        username: this.state.email,
        password: this.state.password
      })
        .then((response) => {
          console.log(response)
        })
        .catch((error) => {
          console.log(error)
        })
    } else {
      alert("Please fill the form.")
    }
  }

  render() {
    return (
      <div className="forms">
        <form>
          <h4>Register</h4>
          <label>User Name
              <input type="text" onChange={this.handleEmailChange} value={this.state.email} />
          </label>
          <br />
          <span style={{ color: "red" }}>{this.state.errors["email"]}</span>
          <br />
          <label>Password
              <input type="text" onChange={this.handlePasswordChange} value={this.state.password} />
          </label>
          <br />
          <span style={{ color: "red" }}>{this.state.errors["password"]}</span>
          <br />
          <label>Verify password
              <input type="text" onChange={this.handleVerifyPasswordChange} value={this.state.verifyPassword} />
          </label>
          <br />
          <span style={{ color: "red" }}>{this.state.errors["verifyPassword"]}</span>
        </form>
        <button className="loginButton" onClick={this.handleClick}>Register</button>
      </div>
    )
  }
}
