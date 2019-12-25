import React from 'react';
import axios from 'axios';
import Feed from './Feed.jsx';
import Sources from './Sources.jsx';
import AddSource from './AddSource.jsx';
import Logout from './Logout.jsx';
import '../public/main.css';

export default class UserPage extends React.Component {
  constructor(props) {
    super(props)

    this.submitSource = this.submitSource.bind(this);
    this.onSourceChange = this.onSourceChange.bind(this);
    this.generateInitialState = this.generateInitialState.bind(this);
    this.onLogoutSuccess = this.onLogoutSuccess.bind(this);
    this.state = this.generateInitialState();
  }

  submitSource(source) {
    axios.post('/api/sources', {
      link: source
    })
      .then((response) => {
        this.props.onRequestFeed();
      })
      .catch((err) => {
        console.log(err)
      })
  }

  generateInitialState() {
    let elements = [];
    let sources = [];
    for (let i of this.props.data) {
      sources.push(i.sourceID);
      elements.push(...i.items);
    }
    elements.sort((a, b) => new Date(b.pubDate) - new Date(a.pubDate));

    const sourceObjects = [];
    for (let i of sources) {
      let url = new URL(i);
      console.log("url", url)
      sourceObjects.push({
        name: url.hostname,
        isChosen: false
      })
    }
    sourceObjects.unshift({ name: "All", isChosen: true });
    return {
      sources: sourceObjects,
      items: elements
    };
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
      for (let i of this.props.data) {
        items.push(...i.items)
      }
      this.setState({ items: items })
    } else {
      for (let i of this.props.data) {
        if ((i.sourceID).includes(this.state.sources[idx].name)) {
          this.setState({ items: i.items })
        }
      }
    }
  }

  onLogoutSuccess() {
    this.props.onLogout();
  }

  render() {
    return (
      <div className="main">
        <div className="leftColumn">
          <h2 className="userName">Happy Reader</h2>
          <Sources sources={this.state.sources} onSourceChange={this.onSourceChange}/>
          <br/>
          <Logout onLogout={this.onLogoutSuccess} />
          <AddSource onSubmit={this.submitSource} />
        </div>
        <div className="rightColumn">
          <Feed data={this.state.items} />
        </div>
      </div>
    )
  }
}
