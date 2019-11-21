import React from 'react';
import '../public/main.css';

export default class Sources extends React.Component {
  constructor(props) {
    super(props)


  }

  render() {
    const sources = this.props.sources.map((source, index) =>
      <li key={index}>
        {source.slice(8, source.indexOf('com')+3)}
      </li>
    )
    return (
      <div>
        <div className="sourceItem">
        <span style={{paddingLeft: "15px" }}>All</span>
        <ul>{sources}</ul>
        </div>
      </div>
    )
  }
}

// style={{display: "block", paddingLeft: "25px", paddingTop: "15px" }}
