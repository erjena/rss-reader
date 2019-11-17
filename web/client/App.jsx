import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      data: []
    }
    this.requestFeed = this.requestFeed.bind(this);
  }

  componentDidMount(event) {
    this.requestFeed();
  }

  requestFeed() {
    axios.get('/list')
    .then((response) => {
      this.setState({ data: response.data.channel.items })
    })
    .catch((error) => {
      console.log(error)
    })
  }

  render() {
    const list = this.state.data.map((item, index) =>
    <li key={index}>
      {item.title}
      <br></br>
      {item.link}
    </li>
    )
    return (
      <div>
        {list}
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'))
