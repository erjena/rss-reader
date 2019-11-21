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
      data: [],
      sources: [],
      items: []
    }
    this.requestFeed = this.requestFeed.bind(this);
    this.submitSource = this.submitSource.bind(this);
    this.processData = this.processData.bind(this);
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
    axios.get('/list', {
      params: { user: "abc@gmail.com" }
    })
      .then((response) => {
        this.setState({ data: response.data })
      })
      .catch((error) => {
        console.log(error)
      })
      .finally(() => {
        this.processData()
      })
  }

  processData() {
    let elements = [];
    let sources = [];
    for (let i = 0; i < this.state.data.length; i++) {
      sources.push(this.state.data[i].sourceID);
      elements.push(...this.state.data[i].items);
    }
    this.setState({
      sources: sources,
      items: elements
    });
  }

  render() {
    return (
      <div className="main">
        <div className="leftColumn">
          <h2 className="userName">Happy Reader</h2>
          <Sources sources={this.state.sources} />
          <AddSource onSubmit={this.submitSource} />
        </div>
        <div className="rightColumn">
          <Feed data={this.state.items} />
        </div>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'))
