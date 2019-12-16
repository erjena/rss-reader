import React from 'react';
import axios from 'axios';

export default class RegisterForm extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      email: '',
      password: '',
      fields: {},
      errors: {}
    }
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
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

  handleValidation() {
    let errors = this.state.errors;
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
      <div className="registerForm">
        <form>
          <fieldset>
            <label>User Name:
              <input type="text" onChange={this.handleEmailChange} value={this.state.email} />
            </label>
            <br />
            <span style={{ color: "red" }}>{this.state.errors["email"]}</span>
            <br />
            <label>Password:
              <input type="text" onChange={this.handlePasswordChange} value={this.state.password} />
            </label>
            <br />
            <span style={{ color: "red" }}>{this.state.errors["password"]}</span>
            <br />
            <button onClick={this.handleClick}>Login</button>
          </fieldset>
        </form>
      </div>
    )
  }
}
