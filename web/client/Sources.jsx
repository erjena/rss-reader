import React from 'react';
import '../public/main.css';

export default class Sources extends React.Component {
  constructor(props) {
    super(props)

    this.onClick = this.onClick.bind(this);
  }

  onClick(event) {
    this.props.onSourceChange(event.target.id);
  }

  render() {
    console.log('inside sources component:', this.props.sources)
    const sources = this.props.sources.map((source, index) => 
      <div id={index} key={index} className="sourceItem"
      style={{ backgroundColor: source.isChosen ? "#B0DDE4" : "transparent"}}
      onClick={this.onClick}>
        {source.name}
      </div>
    )
    return (
      <div style={{ marginTop: "25px" }}>
        {sources}
      </div>
    )
  }
}
