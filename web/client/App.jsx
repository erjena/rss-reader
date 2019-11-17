import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';

class App extends React.Component {
  constructor(props) {
    super(props)

    this.requestFeed = this.requestFeed.bind(this);
  }

  componentDidMount(event) {
    this.requestFeed();
  }

  requestFeed() {
    axios.get('/list')
    .then((response) => {
      console.log(response.data)
    })
    .catch((error) => {
      console.log(error)
    })
  }

  render() {
    return (
      <div>
        <p>Hello</p>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'))
