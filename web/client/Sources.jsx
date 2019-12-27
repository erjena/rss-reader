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
    const sources = this.props.sources.map((source, index) => 
      <div id={index} key={index} className="sourceItem"
      style={{ backgroundColor: source.isChosen ? "#B0DDE4" : "transparent", paddingLeft: index === 0 ? "20px" : "38px" }}
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
