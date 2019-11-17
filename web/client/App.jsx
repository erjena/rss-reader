import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';
import Feed from './Feed.jsx';
import Sources from './Sources.jsx';
import '../public/main.css';

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
    return (
      <div>
        <div>
          <h1>Hello, User</h1>
        </div>
        <div className="main">
          <Sources className="sourcesContainer" />
          <Feed data={this.state.data} className="feedContainer" />
        </div>
        <div>
          <h3 style={{ textAlign: "center" }}>ABOUT</h3>
        </div>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'))
