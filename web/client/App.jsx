import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';
import Feed from './Feed.jsx';
import Sources from './Sources.jsx';
import AddSource from './AddSource.jsx';
import '../public/main.css';

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      data: []
    }
    this.requestFeed = this.requestFeed.bind(this);
    this.submitSource = this.submitSource.bind(this);
  }

  componentDidMount(event) {
    this.requestFeed();
  }

  submitSource(source) {
    console.log("source", source)
    axios.post('/sources', {
      user: "abc@gmail.com",
      link: source
    })
    .then((response) => {
      console.log(response)
    })
    .catch((err) => {
      console.log(err)
    })
  }

  requestFeed() {
    axios.get('/list')
      .then((response) => {
        this.setState({ data: response.data.channel.items })
        //console.log(response.data)
      })
      .catch((error) => {
        console.log(error)
      })
  }

  render() {
    return (
      <div className="main">
        <div className="leftColumn">
          <h2 className="userName">Happy Reader</h2>
          <Sources />
          <AddSource onSubmit={this.submitSource}/>
        </div>
        <div className="rightColumn">
          <Feed data={this.state.data} />
        </div>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'))
