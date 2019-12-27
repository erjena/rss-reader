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
    this.state = {
      data: [],
      sources: [],
      items: []
    }

    this.requestFeed = this.requestFeed.bind(this);
    this.submitSource = this.submitSource.bind(this);
    this.onSourceChange = this.onSourceChange.bind(this);
    this.onLogoutSuccess = this.onLogoutSuccess.bind(this);
  }

  componentDidMount() {
    this.requestFeed();
  }

  requestFeed() {
    axios.get('/api/list')
      .then((response) => {
        if (response.data === null) {
          return;
        }
        let elements = [];
        let sources = [];
        for (let i of response.data) {
          sources.push(i.sourceID);
          elements.push(...i.items);
        }
        elements.sort((a, b) => new Date(b.pubDate) - new Date(a.pubDate));
        const sourceObjects = [];
        for (let i of sources) {
          let url = new URL(i);
          sourceObjects.push({
            name: url.hostname,
            isChosen: false
          })
        }
        sourceObjects.unshift({ name: "All", isChosen: true });
        this.setState({
          data: response.data,
          sources: sourceObjects,
          items: elements
        })
      })
      .catch((error) => {
        console.log(error)
      })
  }

  submitSource(source) {
    axios.post('/api/sources', {
      link: source
    })
      .then((response) => {
        this.requestFeed();
      })
      .catch((err) => {
        if (err.response.status === 409) {
          alert('This source already exists')
        } else {
          console.log(err)
        }
      })
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

  onLogoutSuccess() {
    this.props.onLogout();
  }

  render() {
    return (
      <div className="main">
        <div className="leftColumn">
          <h2 className="greeting">Happy Reader</h2>
          <AddSource onSubmit={this.submitSource} />
          <Sources sources={this.state.sources} onSourceChange={this.onSourceChange} />
          <br />
          <Logout onLogout={this.onLogoutSuccess} />
        </div>
        <div className="rightColumn">
          <Feed data={this.state.items} />
        </div>
      </div>
    )
  }
}
