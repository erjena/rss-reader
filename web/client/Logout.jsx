import React from 'react';
import axios from 'axios';
import '../public/main.css';

export default class Logout extends React.Component {
  constructor(props) {
    super(props)
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e) {
    console.log()
    axios.post('/api/logout')
    .then((response) => {
      this.props.onLogout();
    })
    .catch((error) => {
      console.log(error);
    })
  }

  render() {
    return(
      <div className="signOut">
        <button className="signOutButton" onClick={this.handleClick}>Sign out</button>
      </div>
    )
  }
}