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
      items: [],
      chosenSourse: ''
    }
    this.requestFeed = this.requestFeed.bind(this);
    this.submitSource = this.submitSource.bind(this);
    this.processData = this.processData.bind(this);
    this.modifySources = this.modifySources.bind(this);
    this.onSourceChange = this.onSourceChange.bind(this);
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
        this.modifySources()
      })
  }

  processData() {
    let elements = [];
    let sources = [];
    for (let i of this.state.data) {
      sources.push(i.sourceID);
      elements.push(...i.items);
    }
    this.setState({
      sources: sources,
      items: elements
    });
  }

  modifySources() {
    const temp = [];
    for (let i of this.state.sources) {
      temp.push({
        name: i.slice(8, i.indexOf('com')+3),
        isChosen: false
      })
    }
    temp.unshift({ name: "All", isChosen: true });
    this.setState({ sources: temp });
  }

  onSourceChange(index) {
    const idx = parseInt(index);
    for (let i = 0; i < this.state.sources.length; i++) {
      if (i === idx) {
        this.state.sources[i].isChosen = true;
      } else {
        this.state.sources[i].isChosen = false;
      }
    }
    if (idx === 0) {
      const items = [];
      for (let i of this.state.data) {
        items.push(...i.items)
      }
      this.setState({ items: items })
    } else {
      for (let i of this.state.data) {
        if ((i.sourceID).includes(this.state.sources[idx].name)) {
          this.setState({ items: i.items })
        }
      }
    }
  }

  render() {
    return (
      <div className="main">
        <div className="leftColumn">
          <h2 className="userName">Happy Reader</h2>
          <Sources sources={this.state.sources} onSourceChange={this.onSourceChange}/>
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
